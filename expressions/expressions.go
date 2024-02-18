package expressions

import (
	"regexp"
	"strconv"
	"strings"
	"ya-calc/db"
	"ya-calc/utils"
)

func SplitHumanExpressionToTokens(expression string) []string {
	re := regexp.MustCompile(`(\d+\.\d+|\d+|\+|-|\*|/)`)
	return re.FindAllString(expression, -1)
}

func TokensToRPN(tokens []string) utils.Queue {
	precedence := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
	}

	var output utils.Queue
	var operators utils.Stack

	for _, token := range tokens {
		if _, ok := precedence[token]; ok {
			for !operators.IsEmpty() && precedence[operators.Head()] >= precedence[token] {
				output.Put(operators.Pop())
			}
			operators.Push(token)
		} else {
			output.Put(token)
		}
	}

	for !operators.IsEmpty() {
		output.Put(operators.Pop())
	}

	return output
}

func StrToFloat32(s string) float32 {
	val, _ := strconv.ParseFloat(s, 32)
	return float32(val)
}

func SplitRPNToComputations(tokens utils.Queue) []db.Computation {
	var nums utils.Stack
	var computations []db.Computation
	var comp db.Computation
	var left_link, right_link string
	var left, right float32
	var agentId int32 = 0
	var i int

	db.ExpressionId.Increment()
	exprID := strconv.Itoa(db.ExpressionId.GetValue())
	startCompID := db.ComputationId.GetValue() + 1

	for !tokens.IsEmpty() {
		token := tokens.Get()

		if strings.Contains("+-*/", token) {
			left, right = 0, 0
			right_link = nums.Pop()
			if right_link[0] != '@' {
				right = StrToFloat32(right_link)
				right_link = ""
			} else {
				i, _ = strconv.Atoi(right_link[1:])
				comp = computations[i-startCompID]
				comp.NotifyAgentId = agentId
				computations[i-startCompID] = comp
			}
			left_link = nums.Pop()
			if left_link[0] != '@' {
				left = StrToFloat32(left_link)
				left_link = ""
			} else {
				i, _ = strconv.Atoi(left_link[1:])
				comp = computations[i-startCompID]
				comp.NotifyAgentId = agentId
				computations[i-startCompID] = comp
			}

			db.ComputationId.Increment()
			compID := strconv.Itoa(db.ComputationId.GetValue())
			computations = append(computations,
				db.Computation{
					Id:                compID,
					ExpressionId:      exprID,
					LeftArgument:      left,
					LeftArgumentLink:  left_link,
					RightArgument:     right,
					RightArgumentLink: right_link,
					Operation:         token,
				})
			nums.Push("@" + compID)
			agentId = (agentId + 1) % db.AgentsCnt
		} else {
			nums.Push(token)
		}
	}

	last_index := len(computations) - 1
	comp = computations[last_index]
	comp.Final = true
	computations[last_index] = comp

	return computations
}

func ExpressionComputations(expression string) []db.Computation {
	return SplitRPNToComputations(TokensToRPN(SplitHumanExpressionToTokens(expression)))
}
