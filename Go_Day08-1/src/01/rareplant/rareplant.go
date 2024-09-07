package rareplant

import "fmt"

type UnknownPlant struct {
	FlowerType string
	LeafType   string
	Color      int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
	FlowerColor int
	LeafType    string
	Height      int `unit:"inches"`
}

type Plants interface {
	printPlant()
}

func (p UnknownPlant) printPlant() {
	fmt.Printf("\nFlowerType:%s\nLeafType:%s\nColor:%d\n", p.FlowerType, p.LeafType, p.Color)
}
func (p AnotherUnknownPlant) printPlant() {
	fmt.Printf("\nFlowerColor:%d\nLeafType:%s\nHeight:%d\n", p.FlowerColor, p.LeafType, p.Height)
}
func DescribePlant[T Plants](plant T) {
	plant.printPlant()
}

func DescribePlantSwitch(x interface{}) {
	switch v := x.(type) {
	case UnknownPlant:
		fmt.Printf("\nFlowerType:%s\nLeafType:%s\nColor:%d\n", v.FlowerType, v.LeafType, v.Color)
	case AnotherUnknownPlant:
		fmt.Printf("\nFlowerColor:%d\nLeafType:%s\nHeight:%d\n", v.FlowerColor, v.LeafType, v.Height)
	default:
		fmt.Println("Unknown type")

	}
}
