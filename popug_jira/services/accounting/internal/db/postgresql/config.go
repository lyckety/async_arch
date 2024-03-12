package postgresql

import (
	"fmt"
)

type Config struct {
	DBHost         string
	DBPort         uint16
	DBName         string
	DBUsername     string
	DBPassword     string
	MigrationsPath string
}

func (c *Config) getMigrationsPath() string {
	return c.MigrationsPath
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable port=%d host=%s",
		c.DBUsername,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		c.DBHost,
	)
}
