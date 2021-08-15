package strategy

func ExampleOrder_Discount() {
	o1 := &Order{Id: "1", Typ: 0, Desc: "常规订单"}
	o2 := &Order{Id: "2", Typ: 1, Desc: "团购订单"}
	o1.Discount()
	o2.Discount()
	//Output:
	//normal discount
	//group discount
}
