package interpreter

import (
	"fmt"
	"strconv"
	"strings"
)

/*
告警器负责根据用户提供的数据，判断是否到达用户定义的告警规则，如果达到则发出警告
其中有三个要点
- 获取用户数据
  1.通过用户主动上报数据
  2.用户提供符合规定的接口，在告警服务注册，由告警服务拉取数据
- 判断是否达到告警阀值
  这是告警器的核心功能，我们可以使用解释器模式去实现
- 告警方式
  邮件、短信、微信、日志...

这里我们通过第二部分，告警语句解析判断来介绍 解释器模式
自定义的告警规则只包含“||、&&、>、<、==”这五个运算符，其中，“>、<、==”运算符的优先级高于“||、&&”运算符，“&&”运算符优先级高于“||”。
在表达式中，任意元素之间需要通过空格来分隔。
*/

// 告警规则表达式
type IExpression interface {
	// 判断入参是否符合报警规则
	Interpret(data map[string]int) bool
}

type AlterRuleInterpreter struct {
	exp IExpression
}

func NewAlterRuleInterpreter(rule string) (*AlterRuleInterpreter, error) {
	orExp, err := NewOrExpression(rule)
	if err != nil {
		return nil, err
	}
	return &AlterRuleInterpreter{exp: orExp}, nil
}

func (ar *AlterRuleInterpreter) Interpret(data map[string]int) bool {
	return ar.exp.Interpret(data)
}

// &&
type AndExpression struct {
	exps []IExpression
}

func NewAndExpression(exp string) (*AndExpression, error) {
	strExps := strings.Split(exp, "&&")
	exps := make([]IExpression, 0, len(strExps))
	for _, se := range strExps {
		if strings.Contains(se, ">") {
			ge, err := NewGreaterExpression(se)
			if err != nil {
				return nil, err
			}
			exps = append(exps, ge)
		} else if strings.Contains(se, "<") {
			le, err := NewLessExpression(se)
			if err != nil {
				return nil, err
			}
			exps = append(exps, le)
		} else if strings.Contains(se, "=") {
			ee, err := NewEqualExpression(se)
			if err != nil {
				return nil, err
			}
			exps = append(exps, ee)
		}
	}

	return &AndExpression{
		exps: exps,
	}, nil
}

func (ae AndExpression) Interpret(data map[string]int) bool {
	for _, e := range ae.exps {
		if !e.Interpret(data) {
			return false
		}
	}

	return true
}

// ||
type OrExpression struct {
	exps []IExpression
}

func NewOrExpression(exp string) (*OrExpression, error) {
	strExps := strings.Split(exp, "||")
	exps := make([]IExpression, 0, len(strExps))
	for _, se := range strExps {
		andExp, err := NewAndExpression(se)
		if err != nil {
			return nil, err
		}
		exps = append(exps, andExp)
	}

	return &OrExpression{exps: exps}, nil
}

func (oe OrExpression) Interpret(data map[string]int) bool {
	for _, e := range oe.exps {
		if e.Interpret(data) {
			return true
		}
	}

	return false
}

// >
type GreaterExpression struct {
	key string
	val int
}

func NewGreaterExpression(exp string) (*GreaterExpression, error) {
	k, v, err := getKVByExpression(exp, ">")
	if err != nil {
		return nil, err
	}

	return &GreaterExpression{
		key: k,
		val: v,
	}, nil
}

func (g GreaterExpression) Interpret(data map[string]int) bool {
	v, ok := data[g.key]
	if !ok {
		return false
	}
	return v > g.val
}

// <
type LessExpression struct {
	key string
	val int
}

func NewLessExpression(exp string) (*LessExpression, error) {
	k, v, err := getKVByExpression(exp, "<")
	if err != nil {
		return nil, err
	}

	return &LessExpression{
		key: k,
		val: v,
	}, nil
}

func (g LessExpression) Interpret(data map[string]int) bool {
	v, ok := data[g.key]
	if !ok {
		return false
	}
	return v < g.val
}

// =
type EqualExpression struct {
	key string
	val int
}

func NewEqualExpression(exp string) (*EqualExpression, error) {
	k, v, err := getKVByExpression(exp, "=")
	if err != nil {
		return nil, err
	}

	return &EqualExpression{
		key: k,
		val: v,
	}, nil
}

func (g EqualExpression) Interpret(data map[string]int) bool {
	v, ok := data[g.key]
	if !ok {
		return false
	}
	return v == g.val
}

func getKVByExpression(exp string, operator string) (key string, val int, err error) {
	if len(operator) != 1 {
		return "", 0, fmt.Errorf("operator is invalid: %s", exp)
	}
	data := strings.Split(strings.TrimSpace(exp), " ")
	if len(data) != 3 || data[1] != operator {
		return "", 0, fmt.Errorf("exp is invalid: %s", exp)
	}

	val, err = strconv.Atoi(data[2])
	if err != nil {
		return "", 0, fmt.Errorf("exp is invalid: %s err %v", exp, err)
	}
	key = data[0]

	return
}
