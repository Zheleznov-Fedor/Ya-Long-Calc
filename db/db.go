package db

import (
	"sync"
)

type Expression struct {
	Id              string  `json:"id"`
	StrExpr         string  `json:"str_expr"`
	ComputationsCnt int     `json:"computations_cnt"`
	State           string  `json:"state"`
	Result          float32 `json:"result"`
}

type Computation struct {
	ExpressionId      string  `json:"expressionId"`
	Id                string  `json:"id"`
	LeftArgument      float32 `json:"leftArgument"`
	LeftArgumentLink  string  `json:"LeftArgumentLink"`
	Operation         string  `json:"operation"`
	RightArgument     float32 `json:"rightArgument"`
	RightArgumentLink string  `json:"RightArgumentLink"`
	Result            float32 `json:"result"`
	Final             bool    `json:"final"`
	NotifyAgentId     int32   `json:"notifyAgentId"`
}

var Expressions = make(map[string]Expression)
var Computations = make(map[string]Computation)
var AgentsCnt int32

type Counter struct {
	value int
	mu    sync.RWMutex
}

func (c *Counter) Increment() {
	c.mu.Lock()
	c.value = c.value + 1
	c.mu.Unlock()
}

func (c *Counter) GetValue() int {
	c.mu.RLock()
	data := c.value
	c.mu.RUnlock()
	return data
}

var ExpressionId = Counter{value: 10}
var ComputationId = Counter{value: 10}

type ExpressionResult struct {
	ExpressionId string  `json:"expressionIdId"`
	Result       float32 `json:"result"`
}

type ComputationResult struct {
	ComputationId string  `json:"computationId"`
	Result        float32 `json:"result"`
}
