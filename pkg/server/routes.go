package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/teamroffe/farm/pkg/drinks"
	"github.com/teamroffe/farm/pkg/pumps"
)

//handleLiquid handles /v1/liquids & /v1/liquid/:id
func (server *FarmServer) handleLiquid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if params["id"] == "" {
		var liquids []drinks.Liquid
		results, err := server.DB.Query("SELECT id, liquid_name, url FROM liquids")
		if err != nil {
			glog.Error(err.Error())
			return
		}
		defer results.Close()
		for results.Next() {
			var liquid drinks.Liquid
			err = results.Scan(&liquid.ID, &liquid.Name, &liquid.URL)
			if err != nil {
				glog.Error(err.Error())
				return
			}
			liquids = append(liquids, liquid)

		}
		json.NewEncoder(w).Encode(liquids)
	} else {
		var liquid drinks.Liquid
		err := server.DB.QueryRow("SELECT id, liquid_name FROM liquids WHERE id = ?", params["id"]).Scan(&liquid.ID, &liquid.Name)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(404)
				resp := &farmResponse{
					Status:  404,
					Message: "Liquid not found",
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
			glog.Error(err.Error())
			return
		}
		json.NewEncoder(w).Encode(liquid)
	}
	return
}

func (server *FarmServer) getCategories(w http.ResponseWriter, r *http.Request) {
	var drinkCategories []drinks.Category

	results, err := server.DB.Query("SELECT id,name FROM categories")
	if err != nil {
		glog.Error(err.Error())
		return
	}
	defer results.Close()
	for results.Next() {
		var category drinks.Category
		err = results.Scan(&category.ID, &category.Name)
		if err != nil {
			glog.Error(err.Error())
			return
		}
		drinkCategories = append(drinkCategories, category)

	}
	json.NewEncoder(w).Encode(drinkCategories)
	return
}

//handleDink handles /v1/drinks & /v1/drink/:id
func (server *FarmServer) handleDrink(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if params["id"] != "" {
		var drink drinks.Drink
		params := mux.Vars(r)
		drinkID, err := strconv.Atoi(params["id"])
		if err != nil {
			resp := &farmResponse{
				Status:  503,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			glog.Error(err.Error())
			return
		}
		stmtOut, err := server.DB.Prepare("SELECT id, drink_name, description, url FROM drinks WHERE id = ? LIMIT 1")
		if err != nil {
			resp := &farmResponse{
				Status:  503,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			glog.Error(err.Error())
			return
		}
		defer stmtOut.Close()
		err = stmtOut.QueryRow(drinkID).Scan(&drink.ID, &drink.Name, &drink.Description, &drink.URL)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				w.WriteHeader(404)
				resp := &farmResponse{
					Status:  404,
					Message: "Drink not found",
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
			w.WriteHeader(500)
			resp := &farmResponse{
				Status:  500,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		ingredients, err := server.getingredients(drinkID)
		if err != nil {
			resp := &farmResponse{
				Status:  503,
				Message: err.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			glog.Error(err.Error())
			return
		}

		resp := &drinkResponse{
			ID:          drink.ID,
			Name:        drink.Name,
			Description: drink.Description,
			URL:         drink.URL,
			Ingredients: ingredients,
		}

		json.NewEncoder(w).Encode(resp)
	} else {
		var drinkList []drinks.Drink

		results, err := server.DB.Query("SELECT id, drink_name, url, category, description FROM drinks")
		if err != nil {
			glog.Error(err.Error())
			return
		}
		defer results.Close()
		for results.Next() {
			var drink drinks.Drink

			err = results.Scan(&drink.ID, &drink.Name, &drink.URL, &drink.Category, &drink.Description)
			if err != nil {
				glog.Error(err.Error())
				return
			}
			drinkList = append(drinkList, drink)

		}
		json.NewEncoder(w).Encode(drinkList)
		return
	}
}

//handleOur handles /v1/pour/:id
func (server *FarmServer) handlePour(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if server.Status.Pouring || time.Since(*server.Status.LastPour) < 3*time.Second {
		resp := &farmResponse{
			Status:  509,
			Message: "Pouring limit exceeded, please try again",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	drinkID, err := strconv.Atoi(params["id"])
	if err != nil {
		resp := &farmResponse{
			Status:  503,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp, err := server.getingredients(drinkID)
	if err != nil {
		respns := &farmResponse{
			Status:  503,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(respns)
		return
	}

	if resp == nil {
		msg := fmt.Sprintf("Drink not found")
		respns := &farmResponse{
			Status:  404,
			Message: msg,
		}
		json.NewEncoder(w).Encode(respns)
		return
	}

	err = server.hasLiquids(resp)
	if err != nil {
		respns := &farmResponse{
			Status:  409,
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(respns)
		return
	}

	var wg sync.WaitGroup

	go func() {
		server.pourON()
		defer server.pourOFF()
		for _, liquid := range resp {
			wg.Add(1)
			glog.Infof("Pouring ingredient ID: %d DrinkID: %d LiquidID: %d LiquidName: %s Volume: %d\n", *liquid.ID, *liquid.DrinkID, *liquid.LiquidID, *liquid.LiquidName, *liquid.Volume)
			go func(liq *drinks.DrinkIngredient) {
				defer wg.Done()
				done := make(chan bool)
				duration := time.Duration(int64(time.Second) * int64(*liq.Volume))
				job := &pumps.PumpMSG{
					LiquidID: liq.LiquidID,
					Time:     duration,
					Done:     done,
				}
				server.PM.Queue <- job
				<-done

			}(liquid)
		}
		wg.Wait()
	}()

	respo := &farmResponse{
		Status:  200,
		Message: "Pouring started",
	}

	json.NewEncoder(w).Encode(respo)
}

//getingredients fetches the ingredients for a specified drink by ID
func (server *FarmServer) getingredients(drinkID int) ([]*drinks.DrinkIngredient, error) {
	var ingredients []*drinks.DrinkIngredient
	results, err := server.DB.Query("select drink_ingredients.id as ingredient_id, drinks.id as drink_id, liquids.liquid_name, drink_ingredients.liquid_id, drink_ingredients.volume from drinks left join drink_ingredients on drinks.id = drink_ingredients.drink_id left join liquids on liquids.id = drink_ingredients.liquid_id where drinks.id = ?;", drinkID)
	if err != nil {
		return ingredients, err
	}
	defer results.Close()
	for results.Next() {
		var ingredient drinks.DrinkIngredient
		err = results.Scan(&ingredient.ID, &ingredient.DrinkID, &ingredient.LiquidName, &ingredient.LiquidID, &ingredient.Volume)
		if err != nil {
			return ingredients, err
		}
		ingredients = append(ingredients, &ingredient)
	}
	return ingredients, nil
}

//handlePorts handles /v1/ports
func (server *FarmServer) handlePorts(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(server.PM.Ports)
}
