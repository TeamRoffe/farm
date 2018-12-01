package server

import (
	"fmt"

	"github.com/teamroffe/farm/pkg/drinks"
)

func (server *FarmServer) getingredients(drinkID int) (*[]drinks.DrinkIngredient, error) {
	var ingredients []drinks.DrinkIngredient
	results, err := server.DB.Query("select drink_ingredients.id as ingredient_id, liquids.liquid_name, drink_ingredients.liquid_id, drink_ingredients.volume from drinks left join drink_ingredients on drinks.id = drink_ingredients.drink_id left join liquids on liquids.id = drink_ingredients.liquid_id where drinks.id = ?;", drinkID)
	if err != nil {
		return &ingredients, err
	}
	defer results.Close()
	for results.Next() {
		var ingredient drinks.DrinkIngredient
		err = results.Scan(&ingredient.ID, &ingredient.LiquidName, &ingredient.LiquidID, &ingredient.Volume)
		if err != nil {
			return &ingredients, err
		}
		ingredients = append(ingredients, ingredient)
	}
	return &ingredients, nil
}

func (server *FarmServer) getDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		server.Config.Section("database").Key("username").String(),
		server.Config.Section("database").Key("password").String(),
		server.Config.Section("database").Key("hostname").String(),
		server.Config.Section("database").Key("port").String(),
		server.Config.Section("database").Key("database").String(),
	)
}
