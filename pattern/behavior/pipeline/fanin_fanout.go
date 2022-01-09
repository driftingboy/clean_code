package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// create two sample message and stop channels
	mc1, stop1 := generate("message from generator 1", 200*time.Millisecond)
	mc2, stop2 := generate("message from generator 2", 300*time.Millisecond)

	// multiplex message channels
	mmc, wgmult := multiplex(mc1, mc2)

	// create exit Signal channel for graceful shutdown
	// wait for interrupt or terminate signal
	exitSignal := make(chan error)

	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
		exitSignal <- fmt.Errorf("%s signal received", <-sc)
	}()

	// wait for multiplexed messages
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for m := range mmc {
			fmt.Println(m)
		}
	}()

	// wait for exit Signal
	if err := <-exitSignal; err != nil {
		fmt.Println(err.Error())
		close(exitSignal)
	}

	// 停止生产数据
	stop1()
	stop2()

	// 等待剩余数据聚合接收结束
	wgmult.Wait()

	// 关闭聚合通道
	close(mmc)

	// 等待消费聚合数据结束
	wg.Wait()
}

// 这里还可以通过 ctx 批量
func generate(message string, interval time.Duration) (mc chan string, stop func()) {
	mc = make(chan string)
	sc := make(chan struct{})
	stop = func() { close(sc) }

	go func() {
		// 退出时关闭消息通道
		defer func() {
			close(mc)
		}()

		for {
			select {
			case _, ok := <-sc: // close 或者 sc <- struct{}{} 都可以收到信号
				fmt.Printf("ok %v \n", ok)
				return
			default:
				time.Sleep(interval)

				mc <- message
			}
		}
	}()

	return
}

// 返回 stopall 方法，关闭所有 mcs，在关掉自己
func multiplex(mcs ...chan string) (chan string, *sync.WaitGroup) {
	mmc := make(chan string)
	wg := &sync.WaitGroup{}

	for _, mc := range mcs {
		wg.Add(1)
		go func(mc chan string) {
			for m := range mc {
				mmc <- m
			}
			wg.Done()
		}(mc)
	}

	return mmc, wg
}
