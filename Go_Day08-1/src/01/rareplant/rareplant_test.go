package rareplant_test

import (
	"01/rareplant"
	"testing"
)

func BenchmarkDescribePlant(b *testing.B) {
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
	for i := 0; i < b.N; i++ {

		rareplant.DescribePlant(plant1)
		rareplant.DescribePlant(plant2)
	}
}

// по тестам DescribePlantSwitch быстрее и использует меньше памяти чем DescribePlant
func BenchmarkDescribePlantSwitch(b *testing.B) {
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
	for i := 0; i < b.N; i++ {

		rareplant.DescribePlantSwitch(plant1)
		rareplant.DescribePlantSwitch(plant2)
	}
}
