package dbreader

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
)

type DBReader interface {
	DBRead(str string) (Recipes, error)
	Converter_output(db Recipes) error
}

type Ingredient struct {
	IngredientName  string  `json:"ingredient_name" xml:"itemname"`
	IngredientCount string  `json:"ingredient_count" xml:"itemcount"`
	IngredientUnit  *string `json:"ingredient_unit,omitempty" xml:"itemunit,omitempty" `
}

type Cake struct {
	Name        string       `json:"name" xml:"name"`
	Time        string       `json:"time" xml:"stovetime"`
	Ingredients []Ingredient `json:"ingredients" xml:"ingredients"`
}

type Recipes struct {
	Cakes []Cake `json:"cake" xml:"cake"`
}

//
//
//
//
//
//
//
type xml_Recipes struct {
	Cakes []xml_Cake ` xml:"cake"`
}
type xml_Cake struct {
	Name        string           ` xml:"name"`
	Time        string           ` xml:"stovetime"`
	Ingredients []xml_Ingredient `xml:"ingredients"`
}
type xml_Ingredient struct {
	Items []xml_IngredientItem `xml:"item"`
}
type xml_IngredientItem struct {
	Name  string  ` xml:"itemname"`
	Count string  ` xml:"itemcount"`
	Unit  *string ` xml:"itemunit,omitempty"`
}

//
//
//
//
//
//
//
type DB_json_reader struct{}

func (jr DB_json_reader) Converter_output(db Recipes) error {
	output, err := xml.MarshalIndent(db, "", "    ")
	if err != nil {
		return errors.New("cant mkae output json")
	}
	fmt.Println(string(output))

	return nil
}

func (jr DB_json_reader) DBRead(path string) (Recipes, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Recipes{}, err
	}
	var db Recipes
	err = json.Unmarshal(file, &db)
	if err != nil {
		return Recipes{}, err
	}
	return db, nil
}

type DB_xml_reader struct{}

func (xr DB_xml_reader) Converter_output(db Recipes) error {
	output, err := json.MarshalIndent(db, "", "    ")
	if err != nil {
		return errors.New("cant make json output")
	}
	fmt.Println(string(output))
	return nil
}

func (xr DB_xml_reader) DBRead(path string) (Recipes, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return Recipes{}, err
	}
	var db_xml xml_Recipes
	err = xml.Unmarshal(file, &db_xml)
	if err != nil {
		return Recipes{}, err
	}
	var db_json Recipes
	var empty_Ingredient Ingredient
	empty_Ingredient.IngredientName = ""
	empty_Ingredient.IngredientCount = ""
	empty_Ingredient.IngredientUnit = nil
	var empty_Cake Cake
	empty_Cake.Name = ""
	empty_Cake.Time = ""
	for i := range db_xml.Cakes {
		db_json.Cakes = append(db_json.Cakes, empty_Cake)
		db_json.Cakes[i].Name = db_xml.Cakes[i].Name
		db_json.Cakes[i].Time = db_xml.Cakes[i].Time
		for j := range db_xml.Cakes[i].Ingredients {
			for n := range db_xml.Cakes[i].Ingredients[j].Items {
				db_json.Cakes[i].Ingredients = append(db_json.Cakes[i].Ingredients, empty_Ingredient)
				if len(db_xml.Cakes[i].Ingredients[j].Items) > 0 {
					db_json.Cakes[i].Ingredients[n].IngredientName = db_xml.Cakes[i].Ingredients[j].Items[n].Name
					db_json.Cakes[i].Ingredients[n].IngredientCount = db_xml.Cakes[i].Ingredients[j].Items[n].Count
					db_json.Cakes[i].Ingredients[n].IngredientUnit = db_xml.Cakes[i].Ingredients[j].Items[n].Unit

				}
			}
		}
	}

	return db_json, nil
}
