package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/teamroffe/farm/pkg/drinks"
)

func (server *FarmServer) getDrink(w http.ResponseWriter, r *http.Request) {
	var drink drinks.Drink
	params := mux.Vars(r)
	drinkID, _ := strconv.Atoi(params["id"])
	// Prepare statement for reading data
	stmtOut, err := server.DB.Prepare("SELECT id, drink_name, description, url FROM drinks WHERE id = ? LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer stmtOut.Close()
	// Query the square-number of 13
	err = stmtOut.QueryRow(drinkID).Scan(&drink.ID, &drink.Name, &drink.Description, &drink.URL)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			w.WriteHeader(404)
			resp := &pourResponse{
				Status:  404,
				Message: "Drink not found",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(500)
		resp := &pourResponse{
			Status:  500,
			Message: "Drink not found",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	ingredients := server.getingredients(drinkID)

	resp := &drinkInfo{
		Info:        &drink,
		Ingredients: &ingredients,
	}

	json.NewEncoder(w).Encode(resp)
}

func (server *FarmServer) pour(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	server.Status.mux.Lock()
	defer server.Status.mux.Unlock()
	if server.Status.Pouring {
		resp := &pourResponse{
			Status:  509,
			Message: "Pouring limit exceeded, please try again",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	server.Status.Pouring = true
	server.Status.mux.Unlock()
	server.Relay1.High()
	time.Sleep(5000 * time.Millisecond)
	server.Relay1.Low()
	message := fmt.Sprintf("pouring id %s", params["id"])
	resp := &pourResponse{
		Status:  200,
		Message: message,
	}
	json.NewEncoder(w).Encode(resp)
	server.Status.mux.Lock()
	server.Status.Pouring = false

}
