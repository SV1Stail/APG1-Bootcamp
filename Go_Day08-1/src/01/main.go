package main

import "01/rareplant"

func main() {
	plant1 := rareplant.UnknownPlant{
		FlowerType: "Rose",
		LeafType:   "Broad",
		Color:      16711680, // Красный в формате RGB
	}

	plant2 := rareplant.AnotherUnknownPlant{
		FlowerColor: 255, // Синий
		LeafType:    "Narrow",
		Height:      12,
	}
	rareplant.DescribePlant(plant1)
	rareplant.DescribePlant(plant2)
}
