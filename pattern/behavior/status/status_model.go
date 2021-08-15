package status

/*
对于操作固定，状态不多，但是每个操作都非常复杂的情况；
建议使用状态模式，通过接口，将每个状态的操作实现相互隔离
状态模式通过将事件触发的状态转移和动作执行，拆分到不同的状态类中，来避免分支判断逻辑。
*/

type Consumer interface {
	SmallPurchase()
	LargePurchase()
	NoConsumptionForAMonth()
}

type CivilianConsumer struct {
	// stateMachine
}

type GoldConsumer struct {
	// stateMachine
}

type PlatinumConsumer struct {
	// stateMachine
}

// impl Consumer action
