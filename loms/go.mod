module route256.ozon.ru/project/loms

go 1.22.0

replace (
	route256.ozon.ru/pkg/logger => ../pkg/logger
	route256.ozon.ru/pkg/metrics => ../pkg/metrics
	route256.ozon.ru/pkg/tracer => ../pkg/tracer
)

require (
	github.com/IBM/sarama v1.43.1
	github.com/dnwe/otelsarama v0.0.0-20240308230250-9388d9d40bc0
	github.com/envoyproxy/protoc-gen-validate v1.0.4
	github.com/exaring/otelpgx v0.5.4
	github.com/gojuno/minimock/v3 v3.3.6
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	github.com/jackc/pgx/v5 v5.5.5
	github.com/sethvargo/go-envconfig v1.0.1
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.50.0
	go.opentelemetry.io/contrib/propagators/jaeger v1.25.0
	go.opentelemetry.io/otel v1.25.0
	go.uber.org/goleak v1.3.0
	go.uber.org/zap v1.27.0
	google.golang.org/genproto/googleapis/api v0.0.0-20240314234333-6e1732d8331c
	google.golang.org/grpc v1.62.1
	google.golang.org/protobuf v1.33.0
	route256.ozon.ru/pkg/logger v0.0.0-00010101000000-000000000000
	route256.ozon.ru/pkg/metrics v0.0.0-00010101000000-000000000000
	route256.ozon.ru/pkg/tracer v0.0.0-00010101000000-000000000000
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.6.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/joegasewicz/status-writer v1.0.0 // indirect
	github.com/klauspost/compress v1.17.7 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.19.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0 // indirect
	go.opentelemetry.io/otel/metric v1.25.0 // indirect
	go.opentelemetry.io/otel/sdk v1.25.0 // indirect
	go.opentelemetry.io/otel/trace v1.25.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240311132316-a219d84964c2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
