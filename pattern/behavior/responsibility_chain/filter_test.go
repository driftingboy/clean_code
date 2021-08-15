package responsibilitychain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Filter(t *testing.T) {
	chain := &SensitiveWordFilterChain{}
	chain.AddFilter(&AdWordFilter{}).AddFilter(&SexAndViolenceWordFilter{})
	b := chain.DoFilters("ad1 ........xxxxx")
	assert.False(t, b)
}
