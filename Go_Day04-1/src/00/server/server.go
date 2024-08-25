package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Order struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

type Respond struct {
	Thanks string `json:"thanks,omitempty"`
	Change int    `json:"change,omitempty"`
	Error  string `json:"error,omitempty"`
}

func BuyCandy(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Problem with decoding json", http.StatusBadRequest)
		return
	}
	candies := map[string]int{
		"CE": 10,
		"AA": 15,
		"NT": 17,
		"DE": 21,
		"YR": 23,
	}
	// w.Header().Set("Content-Type", "application/json")
	price, ok := candies[order.CandyType]
	if !ok {
		respond := Respond{Error: "Problem with decoding json"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respond)
		return
	}

	if order.CandyCount < 0 {
		respond := Respond{Error: fmt.Sprintln("invalid candy count")}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(respond)
		return
	}

	if order.Money < order.CandyCount*price {
		respond := Respond{Error: fmt.Sprintf("You need %d more money", order.CandyCount*price-order.Money)}
		w.WriteHeader(http.StatusPaymentRequired)
		json.NewEncoder(w).Encode(respond)
		return
	}

	if order.CandyCount*price <= order.Money {
		respond := Respond{Thanks: "Thank you!", Change: order.Money - order.CandyCount*price}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(respond)

	}

}
