package app

type Config struct {
	GRPCBindAddress string

	Brokers []string

	TopicCUDUsers     string
	GroupIDCUDUsers   string
	PartitionCUDUsers int

	TopicCUDTasks     string
	PartitionCUDTasks int
}
