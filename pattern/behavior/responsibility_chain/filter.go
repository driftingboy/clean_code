package responsibilitychain

import (
	"fmt"
)

// type HandFunc func(context.Context) error

// type HandlersChain []HandFunc

// func (h HandlersChain) add(hfs ...HandFunc) {
// 	h = append(h, hfs...)
// }

type SensitiveWordFilter interface {
	doFilter(content string) bool
}

type SensitiveWordFilterChain struct {
	filters []SensitiveWordFilter
}

func (chain *SensitiveWordFilterChain) AddFilter(filter SensitiveWordFilter) *SensitiveWordFilterChain {
	chain.filters = append(chain.filters, filter)
	return chain
}

func (chain *SensitiveWordFilterChain) DoFilters(context string) bool {
	for _, filter := range chain.filters {
		if !filter.doFilter(context) {
			return false
		}
	}
	return true
}

// AdSensitiveWordFilter 广告
type AdWordFilter struct{}

// Filter 实现过滤算法
func (f *AdWordFilter) doFilter(content string) bool {
	// logic...
	fmt.Println("ad word filter succeed")
	return true
}

// 色情暴力词过滤
type SexAndViolenceWordFilter struct{}

// Filter 实现过滤算法
func (f *SexAndViolenceWordFilter) doFilter(content string) bool {
	// logic...
	fmt.Println("sex、 violence word filter failed")
	return false
}
