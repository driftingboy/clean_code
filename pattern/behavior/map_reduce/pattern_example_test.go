package pattern

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func Test_MapReduce(t *testing.T) {
	userIds := []string{"1", "2", "3", "4", "5"}
	suitableUsers := Filter(userIds, func(userId string) bool {
		id, err := strconv.Atoi(userId)
		if err != nil {
			return false
		}
		if (id & 1) == 0 {
			return true
		}
		return false
	})
	fmt.Printf("suitableUsers: %+v\n", suitableUsers)

	userKeys := MapStr2Str(suitableUsers, func(userId string) string {
		return "userId|" + userId
	})
	fmt.Printf("userKeys: %+v\n", userKeys)

	totalLen := Sum(userKeys, func(userKey string) int {
		return len(userKey)
	})
	fmt.Printf("totalLen: %+v\n", totalLen)

	// generic map
	nums := []int{1, 2, 3, 4}
	squared_arr := Map(nums, func(x int) int {
		return x * x
	})
	fmt.Println(squared_arr)

	strs := []string{"Hao", "Chen", "MegaEase"}
	upstrs := Map(strs, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Println(upstrs)
}
