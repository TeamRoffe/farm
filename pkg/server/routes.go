package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/teamroffe/farm/pkg/drinks"
	"github.com/teamroffe/farm/pkg/pumps"
)

//handleLiquid handles /v1/liquids & /v1/liquid/:id
func (server *FarmServer) handleLiquid(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var err error
	var results *sql.Rows
	if params["id"] == "" {
		results, err = server.DB.Query("SELECT id, liquid_name, url FROM liquids")
	} else {
		results, err = server.DB.Query("SELECT id, liquid_name, url FROM liquids WHERE id = ?", params["id"])
	}
	var liquids []drinks.Liquid
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

	var err error
	var results *sql.Rows
	var drinkList []drinkResponse

	if params["id"] == "" {
		results, err = server.DB.Query("SELECT id, drink_name, category, description, url FROM drinks")
	} else {
		results, err = server.DB.Query("SELECT id, drink_name, category, description, url FROM drinks WHERE id = ? LIMIT 1", params["id"])
	}

	if err != nil {
		glog.Error(err.Error())
		return
	}

	defer results.Close()

	for results.Next() {
		var drink drinkResponse
		err = results.Scan(&drink.ID, &drink.Name, &drink.Category, &drink.Description, &drink.URL)
		if err != nil {
			glog.Error(err.Error())
			return
		}

		ingredients, err := server.getingredients(*drink.ID)

		if err != nil {
			json.NewEncoder(w).Encode(farmResp(503, err.Error()))
			glog.Error(err.Error())
			return
		}

		drink.Ingredients = ingredients
		drinkList = append(drinkList, drink)
	}
	if len(drinkList) == 0 {
		json.NewEncoder(w).Encode(farmResp(404, "Drink not found"))
		return
	}

	js, err := json.Marshal(drinkList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	/*
		json.NewEncoder(w).Encode(drinkList)
		return
	*/
}

var reasons = []string{
	"Le om du tar den i tvåan",
	"Roffe vad håller du på med",
	"Göm spriten, Lindmark kommer!",
	"Här vare snepatchat",
	"Bosse!",
	"Roffe!",
	"Har det hänt en grej?",
	"Jag ser dig i spegeln!",
	"Tvi tvi tvi. Jag spottar vart jag vill",
	"Vem fan har hällt sprit på bastuaggregatet!?",
}

//handleOur handles /v1/pour/:id
func (server *FarmServer) handlePour(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if server.Status.Pouring || time.Since(*server.Status.LastPour) < 3*time.Second {
		json.NewEncoder(w).Encode(farmResp(509, "Pouring limit exceeded, please try again"))
		return
	}

	drinkID, err := strconv.Atoi(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(farmResp(503, err.Error()))
		return
	}
	resp, err := server.getingredients(drinkID)
	if err != nil {
		json.NewEncoder(w).Encode(farmResp(503, err.Error()))
		return
	}

	if resp == nil {
		json.NewEncoder(w).Encode(farmResp(404, "Drink not found"))
		return
	}

	err = server.hasLiquids(resp)
	if err != nil {
		json.NewEncoder(w).Encode(farmResp(409, err.Error()))
		return
	}

	var wg sync.WaitGroup

	server.pourON()
	speech := htgotts.Speech{Folder: "audio", Language: "sv"}
	message := fmt.Sprint(reasons[rand.Intn(len(reasons))])
	go speech.Speak(message)

	go func() {
		for _, liquid := range resp {
			wg.Add(1)
			glog.Infof("Pouring ingredient ID: %d DrinkID: %d LiquidID: %d LiquidName: %s Volume: %d\n", *liquid.ID, *liquid.DrinkID, *liquid.LiquidID, *liquid.LiquidName, *liquid.Volume)
			go func(liquid *drinks.DrinkIngredient) {
				defer wg.Done()
				done := make(chan bool)
				duration := time.Duration(int64(time.Second) * int64(*liquid.Volume))
				job := &pumps.PumpMSG{
					LiquidID: liquid.LiquidID,
					Time:     duration,
					Done:     done,
				}
				server.PM.Queue <- job
				<-done
			}(liquid)
		}
		wg.Wait()
		server.pourOFF()
	}()

	json.NewEncoder(w).Encode(farmResp(200, "Pouring started"))
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
