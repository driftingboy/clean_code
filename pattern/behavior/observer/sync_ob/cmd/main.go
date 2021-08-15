package main

import "github/litao-2071/clean_code/pattern/behavior/observer/sync_ob"

func main() {
	regSubject := sync_ob.NewRegisterSubject()
	regSubject.Subscribe(&sync_ob.RegNotifyObserver{})
	regSubject.Subscribe(&sync_ob.RegPromotionObserver{})
	regSubject.Notify()
}
