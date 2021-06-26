package metrics

import "time"

type data interface {
	SaveRecordInfo(ri RecordInfo) bool
	GetRecordInfos(startTime, endTime time.Time) (apiNameRecordInfosMap map[string]RecordInfos)
	GetRecordInfosWithApi(apiName string, startTime, endTime time.Time) RecordInfos
}

var _ data = (*memData)(nil)

// test
type memData struct {
	apiNameRecordInfosMap map[string]RecordInfos
}

func NewMemData() data {
	return &memData{
		apiNameRecordInfosMap: make(map[string]RecordInfos),
	}
}

func (md *memData) SaveRecordInfo(ri RecordInfo) bool {
	data := md.apiNameRecordInfosMap
	data[ri.ApiName] = append(data[ri.ApiName], ri)
	return true
}

func (md *memData) GetRecordInfos(startTime, endTime time.Time) (apiNameRecordInfosMap map[string]RecordInfos) {
	apiNameRecordInfosMap = make(map[string]RecordInfos)
	for apiName, recordInfos := range md.apiNameRecordInfosMap {
		for _, recordInfo := range recordInfos {
			if startTime.Unix() <= recordInfo.RequestTimestamp && recordInfo.RequestTimestamp < endTime.Unix() {
				apiNameRecordInfosMap[apiName] = append(apiNameRecordInfosMap[apiName], recordInfo)
			}
		}
	}
	return
}

func (md *memData) GetRecordInfosWithApi(apiName string, startTime, endTime time.Time) RecordInfos {
	rs, ok := md.apiNameRecordInfosMap[apiName]
	if !ok {
		return nil
	}
	result := make(RecordInfos, 0, rs.Len())
	copy(result, rs)

	return result
}
