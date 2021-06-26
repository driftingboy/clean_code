package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_spec(t *testing.T) {
	assert.Equal(t, spec(60), "@every 60s")
	assert.Equal(t, spec(0), "")
}
