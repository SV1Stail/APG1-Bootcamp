package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func divide_file(filepath string, chunk_size int) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var temp_files []string
	var lines []string
	file_number := 0

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) == chunk_size {
			temp_file_name, err := make_temp_file(lines, file_number)
			if err != nil {
				return nil, err
			}
			temp_files = append(temp_files, temp_file_name)
			lines = lines[:0]
			file_number++
		}
	}
	if len(lines) > 0 {
		temp_file_name, err := make_temp_file(lines, file_number)
		if err != nil {
			return nil, err
		}
		temp_files = append(temp_files, temp_file_name)
		lines = nil
	}
	return temp_files, nil
}

func make_temp_file(lines []string, file_number int) (string, error) {
	sort.Strings(lines)
	temp_file_name := "temp_" + strconv.Itoa(file_number) + ".txt"
	temp_file, err := os.Create(temp_file_name)
	if err != nil {
		return "", err
	}
	defer temp_file.Close()
	writer := bufio.NewWriter(temp_file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
	return temp_file_name, nil
}

func merge_files(temp_files []string, out_file_name string) error {
	out_file, err := os.Create(out_file_name)
	if err != nil {
		return err
	}
	defer out_file.Close()

	writers := make([]*bufio.Scanner, len(temp_files))

	for i, file_name := range temp_files {
		file, err := os.Open(file_name)
		if err != nil {
			return err
		}
		defer file.Close()
		writers[i] = bufio.NewScanner(file)
		if !writers[i].Scan() {
			return fmt.Errorf("file %s is empty", file_name)
		}
	}

	writer := bufio.NewWriter(out_file)
	for {
		min_index := -1
		var min_string string
		for i, scanner := range writers {
			if scanner == nil {
				continue
			}
			if scanner.Text() == "" {
				continue
			}
			if min_index == -1 || scanner.Text() < min_string {
				min_index = i
				min_string = scanner.Text()
			}
		}
		if min_index == -1 {
			break
		}
		fmt.Fprintln(writer, min_string)
		if !writers[min_index].Scan() {
			if err := writers[min_index].Err(); err != nil {
				return fmt.Errorf("error reading from file: %v", err)
			}
			writers[min_index] = nil
		}
	}
	writer.Flush()

	return nil

}

func outside__sort(file_path string, output_file_name string) error {
	chunk_size := 3
	temp_files, err := divide_file(file_path, chunk_size)
	if err != nil {
		return fmt.Errorf("error in func divide_file, error: %s", err)
	}
	err = merge_files(temp_files, output_file_name)
	if err != nil {
		return fmt.Errorf("error in func merge_files, error: %s", err)
	}
	defer func() {
		for _, temp_file := range temp_files {
			os.Remove(temp_file)
		}
	}()
	//fmt.Println("Sorting complete. Output written to:", output_file_name)
	return nil
}

func read_path() (string, string) {
	var ss1 = flag.String("old", "", "path to snapshot 1")
	var ss2 = flag.String("new", "", "path to snapshot 2")
	flag.Parse()
	return *ss1, *ss2
}

func check(file_path_1 string, file_path_2 string) error {
	file_1, err := os.Open(file_path_1)
	if err != nil {
		return fmt.Errorf("error in func check, error: %s", err)
	}
	defer file_1.Close()
	file_2, err := os.Open(file_path_2)
	if err != nil {
		return fmt.Errorf("error in func check, error: %s", err)
	}
	defer file_2.Close()

	scanner_1 := bufio.NewScanner(file_1)
	scanner_2 := bufio.NewScanner(file_2)
	scanner_1_flag := scanner_1.Scan()
	scanner_2_flag := scanner_2.Scan()

	for scanner_2_flag && scanner_1_flag {
		line_1 := scanner_1.Text()
		line_2 := scanner_2.Text()
		if line_1 < line_2 {
			fmt.Printf("REMOVED %s\n", line_1)
			scanner_1_flag = scanner_1.Scan()
		} else if line_1 > line_2 {
			fmt.Printf("ADDED %s\n", line_2)
			scanner_2_flag = scanner_2.Scan()
		} else {
			scanner_1_flag = scanner_1.Scan()
			scanner_2_flag = scanner_2.Scan()
		}
	}
	for scanner_1_flag {
		fmt.Printf("REMOVED %s\n", scanner_1.Text())
		scanner_1_flag = scanner_1.Scan()
	}
	for scanner_2_flag {
		fmt.Printf("ADDED+ %s\n", scanner_2.Text())
		scanner_2_flag = scanner_2.Scan()

	}
	return nil
}

func main() {
	path_ss1, path_ss2 := read_path()
	sorted_ss1 := "path_ss1_sorted.txt"
	sorted_ss2 := "path_ss2_sorted.txt"

	if path_ss1 == "" {
		fmt.Println("cant find snapshot 1")
		os.Exit(1)
	} else if path_ss2 == "" {
		fmt.Println("cant find snapshot 2")
		os.Exit(1)
	}
	err := outside__sort(path_ss1, sorted_ss1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = outside__sort(path_ss2, sorted_ss2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = check(sorted_ss1, sorted_ss2)
	os.Remove(sorted_ss1)
	os.Remove(sorted_ss2)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
