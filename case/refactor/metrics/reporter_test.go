package metrics

import (
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
)

func Test_spec(t *testing.T) {
	assert.Equal(t, spec(60), "@every 60s")
	assert.Equal(t, spec(0), "")
}

func Test_specTime(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := specTime(tt.args.t); got != tt.want {
				t.Errorf("specTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scheduledReporter_doAnalysisAndReport(t *testing.T) {
	type fields struct {
		c    *cron.Cron
		data data
		v    Viewer
	}
	type args struct {
		startTime time.Time
		endTime   time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &scheduledReporter{
				c:    tt.fields.c,
				data: tt.fields.data,
				v:    tt.fields.v,
			}
			s.doAnalysisAndReport(tt.args.startTime, tt.args.endTime)
		})
	}
}
