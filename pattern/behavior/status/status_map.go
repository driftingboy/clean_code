package status

// 假设有一个平台币充值功能
// 状态：平民，黄金会员，白金会员
// 事件：充值100平台币、充值500平台币、长时间不购买
// 动作：不同状态触发不同事件产生不同动作，3*3种动作
/*
|              | 平民          | 黄金会员      | 白金会员      |
| ------------ | ------------- | ------------ | ------------- |
| 一次小额购买   | +100/-        | +200/-       | +500/-        |
| 一次大额购买 	 | +500/黄金会员  | +700/白金会员  | +1000/-       |
| 一个月未消费   |       -       | -100/平民     | -300/黄金会员 |
*/

var (
	CurrentLevel = Civilian
	CurrentToken = 0
)

type ConsumerLevel int

const (
	Civilian ConsumerLevel = iota
	Gold
	Platinum
)

// action
const (
	SmallPurchase = iota
	LargePurchase
	NoConsumptionForAMonth
)

// StateMachine
var (
	levelTable = [][]ConsumerLevel{
		{Civilian, Gold, Platinum},
		{Gold, Platinum, Platinum},
		{Civilian, Civilian, Gold},
	}
	actionTable = [][]int{
		{100, 200, 500},
		{500, 700, 1000},
		{0, -100, -300},
	}
)

func AppearSmallPurchase() {
	executeAction(SmallPurchase)
}

func AppearLargePurchase() {
	executeAction(LargePurchase)
}

func AppearNoConsumptionForAMonth() {
	executeAction(NoConsumptionForAMonth)
}

func executeAction(action int) {
	nowLevel := levelTable[action][CurrentLevel]
	addToken := actionTable[action][CurrentLevel]
	CurrentLevel = nowLevel
	CurrentToken += addToken
}
