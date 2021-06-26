package metrics

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hokaccha/go-prettyjson"
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

type consoleReporter struct {
	timeCycleInSec int64
	timeRangeInSec int64
	c              *cron.Cron
	data           data
}

// timeCycle：统计时间周期
// timeRange：统计时间范围
// example：NewConsoleReporter(data, 60, 60) 每60秒统计前60秒
func NewConsoleReporter(d data, timeCycle, timeRange int64) *consoleReporter {
	return &consoleReporter{
		timeCycleInSec: timeCycle,
		timeRangeInSec: timeRange,
		c:              cron.New(),
		data:           d,
	}
}

func (cr *consoleReporter) Start() error {
	_, err := cr.c.AddFunc(spec(cr.timeCycleInSec), func() {
		// 获取数据
		endTime := time.Now()
		startTime := endTime.Add(-1 * time.Duration(cr.timeRangeInSec) * time.Second)
		apiNameRecordInfosMap := cr.data.GetRecordInfos(startTime, endTime)

		// 解析数据
		rsMap := Analysis(cr.timeRangeInSec, apiNameRecordInfosMap)

		// 报道数据
		data, err := prettyjson.Marshal(rsMap)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("%v", string(data))
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
	c              *cron.Cron
	emailAddr      string
	emailPwd       string
	data           data
}

// timeCycle：统计时间周期
// timeRange：统计时间范围
// example：NewConsoleReporter(data, 60, 60) 每60秒统计前60秒
func NewWebReporter(d data, emailAddr, emailPwd string, timeCycle, timeRange int64) *emailReporter {
	return &emailReporter{
		// execTime: execTime,
		timeCycleInSec: timeCycle,
		timeRangeInSec: timeRange,
		c:              cron.New(),
		emailAddr:      emailAddr,
		emailPwd:       emailPwd,
		data:           d,
	}
}

func (cr *emailReporter) Start() error {
	_, err := cr.c.AddFunc(spec(cr.timeCycleInSec), func() {
		// 获取数据
		endTime := time.Now()
		startTime := endTime.Add(-1 * time.Duration(cr.timeRangeInSec) * time.Second)
		apiNameRecordInfosMap := cr.data.GetRecordInfos(startTime, endTime)

		// 解析数据
		rsMap := Analysis(cr.timeRangeInSec, apiNameRecordInfosMap)

		// report data: 邮箱的发送逻辑可能比较复杂，这时候可以抽出一个viewer interface
		// monitor send to email
		fmt.Printf("emailAddr: %s send success. %v", cr.emailAddr, rsMap)
	})
	if err != nil {
		return err
	}

	cr.c.Start()
	return nil
}

func spec(timeCycleInSec int64) string {
	if timeCycleInSec <= 0 {
		return ""
	}
	return "@every " + strconv.FormatInt(timeCycleInSec, 10) + "s"
}

func specTime(t time.Time) string {
	return ""
}
