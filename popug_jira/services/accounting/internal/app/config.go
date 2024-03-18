package app

type Config struct {
	GRPCBindAddress string

	Brokers       []string
	EventsGroupID string

	TopicCUDUsers     string
	PartitionCUDUsers int

	TopicCUDTasks     string
	PartitionCUDTasks int

	TopicBETasks     string
	PartitionBETasks int

	TopicCUDTxs     string
	PartitionCUDTxs int
}
