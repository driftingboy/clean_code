package metrics

import "sort"

type RequestStatus struct {
	Max   float32 `json:"max"`
	Min   float32 `json:"min"`
	Avg   float32 `json:"avg"`
	Count int64   `json:"count"` // big.int
	Tps   float64 `json:"tps"`
}

// 如果数据量大可以起多个goroutine
func Analysis(timeRange int64, apiNameRecordInfosMap map[string]RecordInfos) map[string]*RequestStatus {
	rsMap := make(map[string]*RequestStatus, len(apiNameRecordInfosMap))
	for apiName, recordInfos := range apiNameRecordInfosMap {
		rsMap[apiName] = doAnalysis(recordInfos, timeRange) // Analysis静态没法mock
	}
	return rsMap
}

// max min 可以使用最小堆，查询复杂度 o1， 插入o logN
// 但是插入操作很频繁，统计查找操作又很少是运营平台查询或主动定时推送，所以还是使用数组保证插入不影响业务代码性能

// 当max、min、...统计操作变得复杂，最好抽离成单独的函数独立测试
func doAnalysis(records RecordInfos, timeRangeInSec int64) (rs *RequestStatus) {
	// 注意要考虑是否拷贝一份recordInfos，这里只是做统计（只会使用一次），不会存储，所以可以不用拷贝
	rs = &RequestStatus{}
	recordsCount := records.Len()
	if recordsCount == 0 {
		return
	}

	// request count、tps
	rs.Count = int64(recordsCount)
	rs.Tps = float64(recordsCount) / float64(timeRangeInSec)

	// response max、min、avg
	sort.Sort(records)
	rs.Min = records[0].ResponseDurationInMillis
	rs.Max = records[recordsCount-1].ResponseDurationInMillis
	var sum float32
	for _, r := range records {
		sum += r.ResponseDurationInMillis
	}
	rs.Avg = sum / float32(recordsCount)

	return rs
}

// TODO
// 使用最小堆 实现一个topN效果（不仅可以是内存中的结构，也可以是数据库结构）
// 为了防止影响业务代码，异步的更新这个堆中的数据
