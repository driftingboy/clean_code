package metrics

type collector struct {
	data data
}

type CollectorOption func(*collector)

func WitCollectorData(d data) CollectorOption {
	return func(c *collector) {
		c.data = d
	}
}

func NewCollector(opts ...CollectorOption) *collector {
	c := &collector{
		data: NewMemData(),
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (c *collector) Record(ri RecordInfo) {
	c.data.SaveRecordInfo(ri)
}

type RecordInfo struct {
	ApiName                  string
	ResponseDurationInMillis float32
	RequestTimestamp         int64
}

type RecordInfos []RecordInfo

func (rs RecordInfos) Len() int      { return len(rs) }
func (rs RecordInfos) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }
func (rs RecordInfos) Less(i, j int) bool {
	return rs[i].ResponseDurationInMillis < rs[j].ResponseDurationInMillis
}
