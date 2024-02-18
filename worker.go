package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
	"ya-calc/db"
	"ya-calc/utils"
)

var waiters []db.Computation
var computationsResults = make(map[string]float32)

func ProcessNewComputation(data []byte) {
	var computation db.Computation
	err := json.Unmarshal(data, &computation)
	if err != nil {
		println(err)
		return
	}

	if value, ok := computationsResults[computation.LeftArgumentLink]; ok {
		computation.LeftArgument = value
		delete(computationsResults, computation.RightArgumentLink)
		computation.LeftArgumentLink = ""
	}

	if value, ok := computationsResults[computation.RightArgumentLink]; ok {
		computation.RightArgument = value
		delete(computationsResults, computation.RightArgumentLink)
		computation.RightArgumentLink = ""
	}

	if computation.LeftArgumentLink != "" {
		waiters = append(waiters, computation)
	} else if computation.RightArgumentLink != "" {
		waiters = append(waiters, computation)
	} else {
		Calc(computation)
	}
}

func ProcessCalculatedComputation(data []byte) {
	var computationResult db.ComputationResult
	var computation db.Computation
	var i int
	err := json.Unmarshal(data, &computationResult)
	if err != nil {
		println(err)
		return
	}

	for i, computation = range waiters {
		if computation.LeftArgumentLink == "@"+computationResult.ComputationId {
			computation.LeftArgument = computationResult.Result
			computation.LeftArgumentLink = ""
			if computation.RightArgumentLink == "" {
				Calc(computation)
			}
		} else if computation.RightArgumentLink == "@"+computationResult.ComputationId {
			computation.RightArgument = computationResult.Result
			computation.RightArgumentLink = ""
			if computation.LeftArgumentLink == "" {
				Calc(computation)
			}
		} else {
			computationsResults["@"+computationResult.ComputationId] = computationResult.Result
		}
		waiters[i] = computation
	}

	if len(waiters) == 0 {
		computationsResults["@"+computationResult.ComputationId] = computationResult.Result
	}

}

func SendEventComputationCalculated(expressionId string, computationId string, notifyAgentId int32, result float32, final bool) {
	if final {
		utils.PostMessage(db.ExpressionResult{expressionId, result}, "CalculatedExpressions", os.Getenv("KAFKA_URL"), 0)
	} else {
		utils.PostMessage(db.ComputationResult{computationId, result}, "CalculatedComputations", os.Getenv("KAFKA_URL"), notifyAgentId)
	}
}

func Calc(computation db.Computation) {
	var n int
	var res float32

	switch computation.Operation {
	case "+":
		n, _ = strconv.Atoi(os.Getenv("TIME_ADD"))
		res = computation.LeftArgument + computation.RightArgument
		break
	case "-":
		n, _ = strconv.Atoi(os.Getenv("TIME_SUBSTR"))
		res = computation.LeftArgument - computation.RightArgument
	case "*":
		n, _ = strconv.Atoi(os.Getenv("TIME_MULT"))
		res = computation.LeftArgument * computation.RightArgument
	case "/":
		n, _ = strconv.Atoi(os.Getenv("TIME_DIVISION"))
		res = computation.LeftArgument / computation.RightArgument
	}

	time.Sleep(time.Duration(n) * time.Second)
	SendEventComputationCalculated(computation.ExpressionId, computation.Id, computation.NotifyAgentId, res, computation.Final)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	t, _ := strconv.Atoi(os.Args[1])
	ind := int32(t)

	forever := make(chan bool)
	go func() {
		utils.SubscribeHandlerToTopic("NewComputations", ind, ProcessNewComputation, os.Getenv("KAFKA_URL"), "workers")
	}()
	go func() {
		utils.SubscribeHandlerToTopic("CalculatedComputations", ind, ProcessCalculatedComputation, os.Getenv("KAFKA_URL"), "workers")
	}()
	fmt.Println("Ready!")
	<-forever
}
