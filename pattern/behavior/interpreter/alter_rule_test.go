package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlterRuleInterpreter_Interpret(t *testing.T) {

	ar, _ := NewAlterRuleInterpreter("req_total_peer_sec < 1000 || req_error_peer_sec > 100 && flag = 5")

	type args struct {
		data map[string]int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "case1",
			args: args{
				data: map[string]int{
					"req_total_peer_sec": 800, // true
					"req_error_peer_sec": 200, // true
					"flag":               5,   // true
				},
			},
			want: true,
		},
		{name: "case4", args: args{
			data: map[string]int{
				// "req_total_peer_sec":800, 不填则直接认为 false，因为规则是业务定的，如果key缺失，说明有问题；默认不填为0，可能导致传输问题一直不被发现
				"req_error_peer_sec": 80, // false
				"flag":               4,  //false
			},
		}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ar.Interpret(tt.args.data)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewInterpret(t *testing.T) {
	// 测试 > < = && ||
	data := map[string]int{
		"a": 1,
		"b": 2,
	}
	tests := []struct {
		name    string
		rule    string
		wantErr bool
		want    bool
	}{
		{name: "test > ", rule: "a > 2", want: false},
		{name: "test < ", rule: "a < 2", want: true},
		{name: "test = ", rule: "a = 1", want: true},
		{name: "test && 1", rule: "a = 1 && b > 2", want: false},
		{name: "test && 2", rule: "a = 1 && b = 2", want: true},
		{name: "test || 1", rule: "a < 1 || b < 2", want: false},
		{name: "test || 2", rule: "a < 1 || b = 2", want: true},
		{name: "test &| priority", rule: "a < 1 || b = 2 && a = 1", want: true},        // false || true
		{name: "test invalid format", rule: "  a < 1 || b = 2 && a = 1  ", want: true}, // false || true
		{name: "test invalid format", rule: "a  < 1", wantErr: true},                   // err
		{name: "test invalid format", rule: "a  < 1 ||", wantErr: true},                // err
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ar, err := NewAlterRuleInterpreter(tt.rule)
			assert.Equal(t, tt.wantErr, err != nil)
			if err != nil {
				return
			}

			got := ar.Interpret(data)
			assert.Equal(t, tt.want, got)
		})
	}
}
