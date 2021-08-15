package status

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StatusMap(t *testing.T) {
	AppearSmallPurchase()
	assert.Equal(t, CurrentLevel, Civilian)
	assert.Equal(t, CurrentToken, 100)
	AppearLargePurchase()
	assert.Equal(t, CurrentLevel, Gold)
	assert.Equal(t, CurrentToken, 600)
	AppearNoConsumptionForAMonth()
	assert.Equal(t, CurrentLevel, Civilian)
	assert.Equal(t, CurrentToken, 500)
	AppearLargePurchase()
	assert.Equal(t, CurrentLevel, Gold)
	assert.Equal(t, CurrentToken, 1000)
	AppearLargePurchase()
	assert.Equal(t, CurrentLevel, Platinum)
	assert.Equal(t, CurrentToken, 1700)
}
