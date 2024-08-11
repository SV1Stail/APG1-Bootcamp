package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type _flags struct {
	a           string
	files_names []string
}

func path_parse(flags *_flags) {

	flag.StringVar(&flags.a, "a", "", "where save archives")
	flag.Parse()
	flags.files_names = flag.Args()
}

func archive(a string, file_path string, cntx context.Context, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	select {
	case <-cntx.Done():
		return
	default:
	}

	file_inf0, err := os.Stat(file_path)
	if err != nil {
		errChan <- fmt.Errorf("ERROR!!! func: \"archive\". Can't take info - %v", err)
		return
	}
	timestamp := file_inf0.ModTime().Unix()
	// fmt.Println(file_inf0.Name())
	var arch_name string
	if a == "" {
		treimmed_path := strings.TrimSuffix(file_path, ".log")
		arch_name = fmt.Sprintf("%s%d.tar.gz", treimmed_path, timestamp)
	} else {
		treimmed_path := strings.TrimSuffix(file_inf0.Name(), ".log")
		arch_name = fmt.Sprintf("%s%s%d.tar.gz", a, treimmed_path, timestamp)

	}
	// fmt.Println(arch_name)
	out_file, err := os.Create(arch_name)
	if err != nil {
		errChan <- fmt.Errorf("ERROR!!! func: \"archive\". Can't create file - %v", err)
		return
	}
	defer out_file.Close()

	gzip_writer := gzip.NewWriter(out_file)
	defer gzip_writer.Close()
	tar_writer := tar.NewWriter(gzip_writer)
	defer tar_writer.Close()
	file, err := os.Open(file_path)
	if err != nil {
		errChan <- fmt.Errorf("ERROR!!! func: \"archive\". Can't open .log file %s - %v", file_path, err)
		return
	}
	defer file.Close()

	header := &tar.Header{
		Name:    filepath.Base(file_path),
		Mode:    int64(file_inf0.Mode()),
		Size:    file_inf0.Size(),
		ModTime: file_inf0.ModTime(),
		Uname:   "hello go", // Ваше кастомное сообщение
	}

	if err := tar_writer.WriteHeader(header); err != nil {
		errChan <- fmt.Errorf("ERROR!!! func: \"archive\". Can't write header  - %v", err)
		return
	}
	if _, err := io.Copy(tar_writer, file); err != nil {
		errChan <- err
		return
	}

}

func archives_sync(a string, files_paths []string) error {
	cntx, cancel_cntx := context.WithCancel(context.Background())
	defer cancel_cntx()

	var wg sync.WaitGroup
	errChan := make(chan error, len(files_paths))

	for _, file := range files_paths {
		file_info, err := os.Stat(file)
		if os.IsNotExist(err) {
			fmt.Printf("file %s doesn't exist\n", file)
			continue
		} else if file_info.IsDir() {
			fmt.Printf("%s is a directory\n", file)
			continue
		}
		wg.Add(1)
		go archive(a, file, cntx, &wg, errChan)

	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			cancel_cntx()
			return fmt.Errorf("ERROR!!! func archives_sync - %v", err)
		}
	}

	return nil
}

func main() {
	var flags _flags
	path_parse(&flags)
	if len(flags.files_names) == 0 {
		fmt.Println("Pls input files")
		os.Exit(1)
	}
	if err := archives_sync(flags.a, flags.files_names); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
