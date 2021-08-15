package strategy

import "fmt"

type DiscountStrategy interface {
	CalDiscount() float32
}

type NormalDiscount struct{}

func (n *NormalDiscount) CalDiscount() float32 {
	fmt.Println("normal discount")
	return 2.2
}

type GroupDiscount struct{}

func (n *GroupDiscount) CalDiscount() float32 {
	fmt.Println("group discount")
	return 3.3
}

type PromotionDiscount struct{}

func (n *PromotionDiscount) CalDiscount() float32 {
	fmt.Println("promotion discount")
	return 4.4
}

// Strategy factory
var Strategies = map[orderTyp]DiscountStrategy{
	normal:    &NormalDiscount{},
	group:     &GroupDiscount{},
	promotion: &PromotionDiscount{},
}

func NewDiscountStrategy(typ orderTyp) DiscountStrategy {
	return Strategies[typ]
}
