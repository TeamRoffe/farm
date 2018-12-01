ritchlista
select drinks.id, drinks.drink_name, liquids.liquid_name, drink_ingredients.liquid_id, drink_ingredients.volume from drinks left join drink_ingredients on drinks.id =
                 drink_ingredients.drink_id left join liquids on liquids.id = drink_ingredients.liquid_id;


id, drink_name, liquid_name, liquid_id


från ingredients: id, drink_id, liquid_id, joina in liquid_name, volume


-----
select drink_ingredients.id as ingredient_id, drinks.id as drink_id, liquids.liquid_name, drink_ingredients.liquid_id, drink_ingredients.volume from drinks left join drink_ingredients on drinks.id = drink_ingredients.drink_id left join liquids on liquids.id = drink_ingredients.liquid_id where drinks.id = ?;