package drinks

//Drink info
type Drink struct {
	ID          *int    `json:"id"`
	Name        *string `json:"drink_name"`
	Category    *int    `json:"category,omitempty"`
	Description *string `json:"description,omitempty"`
	URL         *string `json:"url,omitempty"`
}

//Liquid info
type Liquid struct {
	ID   *int    `json:"id"`
	Name *string `json:"name"`
	URL  *string `json:"url,omitempty"`
}

//DrinkIngredient info
type DrinkIngredient struct {
	ID         *int    `json:"id"`
	DrinkID    *int    `json:"drink_id,omitempty"`
	LiquidName *string `json:"liquid_name"`
	LiquidID   *int    `json:"liquid_id"`
	Volume     *int    `json:"volume"`
}

//Category info
type Category struct {
	ID   *int    `json:"id"`
	Name *string `json:"name"`
}
