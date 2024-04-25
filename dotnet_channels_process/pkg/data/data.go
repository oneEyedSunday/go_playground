package data

type DataWithKey struct {
	Key  string
	Data string
}

type IDataProcessor interface {
	Schedule(DataWithKey)
}
