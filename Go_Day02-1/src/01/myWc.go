package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"unicode"
)

type _flags struct {
	l           bool
	w           bool
	m           bool
	files_names []string
}

func files_parse() (_flags, error) {
	var flags _flags
	flag.BoolVar(&flags.l, "l", false, "for counting lines")
	flag.BoolVar(&flags.w, "w", false, "for counting words")
	flag.BoolVar(&flags.m, "m", false, "for counting characters")

	flag.Parse()
	counter := 0
	if flags.l {
		counter++
	}
	if flags.w {
		counter++
	}
	if flags.m {
		counter++
	}
	if counter > 1 {
		return _flags{}, fmt.Errorf("to mach flags")
	} else if counter == 0 {
		flags.w = true
	}
	flags.files_names = flag.Args()
	if len(flags.files_names) == 0 {
		return _flags{}, fmt.Errorf("no paths to files")
	}
	return flags, nil
}

func count_lines(ctx context.Context, file_name string, wg *sync.WaitGroup, err_chan chan<- error) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	file, err := os.Open(file_name)
	if err != nil {
		err_chan <- err
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	l_counter := 0
	for scanner.Scan() {
		select {
		case <-ctx.Done():
		default:
			l_counter++
		}
	}
	if err := scanner.Err(); err != nil {
		err_chan <- err
		return
	}
	fmt.Printf("%d\t%s\n", l_counter, file_name)
}

func count_l(files_names []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, len(files_names))

	for _, file := range files_names {
		wg.Add(1)
		go count_lines(ctx, file, &wg, errChan)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		if err != nil {
			cancel()
			return err
		}
	}

	return nil
}
func count_words(ctx context.Context, file_name string, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	file, err := os.Open(file_name)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	w_counter := 0
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			w_counter++
		}
	}
	if err := scanner.Err(); err != nil {
		errChan <- err
		return
	}
	fmt.Printf("%d\t%s\n", w_counter, file_name)

}

func count_w(files_names []string) error {
	cntx, cancel_cntx := context.WithCancel(context.Background())
	defer cancel_cntx()

	var wg sync.WaitGroup
	errChan := make(chan error, len(files_names))

	for _, file := range files_names {
		wg.Add(1)
		go count_words(cntx, file, &wg, errChan)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()
	for err := range errChan {
		if err != nil {
			cancel_cntx()
			return err
		}
	}

	return nil
}

func count_members(cntx context.Context, file_name string, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	select {
	case <-cntx.Done():
		return
	default:
	}

	file, err := os.Open(file_name)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	l_counter := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, l := range line {
			select {
			case <-cntx.Done():
				return
			default:
				if unicode.IsLetter(l) || unicode.IsDigit(l) {
					l_counter++
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		errChan <- err
		return
	}
	fmt.Printf("%d\t%s\n", l_counter, file_name)

}

func count_m(files_names []string) error {
	cntx, cancel_cntx := context.WithCancel(context.Background())
	defer cancel_cntx()

	var wg sync.WaitGroup
	errChan := make(chan error, len(files_names))

	for _, file := range files_names {
		wg.Add(1)
		go count_members(cntx, file, &wg, errChan)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			cancel_cntx()
			return err
		}
	}

	return nil
}

func main() {
	var flags _flags
	flags, err := files_parse()
	if err != nil {
		fmt.Printf("ERROR!!! func \"files_parse\": %v\n", err)
		os.Exit(1)
	}
	if flags.l {
		if err := count_l(flags.files_names); err != nil {
			fmt.Printf("ERROR!!! func \"count_l\": %v\n", err)
			os.Exit(1)
		}
	} else if flags.m {
		if err := count_m(flags.files_names); err != nil {
			fmt.Printf("ERROR!!! func \"count_m\": %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := count_w(flags.files_names); err != nil {
			fmt.Printf("ERROR!!! func \"count_w\": %v\n", err)
			os.Exit(1)
		}
	}

}
