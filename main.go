package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"ya-calc/db"
	"ya-calc/expressions"
	"ya-calc/utils"
	"ya-calc/workers-pilot"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	i, _ := strconv.Atoi(os.Getenv("AGENTS_CNT"))
	db.AgentsCnt = int32(i)
	go func() {
		workers_pilot.StartControlling()
	}()

	http.HandleFunc("/expression", handleExpression)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("I am listening at http://localhost:8080!")
}

func handleExpression(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		var expr db.Expression
		err = json.NewDecoder(r.Body).Decode(&expr)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		var comp db.Computation
		var i int32
		i = 0
		comps := expressions.ExpressionComputations(expr.StrExpr)
		expr.Id = comps[0].ExpressionId
		expr.ComputationsCnt = len(comps)
		expr.ComputationsReadyCnt = 0
		expr.State = "Calculating"
		db.Expressions[expr.Id] = expr

		for _, comp = range comps {
			db.Computations[comp.Id] = comp
			utils.PostMessage(comp, "NewComputations", "localhost:9092", i)
			i = (i + 1) % db.AgentsCnt
		}

		_, err = fmt.Fprintf(w, expr.Id)
		if err != nil {
			panic(err)
		}
	} else if r.Method == http.MethodGet {
		e := db.Expressions[r.URL.Query().Get("id")]
		type PartialExpression struct {
			Id              string
			ComputationsCnt int
			State           string
			Result          float32
		}

		partialExpr := PartialExpression{
			Id:              e.Id,
			ComputationsCnt: e.ComputationsCnt,
			State:           e.State,
			Result:          e.Result,
		}

		jsonData, err := json.Marshal(partialExpr)
		if err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
		return
	}
}
