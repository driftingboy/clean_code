package metrics_test

import (
	"github/litao-2071/clean_code/case/refactor/metrics"
	"time"
)

func Example_metrics() {
	// data := metrics.NewMemData()
	// view := metrics.NewConsoleViewer()

	// cron
	if err := metrics.StartReport(metrics.NewConsoleReporter(2, 2)); err != nil {
		panic(err)
	}

	// collector
	collector := metrics.NewCollector()
	collector.Record(metrics.RecordInfo{ApiName: "register", ResponseDurationInMillis: 90, RequestTimestamp: time.Now().Unix() - 3}) // 不在时间范围内
	collector.Record(metrics.RecordInfo{ApiName: "register", ResponseDurationInMillis: 500, RequestTimestamp: time.Now().Unix()})
	collector.Record(metrics.RecordInfo{ApiName: "register", ResponseDurationInMillis: 500, RequestTimestamp: time.Now().Unix()})
	collector.Record(metrics.RecordInfo{ApiName: "register", ResponseDurationInMillis: 100, RequestTimestamp: time.Now().Unix()})
	collector.Record(metrics.RecordInfo{ApiName: "login", ResponseDurationInMillis: 10, RequestTimestamp: time.Now().Unix()})
	collector.Record(metrics.RecordInfo{ApiName: "login", ResponseDurationInMillis: 30, RequestTimestamp: time.Now().Unix()})

	time.Sleep(time.Second * 3)
	// Output:
	// {
	//   "login": {
	//     "avg": 20,
	//     "count": 2,
	//     "max": 30,
	//     "min": 10,
	//     "tps": 1
	//   },
	//   "register": {
	//     "avg": 366.66666,
	//     "count": 3,
	//     "max": 500,
	//     "min": 100,
	//     "tps": 1.5
	//   }
	// }
}

// 建议不要使用TestMain, 使用只是测output 写为example, 与外部有交互直接main函数（demo）
// func TestMain(m *testing.M) {
// 	data := metrics.NewMemData()

// 	// cron
// 	if err := metrics.StartReport(metrics.NewConsoleReporter(data, 10, 10)); err != nil {
// 		panic(err)
// 	}

// 	// collector
// 	collector := metrics.NewCollector(data)
// 	collector.Record(metrics.RecordInfo{ApiName: "register", ResponseDurationInMillis: 90, RequestTimestamp: time.Now().Unix() - 20}) // 不在时间范围内
// 	collector.Record(metrics.RecordInfo{ApiName: "register", ResponseDurationInMillis: 500, RequestTimestamp: time.Now().Unix()})
// 	collector.Record(metrics.RecordInfo{ApiName: "register", ResponseDurationInMillis: 500, RequestTimestamp: time.Now().Unix()})
// 	collector.Record(metrics.RecordInfo{ApiName: "register", ResponseDurationInMillis: 100, RequestTimestamp: time.Now().Unix()})
// 	collector.Record(metrics.RecordInfo{ApiName: "login", ResponseDurationInMillis: 10, RequestTimestamp: time.Now().Unix()})
// 	collector.Record(metrics.RecordInfo{ApiName: "login", ResponseDurationInMillis: 30, RequestTimestamp: time.Now().Unix()})

// 	time.Sleep(time.Second * 11)
// 	os.Exit(m.Run())

// 	os.Exit(m.Run())
// }
