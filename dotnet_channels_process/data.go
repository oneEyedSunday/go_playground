package data

type DataWithKey struct {
	Key  string
	Data string
}

type IDataProcessor interface {
	Schedule(DataWithKey)
}

type KeySpecificDataProcessor struct {
	ProcessorKey string
	c            chan DataWithKey
}

var _ IDataProcessor = (*KeySpecificDataProcessor)(nil)
