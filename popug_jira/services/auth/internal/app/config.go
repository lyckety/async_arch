package app

import "time"

type Config struct {
	GRPCBindAddress string
	MBTopic         string
	Brokers         []string
	MBPartition     int
	JWTTokenExpired time.Duration
}
