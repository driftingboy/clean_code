package metrics_test

import (
	"github/litao-2071/clean_code/case/principle/metrics"
	"time"
)

func Example_metrics() {

	data := metrics.NewMemData()

	// cron
	if err := metrics.StartReport(metrics.NewConsoleReporter(data, 2, 2)); err != nil {
		panic(err)
	}

	// collector
	collector := metrics.NewCollector(data)
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

// 优化点1.提供更方便的查询 2.非入侵代码 3.存储异步可选
// 优化类：
// 1. aggregator现在还好，以后逻辑复杂会越来越难维护，职责不单一；
// 2. reporter类多个reporter代码重复违背diy原则(已解决，抽出)；还有静态方法不好测试；如果有email发送比较复杂还需要在拆一个类

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
