package drinks

//Drink info
type Drink struct {
	ID          int    `json:"id"`
	Name        string `json:"drink_name"`
	Description string `json:"description"`
	URL         string `json:"url,omitempty"`
}

//Liquid info
type Liquid struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

//DrinkIngredient info
type DrinkIngredient struct {
	ID       int `json:"id"`
	DrinkID  int `json:"drink_id"`
	LiquidID int `json:"liquid_id"`
	Volume   int `json:"volume"`
}
