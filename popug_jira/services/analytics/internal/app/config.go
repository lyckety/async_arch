package app

type Config struct {
	Brokers       []string
	EventsGroupID string

	TopicCUDUsers     string
	PartitionCUDUsers int

	TopicCUDTasks     string
	PartitionCUDTasks int

	TopicCUDTxs     string
	PartitionCUDTxs int
}
