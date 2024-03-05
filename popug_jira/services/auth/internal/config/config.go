package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	cleanenv "github.com/ilyakaznacheev/cleanenv"
	logger "github.com/sirupsen/logrus"
)

const (
	envConfigFileName = ".env"
)

var (
	ErrReadEnvSettings = errors.New("error read environment settings")
)

type EnvSettings struct {
	GRPCBindAddress string `env:"GRPCBINDADDR" env-default:":50051" env-description:"IP:PORT to bind grpc"`                                //nolint:lll
	LogLevel        string `env:"LOG_LEVEL" env-default:"info" env-description:"log level: trace, debug, info, warn, error, fatal, panic"` //nolint:lll

	JWTTokenExpired time.Duration `env:"JWT_TOKEN_EXPIRED" env-default:"5m" env-description:"lexpired jwt token"` //nolint:lll

	// Kafka settings
	BrokerURLs string `env:"BROKER_URLS" env-default:"127.0.0.1:19093,127.0.0.1:29093,127.0.0.1:39093" env-description:"kafka broker urls: host1:port,host2:port,host3:port"` //nolint:lll
	TopicName  string `env:"MB_TOPIC" env-required:"true" env-description:"topic for producing events"`                                                                       //nolint:lll
	Partition  int    `env:"MB_PARTITION" env-default:"0" env-description:"listen partition. Must set GROUP_ID or PARTITION"`                                                 //nolint:lll

	// Database settings
	DBHost            string `env:"DB_HOST" env-default:"localhost" env-description:"ip or domain name database host"`                            //nolint:lll
	DBPort            uint16 `env:"DB_PORT" env-default:"55432" env-description:"port for database host"`                                         //nolint:lll
	DBName            string `env:"DB_NAME" env-default:"users" env-description:"database name"`                                                  //nolint:lll
	DBUsername        string `env:"DB_USERNAME" env-default:"postgres" env-description:"database username"`                                       //nolint:lll
	DBPassword        string `env:"DB_PASSWORD" env-default:"postgres" env-description:"database password"`                                       //nolint:lll
	MigrationsSQLPath string `env:"MIGRATIONS_SQL_PATH" env-default:"./migrations/postgresql" env-description:"Path where Migration SQL scripts"` //nolint:lll
}

func (e *EnvSettings) GetHelpString() (string, error) {
	customHeader := "options which can be set via env: "

	helpString, err := cleanenv.GetDescription(e, &customHeader)
	if err != nil {
		return "", fmt.Errorf("get help string failed: %w", err)
	}

	return helpString, nil
}

type Config struct {
	env        *EnvSettings
	brokerURLs []string
}

func New() (*Config, error) {
	cfg := &Config{}

	cfg.env = &EnvSettings{} //nolint:exhaustruct

	helpString, err := cfg.env.GetHelpString()
	if err != nil {
		return nil, fmt.Errorf("getting help string of env settings failed: %w", err)
	}

	logger.Info(helpString)

	if issetEnvConfigFile() {
		if err := cleanenv.ReadConfig(envConfigFileName, cfg.env); err != nil {
			return nil, fmt.Errorf("read env cofig file failed: %w", err)
		}
	} else if err := cleanenv.ReadEnv(cfg.env); err != nil {
		return nil, fmt.Errorf("read env config failed: %w", err)
	}

	if err := cfg.validateBrokerURLs(); err != nil {
		return nil, fmt.Errorf("cfg.validateBrokerURLs(): %w", err)
	}

	return cfg, nil
}

func (c *Config) validateBrokerURLs() error {
	sliceBrokers := strings.Split(c.env.BrokerURLs, ",")

	if len(sliceBrokers) == 0 {
		return fmt.Errorf("%w: must be set broker urls", ErrReadEnvSettings)
	}

	// for _, brokerUrl := range sliceBrokers {
	// 	_, err := url.Parse(brokerUrl)
	// 	if err != nil {
	// 		return fmt.Errorf("%w: not correct broker url %q", ErrReadEnvSettings, brokerUrl)
	// 	}
	// }

	c.brokerURLs = sliceBrokers

	return nil
}

func (c *Config) GetGRPCBindAddress() string {
	return c.env.GRPCBindAddress
}

func (c *Config) GetLogLevel() logger.Level {
	lvl, err := logger.ParseLevel(c.env.LogLevel)
	if err != nil {
		logger.Error(err)

		return logger.InfoLevel
	}

	return lvl
}

func (c *Config) GetJWTTokenExpired() time.Duration {
	return c.env.JWTTokenExpired
}

func issetEnvConfigFile() bool {
	_, err := os.Stat(envConfigFileName)

	return err == nil
}

func (c *Config) GetBrokers() []string {
	return c.brokerURLs
}

func (c *Config) GetPartition() int {
	return c.env.Partition
}

func (c *Config) GetTopicName() string {
	return c.env.TopicName
}

func (c *Config) GetDBHost() string {
	return c.env.DBHost
}

func (c *Config) GetDBPort() uint16 {
	return c.env.DBPort
}

func (c *Config) GetDBUserName() string {
	return c.env.DBUsername
}

func (c *Config) GetDBPassword() string {
	return c.env.DBPassword
}

func (c *Config) GetDBName() string {
	return c.env.DBName
}

func (c *Config) GetMigrationsPath() string {
	return c.env.MigrationsSQLPath
}
