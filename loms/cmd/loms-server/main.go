package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"net/http"
	"route256.ozon.ru/project/loms/internal/config"
	pgs "route256.ozon.ru/project/loms/internal/repository/pgs"
	order_pgs_repository "route256.ozon.ru/project/loms/internal/repository/pgs/order"
	stock_pgs_repository "route256.ozon.ru/project/loms/internal/repository/pgs/stock"
	notes_usecase "route256.ozon.ru/project/loms/internal/service/loms"
	"strings"

	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc/reflection"
	"route256.ozon.ru/project/loms/internal/app/loms"
	"route256.ozon.ru/project/loms/internal/mw"
	lomsDesc "route256.ozon.ru/project/loms/pkg/api/v1"

	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	lomsConfig, err := config.GetConfig(ctx)
	if err != nil {
		panic(err)
	}

	dbConnection := connectToDB(ctx, lomsConfig)
	defer dbConnection.Close()

	grpcServer := createGRPCServer()
	controller := createLomsServer(dbConnection)

	lomsDesc.RegisterLomsServer(grpcServer, controller)

	startGRPCServer(grpcServer, lomsConfig.LomsGrpcPort)
	startHttpServer(grpcServer, lomsConfig.LomsHttpPort, lomsConfig.LomsGrpcPort)

	dbConnection.Close()
}

func connectToDB(ctx context.Context, config *config.Config) *pgs.DB {
	dbPool, err := pgs.ConnectToPgsDb(ctx, config, false)
	if err != nil {
		log.Fatalln("Cannot initialize connection to postgres")
	}

	return dbPool
}

func createGRPCServer() *grpc.Server {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mw.Panic,
			mw.Logger,
			mw.Validate,
		),
	)
	reflection.Register(grpcServer)
	return grpcServer
}

func createLomsServer(dbConnection *pgs.DB) *loms.Server {
	orderRepository := order_pgs_repository.NewOrderPgsRepository(dbConnection)
	stockRepository := stock_pgs_repository.NewStocksPgRepository(dbConnection)
	useCase := notes_usecase.NewService(orderRepository, stockRepository)
	return loms.NewServer(useCase)
}

func startGRPCServer(grpcServer *grpc.Server, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func startHttpServer(grpcServer *grpc.Server, httpPort, grpcPort string) {
	conn, err := grpc.Dial(grpcPort, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()

	mux.Handle("/", gwmux)
	serveSwagger(mux)

	if err = lomsDesc.RegisterLomsHandler(context.Background(), gwmux, conn); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	handler := grpcHandlerFunc(grpcServer, mux)
	handler = mw.WithHTTPLoggingMiddleware(mw.WithHTTPCorsMiddleware(handler)) // todo chain
	gwServer := &http.Server{
		Addr:    httpPort,
		Handler: handler,
	}
	log.Printf("Serving gRPC-Gateway on %s\n", gwServer.Addr)
	log.Fatal(gwServer.ListenAndServe())
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}
func serveSwagger(mux *http.ServeMux) {
	prefix := "/docs/"
	mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir("./swagger-ui"))))
	mux.HandleFunc(prefix+"swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/openapiv2/loms.swagger.json")
	})
}
