package metrics

import (
	"reflect"
	"testing"
	"time"
)

func TestAnalysis(t *testing.T) {
	type args struct {
		records        RecordInfos
		timeRangeInSec int64
	}
	tests := []struct {
		name   string
		args   args
		wantRs *RequestStatus
	}{
		{
			name:   "zero",
			args:   args{records: []RecordInfo{}, timeRangeInSec: 2},
			wantRs: &RequestStatus{Max: 0.0, Min: 0.0, Avg: 0.0, Count: 0, Tps: 0.0},
		},
		{
			name: "base",
			args: args{
				records: []RecordInfo{
					{ApiName: "register", ResponseDurationInMillis: 100, RequestTimestamp: time.Now().Unix()},
					{ApiName: "register", ResponseDurationInMillis: 50.2, RequestTimestamp: time.Now().Unix()},
				}, timeRangeInSec: 2,
			},
			wantRs: &RequestStatus{Max: 100.0, Min: 50.2, Avg: 75.1, Count: 2, Tps: 1.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRs := doAnalysis(tt.args.records, tt.args.timeRangeInSec); !reflect.DeepEqual(gotRs, tt.wantRs) {
				t.Errorf("Analysis() = %v, want %v", gotRs, tt.wantRs)
			}
		})
	}
}
