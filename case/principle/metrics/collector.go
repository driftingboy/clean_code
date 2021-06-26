package metrics

type collector struct {
	data data
}

func NewCollector(data data) *collector {
	return &collector{
		data: data,
	}
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
