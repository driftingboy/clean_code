package strategy

const (
	normal orderTyp = iota
	group
	promotion
)

type orderTyp int

type Order struct {
	Id   string   `json:"id,omitempty"`
	Typ  orderTyp `json:"typ,omitempty"`
	Desc string   `json:"desc,omitempty"`
}

// 动态的更具条件选择不同策略才是策略模式的应用场景
// 否则，静态注入的方式只是面向接口编程
func (o *Order) Discount() float32 {
	ds := NewDiscountStrategy(o.Typ)
	return ds.CalDiscount()
}
