package main

import "github/litao-2071/clean_code/pattern/behavior/observer/base"

func main() {
	regSubject := base.NewRegisterSubject()
	regSubject.Subscribe(&base.RegNotifyObserver{})
	regSubject.Subscribe(&base.RegPromotionObserver{})
	regSubject.Notify()
}
