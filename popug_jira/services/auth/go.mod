module github.com/lyckety/async_arch/popug_jira/services/auth

go 1.21.4

replace github.com/lyckety/async_arch/popug_jira/schema-registry => ../../schema-registry/go

require (
	github.com/golang-migrate/migrate/v4 v4.16.2
	github.com/ilyakaznacheev/cleanenv v1.4.2
	github.com/lyckety/async_arch/popug_jira/schema-registry v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.42
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/sync v0.3.0
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.33.0
	gorm.io/driver/postgres v1.5.2
	gorm.io/gorm v1.25.2
)

require (
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/stretchr/testify v1.8.3 // indirect
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.6.0
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.3.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/klauspost/compress v1.16.4 // indirect
	github.com/lib/pq v1.10.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230525234030-28d5490b6b19 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
