package eventbus

import (
	"fmt"
	"reflect"
	"sync"
	// "github.com/gammazero/workerpool"
)

type Bus interface {
	Subscribe(topic string, handler interface{}) error
	Publish(topic string, args ...interface{})
}

var _ Bus = (*SimpleBus)(nil)

type SimpleBus struct {
	mu       sync.RWMutex
	handlers map[string][]reflect.Value
	// wp workerpool 可以考虑加入协程池
}

// async
func NewSimpleBus() *SimpleBus {
	return &SimpleBus{
		handlers: make(map[string][]reflect.Value),
	}
}

func (sb *SimpleBus) Subscribe(topic string, handler interface{}) error {
	sb.mu.Lock()
	defer sb.mu.Unlock()

	v := reflect.ValueOf(handler)
	if v.Type().Kind() != reflect.Func {
		return fmt.Errorf("handler type %s is not reflect.func", v.Type().Kind())
	}
	// topic []handler
	sb.handlers[topic] = append(sb.handlers[topic], v)
	return nil
}

func (sb *SimpleBus) Publish(topic string, args ...interface{}) {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	handlers, ok := sb.handlers[topic]
	if !ok || len(handlers) == 0 {
		return
	}

	callArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		callArgs[i] = reflect.ValueOf(arg)
	}

	for i := range handlers {
		go handlers[i].Call(callArgs)
	}
}

// handler config
// type eventHandler struct {
// 	callBack      reflect.Value
// 	flagOnce      bool
// 	async         bool
// 	transactional bool
// 	sync.Mutex    // lock for an event handler - useful for running async callbacks serially
// }
