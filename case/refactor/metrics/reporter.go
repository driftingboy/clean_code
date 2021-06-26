package metrics

import (
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
)

// func StartReport(rs ...Report) {
// 	for _, r := range rs {
// 		go func() {
// 			err := r(60,60)
// 		}()
// 	}
// }

// type Report func(timeCycle, timeRange int) error
type reporter interface {
	Start() error
}

func StartReport(rs ...reporter) error {
	for _, r := range rs {
		if err := r.Start(); err != nil {
			return err
		}
	}
	return nil
}

var _ reporter = (*consoleReporter)(nil)

type consoleReporterOption func(*consoleReporter)

func WithConsoleReporterData(d data) consoleReporterOption {
	return func(c *consoleReporter) {
		c.data = d
	}
}

func WithConsoleReporterViewer(v Viewer) consoleReporterOption {
	return func(c *consoleReporter) {
		c.v = v
	}
}

type consoleReporter struct {
	timeCycleInSec int64
	timeRangeInSec int64
	scheduledReporter
}

// timeCycle：统计时间周期
// timeRange：统计时间范围
// example：NewConsoleReporter(data, 60, 60) 每60秒统计前60秒
func NewConsoleReporter(timeCycle, timeRange int64, opts ...consoleReporterOption) *consoleReporter {
	cr := &consoleReporter{
		timeCycleInSec: timeCycle,
		timeRangeInSec: timeRange,
		scheduledReporter: scheduledReporter{
			c:    cron.New(),
			data: NewMemData(),
			v:    NewConsoleViewer(),
		},
	}

	for _, o := range opts {
		o(cr)
	}
	return cr
}

func (cr *consoleReporter) Start() error {
	_, err := cr.c.AddFunc(spec(cr.timeCycleInSec), func() {
		// 获取数据
		endTime := time.Now()
		startTime := endTime.Add(-1 * time.Duration(cr.timeRangeInSec) * time.Second)
		cr.doAnalysisAndReport(startTime, endTime)
	})
	if err != nil {
		return err
	}

	cr.c.Start()
	return nil
}

type emailReporter struct {
	timeCycleInSec int64
	timeRangeInSec int64
	scheduledReporter
}

// timeCycle：统计时间周期
// timeRange：统计时间范围
// example：NewEmailReporter(data,"1425xxxxxx@qq.com","xxxxxxxx", 60, 60) 每60秒统计前60秒
func NewEmailReporter(d data, v Viewer, timeCycle, timeRange int64) *emailReporter {
	return &emailReporter{
		// execTime: execTime,
		timeCycleInSec: timeCycle,
		timeRangeInSec: timeRange,
		scheduledReporter: scheduledReporter{
			c:    cron.New(),
			data: d,
			v:    v,
		},
	}
}

func (cr *emailReporter) Start() error {
	_, err := cr.c.AddFunc(spec(cr.timeCycleInSec), func() {
		endTime := time.Now()
		startTime := endTime.Add(-1 * time.Duration(cr.timeRangeInSec) * time.Second)
		cr.doAnalysisAndReport(startTime, endTime)
	})
	if err != nil {
		return err
	}

	cr.c.Start()
	return nil
}

// 每多长时间周期执行一次
func spec(timeCycleInSec int64) string {
	if timeCycleInSec <= 0 {
		return ""
	}
	return "@every " + strconv.FormatInt(timeCycleInSec, 10) + "s"
}

// 在什么时间点执行
func specTime(t time.Time) string {
	return ""
}

type scheduledReporter struct {
	c    *cron.Cron
	data data
	v    Viewer
}

func (s *scheduledReporter) doAnalysisAndReport(startTime, endTime time.Time) {
	// 获取数据
	apiNameRecordInfosMap := s.data.GetRecordInfos(startTime, endTime)

	// 解析数据
	duration := endTime.Sub(startTime)
	rsMap := Analysis(int64(duration.Seconds()), apiNameRecordInfosMap)

	// report data
	s.v.Output(rsMap)
}
