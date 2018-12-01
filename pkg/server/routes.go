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

func (server *FarmServer) getLiquids(w http.ResponseWriter, r *http.Request) {
	var liquids []drinks.Liquid
	results, err := server.DB.Query("SELECT id, liquid_name, url FROM liquids")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer results.Close()
	for results.Next() {
		var liquid drinks.Liquid
		err = results.Scan(&liquid.ID, &liquid.Name, &liquid.URL)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		} else {
			liquids = append(liquids, liquid)
		}
	}
	json.NewEncoder(w).Encode(liquids)
	return
}

func (server *FarmServer) getLiquid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var liquid drinks.Liquid
	liquidID, err := strconv.Atoi(params["id"])
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	err = server.DB.QueryRow("SELECT id, liquid_name FROM liquids WHERE id = ?", liquidID).Scan(&liquid.ID, &liquid.Name)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	json.NewEncoder(w).Encode(liquid)
	return
}

func (server *FarmServer) getCategories(w http.ResponseWriter, r *http.Request) {
	var drinkCategories []drinks.Category

	results, err := server.DB.Query("SELECT id,name FROM categories")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer results.Close()
	for results.Next() {
		var category drinks.Category
		err = results.Scan(&category.ID, &category.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		} else {
			drinkCategories = append(drinkCategories, category)
		}
	}
	json.NewEncoder(w).Encode(drinkCategories)
	return
}

func (server *FarmServer) getDrinks(w http.ResponseWriter, r *http.Request) {
	var drinkList []drinks.Drink

	//	results, err := server.DB.Query("SELECT id, drink_id, liquid_id, volume FROM drink_ingredients WHERE drink_id = ?", drinkID)
	results, err := server.DB.Query("SELECT id, drink_name, url, category, description FROM drinks")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer results.Close()
	for results.Next() {
		var drink drinks.Drink
		// for each row, scan the result into our tag composite object
		err = results.Scan(&drink.ID, &drink.Name, &drink.URL, &drink.Category, &drink.Description)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		} else {
			drinkList = append(drinkList, drink)
		}
		// and then print out the tag's Name attribute
	}
	json.NewEncoder(w).Encode(drinkList)
	return
}

func (server *FarmServer) getDrink(w http.ResponseWriter, r *http.Request) {
	var drink drinks.Drink
	params := mux.Vars(r)
	drinkID, _ := strconv.Atoi(params["id"])
	stmtOut, err := server.DB.Prepare("SELECT id, drink_name, description, url FROM drinks WHERE id = ? LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer stmtOut.Close()
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
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	ingredients, err := server.getingredients(drinkID)
	if err != nil {
		panic(err)
	}

	resp := &drinkResponse{
		ID:          drink.ID,
		Name:        drink.Name,
		Description: drink.Description,
		URL:         drink.URL,
		Ingredients: ingredients,
	}

	json.NewEncoder(w).Encode(resp)
}

type drinkResponse struct {
	ID          *int    `json:"id"`
	Name        *string `json:"drink_name"`
	Description *string `json:"description"`
	URL         *string `json:"url"`
	Ingredients *[]drinks.DrinkIngredient
}

func (server *FarmServer) pour(w http.ResponseWriter, r *http.Request) {
	if server.Status.Pouring {
		resp := &pourResponse{
			Status:  509,
			Message: "Pouring limit exceeded, please try again",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	params := mux.Vars(r)

	server.Status.mux.Lock()
	defer server.Status.mux.Unlock()
	server.Status.Pouring = true
	server.Status.mux.Unlock()

	if server.RpiHW {
		server.Relay1.High()
		time.Sleep(5000 * time.Millisecond)
		server.Relay1.Low()
	}
	message := fmt.Sprintf("pouring id %s", params["id"])

	resp := &pourResponse{
		Status:  200,
		Message: message,
	}

	json.NewEncoder(w).Encode(resp)

	server.Status.mux.Lock()
	server.Status.Pouring = false

}
