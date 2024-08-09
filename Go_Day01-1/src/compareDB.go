package main

import (
	"day1/dbreader"
	"flag"
	"fmt"
	"os"
	"strings"
)

func Compare_dbs(path_xml string, path_json string) error {
	var xml_db dbreader.Recipes
	var json_db dbreader.Recipes
	var err error
	if strings.HasSuffix(path_xml, "xml") {
		reader := dbreader.DB_xml_reader{}
		if xml_db, err = reader.DBRead(path_xml); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("it's not xml or cant open: %s", path_xml)
	}

	if strings.HasSuffix(path_json, "json") {
		reader := dbreader.DB_json_reader{}
		if json_db, err = reader.DBRead(path_json); err != nil {
			return fmt.Errorf("it's not json or cant open: %s", path_json)
		}
	}
	deleted_cakes, added_cakes, same_cakes := deleted_added_same_cakes_map(xml_db, json_db)
	for cake_name, _ := range added_cakes {
		fmt.Printf("ADDED cake \"%s\"\n", cake_name)
	}
	for cake_name, _ := range deleted_cakes {
		fmt.Printf("REMOVED cake \"%s\"\n", cake_name)
	}

	check_cooking_time(same_cakes, json_db)
	compare_ingr(same_cakes, json_db)

	return nil
}

func compare_units(ingr, json_ingr dbreader.Ingredient, cake_name string) {
	if ingr.IngredientCount != json_ingr.IngredientCount {
		fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n", ingr.IngredientName, cake_name, json_ingr.IngredientCount, ingr.IngredientCount)
	}

	if *ingr.IngredientUnit == "" && json_ingr.IngredientUnit != nil {
		fmt.Printf("ADDED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", *json_ingr.IngredientUnit, ingr.IngredientName, cake_name)
	} else if *ingr.IngredientUnit == "" && json_ingr.IngredientUnit == nil {
		fmt.Printf("THERE WAS NO unit for ingredient \"%s\" for cake \"%s\"\n", ingr.IngredientName, cake_name)
	} else if *ingr.IngredientUnit != "" && json_ingr.IngredientUnit == nil {
		fmt.Printf("REMOVED unit \"%s\" for ingredient \"%s\" for cake \"%s\"\n", *ingr.IngredientUnit, ingr.IngredientName, cake_name)
	} else if json_ingr.IngredientUnit != nil {
		if *ingr.IngredientUnit != *json_ingr.IngredientUnit {
			fmt.Printf("CHANGED unit for ingredient \"%s\" for cake  \"%s\" - \"%s\" instead of \"%s\"\n", ingr.IngredientName, cake_name, *json_ingr.IngredientUnit, *ingr.IngredientUnit)
		}

	}

}
func compare_ingr(same_cakes map[string]dbreader.Cake, json_db dbreader.Recipes) {
	// same_ingredients := make(map[string]dbreader.Ingredient)
	json_cakes := make(map[string]dbreader.Cake)
	for _, cake := range json_db.Cakes {
		json_cakes[cake.Name] = cake
	}
	same_ingr := make([]string, 0)
	flag := false
	for _, cake := range same_cakes {
		if json_cake, ok := json_cakes[cake.Name]; ok {
			for _, ingr := range cake.Ingredients {
				for _, json_ingr := range json_cake.Ingredients {
					if json_ingr.IngredientName == ingr.IngredientName {
						same_ingr = append(same_ingr, ingr.IngredientName)
						flag = true
						compare_units(ingr, json_ingr, cake.Name)
						break
					}
				}
				if !flag {
					fmt.Printf("REMOVED ingredient \"%s\" for cake  \"%s\"\n", ingr.IngredientName, cake.Name)
				} else {
					flag = false
				}
			}
			for _, json_ingr := range json_cake.Ingredients {
				flag = false
				for _, ingr_name := range same_ingr {
					if json_ingr.IngredientName == ingr_name {
						flag = true
						break
					}
				}
				if !flag {
					fmt.Printf("ADDED ingredient \"%s\" for cake  \"%s\"\n", json_ingr.IngredientName, cake.Name)
					flag = false
				}
			}
		}
	}

}
func check_cooking_time(same_cakes map[string]dbreader.Cake, json_db dbreader.Recipes) {
	json_cakes := make(map[string]dbreader.Cake)
	for _, cake := range json_db.Cakes {
		json_cakes[cake.Name] = cake
	}
	for _, cake := range same_cakes {
		if cake_stolen, ok := json_cakes[cake.Name]; ok {
			if cake.Time != cake_stolen.Time {
				fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", cake.Name, cake_stolen.Time, cake.Time)
			}
		}
	}
}

// same_cakes содержит торт из оригинальной баз
func deleted_added_same_cakes_map(xml_db, json_db dbreader.Recipes) (map[string]dbreader.Cake, map[string]dbreader.Cake, map[string]dbreader.Cake) {
	deleted_cakes := make(map[string]dbreader.Cake)
	added_cakes := make(map[string]dbreader.Cake)
	same_cakes := make(map[string]dbreader.Cake)
	for _, cake := range xml_db.Cakes {
		deleted_cakes[cake.Name] = cake
	}
	for _, cake := range json_db.Cakes {
		added_cakes[cake.Name] = cake
	}

	for cake_name, cake := range deleted_cakes {
		if _, ok := added_cakes[cake_name]; ok {
			delete(added_cakes, cake_name)
			delete(deleted_cakes, cake_name)
			same_cakes[cake_name] = cake
		}
	}
	return deleted_cakes, added_cakes, same_cakes
}

func read_path_xml() (string, string) {
	var xml_path = flag.String("old", "", "read path xml")
	var json_path = flag.String("new", "", "read path json")
	flag.Parse()
	return *xml_path, *json_path
}

func main() {
	path_xml, path_json := read_path_xml()
	if path_xml == "" {
		fmt.Println("cant find xml ")
		os.Exit(1)
	} else if path_json == "" {
		fmt.Println("cant find  json")
		os.Exit(1)
	}
	if err := Compare_dbs(path_xml, path_json); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
