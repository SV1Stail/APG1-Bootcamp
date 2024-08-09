package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// work
func Mean(mas []int) float64 {
	var mean float64 = 0

	for i := 0; i < len(mas); i++ {
		mean += float64(mas[i])
	}
	mean = mean / float64(len(mas))
	return mean
}

// work
func Mode(mas []int) int {
	var noda int = 0
	var noda_kol_vo int = 0
	if len(mas) > 0 {
		noda = int(mas[0])
		noda_kol_vo = 1
	}
	var tmp_map map[int]int = make(map[int]int)

	for i := 0; i < len(mas); i++ {
		tmp_map[mas[i]]++
		if noda_kol_vo < tmp_map[mas[i]] {
			noda = mas[i]
			noda_kol_vo = tmp_map[mas[i]]
		} else if noda > mas[i] && noda_kol_vo == tmp_map[mas[i]] {
			noda = mas[i]
		}
	}
	return noda
}

// work
func Mediana(mas []int) float64 {
	var x float64 = 0
	if !sort.IntsAreSorted(mas) {
		sort.Ints(mas)
	}
	if len(mas)%2 != 0 {
		x = float64(mas[len(mas)/2])
	} else {
		x = float64(mas[len(mas)/2]+mas[len(mas)/2-1]) / 2
	}

	return x
}

// work
func Sd(mas []int) float64 {

	var sr_znach float64 = Mean(mas)
	var all_sum float64 = 0
	for i := 0; i < len(mas); i++ {
		all_sum = all_sum + math.Pow(float64(mas[i])-sr_znach, 2)
	}
	sr_znach = all_sum / sr_znach
	return math.Sqrt(sr_znach)
}

// work
func Scaner() ([]int, error) {
	var slice []int = make([]int, 0)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите числа, каждое c новой строки (нажмите Ctrl+D для завершения) или напишите \"end\":")
	for scanner.Scan() {
		input := scanner.Text()
		input = strings.TrimSpace(input)
		if input == "end" {
			break
		}
		if x, err := strconv.Atoi(input); err != nil {
			return slice, errors.New("Ошибка: введенное значение не является целым числом.\n")
		} else {
			if x < -100000 || x > 100000 {
				return slice, errors.New("Ошибка: введенное значение вне диапазона допустимых значений.\n")
			}
			slice = append(slice, x)
		}
		if err := scanner.Err(); err != nil {
			return slice, errors.New("Ошибка чтения:")
		}
	}
	if len(slice) < 1 {
		return slice, errors.New("Введено 0 символов")
	}
	return slice, nil
}
func main() {
	var slice []int
	slice, err := Scaner()
	var (
		mn  = flag.Bool("mn", false, "flag to show Mean")
		md  = flag.Bool("md", false, "flag to show Mode")
		mdn = flag.Bool("mdn", false, "flag to show Mediana")
		sd  = flag.Bool("sd", false, "flag to show SD")
	)
	flag.Parse()
	if err == nil {
		if *mn {
			fmt.Printf("Mean: %.2f\n", Mean(slice))
		}
		if *mdn {
			fmt.Printf("Median: %.2f\n", Mediana(slice))
		}
		if *md {
			fmt.Println("Mode:", Mode(slice))
		}
		if *sd {
			fmt.Printf("SD: %.2f\n", Sd(slice))
		}
	} else {
		fmt.Println(err)
	}

}
