package workers_pilot

import (
	"encoding/json"
	"ya-calc/db"
	"ya-calc/utils"
)

func ProcessCalculatedExpression(data []byte) {
	var res db.ExpressionResult
	var e db.Expression
	err := json.Unmarshal(data, &res)
	if err != nil {
		println(err)
		return
	}
	e = db.Expressions[res.ExpressionId]
	e.State = "Ready"
	e.Result = res.Result
	db.Expressions[res.ExpressionId] = e
}

func StartControlling() {
	utils.SubscribeHandlerToTopic("CalculatedExpressions", 0, ProcessCalculatedExpression, "localhost:9092", "workers")

}
