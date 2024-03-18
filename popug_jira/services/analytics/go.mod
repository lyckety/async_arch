module github.com/lyckety/async_arch/popug_jira/services/analytics

go 1.21.4

replace github.com/lyckety/async_arch/popug_jira/schema-registry => ../../schema-registry/go

require (
	github.com/ilyakaznacheev/cleanenv v1.4.2
	github.com/lyckety/async_arch/popug_jira/schema-registry v0.0.0-00010101000000-000000000000
	github.com/segmentio/kafka-go v0.4.42
	github.com/sirupsen/logrus v1.9.3
	golang.org/x/sync v0.2.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/klauspost/compress v1.16.4 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/stretchr/testify v1.8.1 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
