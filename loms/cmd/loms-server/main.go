package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"route256.ozon.ru/project/loms/internal/app/loms"
	"route256.ozon.ru/project/loms/internal/config"
	"route256.ozon.ru/project/loms/internal/events"
	kafka_events "route256.ozon.ru/project/loms/internal/events/kafka"
	"route256.ozon.ru/project/loms/internal/mw"
	pgs "route256.ozon.ru/project/loms/internal/repository/pgs"
	order_pgs_repository "route256.ozon.ru/project/loms/internal/repository/pgs/order"
	stock_pgs_repository "route256.ozon.ru/project/loms/internal/repository/pgs/stock"
	loms_usecase "route256.ozon.ru/project/loms/internal/service/loms"
	lomsDesc "route256.ozon.ru/project/loms/pkg/api/v1"
	"strings"
	"sync"
	"syscall"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	wg := &sync.WaitGroup{}
	lomsConfig, err := config.GetConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dbConnection := connectToDB(ctx, lomsConfig)
	defer dbConnection.Close()

	grpcServer := createGRPCServer()
	eventManager := events.NewEventManager()
	kafka_events.RegisterEvents(ctx, wg, lomsConfig, eventManager)

	controller := createLomsServer(dbConnection, eventManager)
	lomsDesc.RegisterLomsServer(grpcServer, controller)

	startGRPCServer(ctx, grpcServer, lomsConfig.LomsGrpcPort)
	startHttpServer(ctx, grpcServer, controller, lomsConfig.LomsHttpPort)
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

func createLomsServer(dbConnection *pgs.DB, eventManager loms_usecase.EventManagers) *loms.Server {
	orderRepository := order_pgs_repository.NewOrderPgsRepository(dbConnection)
	stockRepository := stock_pgs_repository.NewStocksPgRepository(dbConnection)

	useCase := loms_usecase.NewService(orderRepository, stockRepository, eventManager)
	return loms.NewServer(useCase)
}

func startGRPCServer(ctx context.Context, grpcServer *grpc.Server, port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server listening at %v", lis.Addr())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	go func() {
		<-ctx.Done()
		log.Println("Shutting down gprc server")
		grpcServer.Stop()
	}()
}

func startHttpServer(ctx context.Context, grpcServer *grpc.Server, controller *loms.Server, httpPort string) {
	mux := http.NewServeMux()
	gwMux := runtime.NewServeMux()

	mux.Handle("/", gwMux)
	serveSwagger(mux)

	if err := lomsDesc.RegisterLomsHandlerServer(ctx, gwMux, controller); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	handler := grpcHandlerFunc(grpcServer, mux)
	handler = mw.WithHTTPLoggingMiddleware(mw.WithHTTPCorsMiddleware(handler)) // todo chain
	gwServer := &http.Server{
		Addr:        httpPort,
		Handler:     handler,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	log.Printf("Serving gRPC-Gateway on %s\n", gwServer.Addr)
	go func() {
		<-ctx.Done()
		log.Println("Shutting down http server")
		if err := gwServer.Shutdown(ctx); err != nil {
			log.Println("Failed to shutdown server: ", err)
		}
	}()
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
