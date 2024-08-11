package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type parsed_input struct {
	f        bool
	sl       bool
	d        bool
	ext      string
	filePath string
}

func read_path() (parsed_input, error) {
	var input parsed_input
	flag.BoolVar(&input.d, "d", false, "find symbol links")
	flag.BoolVar(&input.sl, "sl", false, "find symbolic links")
	flag.BoolVar(&input.f, "f", false, "find files")
	flag.StringVar(&input.ext, "ext", "", "suffix, only with -f")

	flag.Parse()
	args := flag.Args()

	if input.ext != "" && !input.f {
		return parsed_input{}, fmt.Errorf("-ext without -f")
	}

	if len(args) < 1 {
		return parsed_input{}, fmt.Errorf("no path")
	}
	input.filePath = args[0]
	return input, nil
}

func find_d(path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			if _, err := os.Open(path); err != nil {
				if os.IsPermission(err) {
					fmt.Printf("Permission denied dir: %s\n", path)
					return filepath.SkipDir
				}
				return err
			}
			fmt.Println(path)
		}

		return nil

	})
	if err != nil {
		return fmt.Errorf("Error for -d walking the path: %v", err)
	}
	return nil
}

func find_f(path string, ext string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("Permission denied file: %s\n", path)
				return nil
			}
			return err
		}
		if info.Mode().IsRegular() {
			if ext == "" || filepath.Ext(info.Name()) == "."+ext {
				fmt.Println(path)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error for -f walking the path: %v", err)
	}
	return nil
}

func find_sl(path string) error {
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("Permission denied s-link: %s\n", path)
				return nil
			}
			return err
		}

		if info.Mode()&os.ModeSymlink != 0 {
			target, err := os.Readlink(path)
			if err != nil {
				fmt.Printf("%s -> [broken]\n", path)
				return nil
			}
			if !filepath.IsAbs(target) {
				target = filepath.Join(filepath.Dir(path), target)
			}
			if _, err = os.Stat(target); os.IsNotExist(err) {
				fmt.Printf("%s -> [broken]\n", path)
				return nil
			}

			fmt.Printf("%s -> %s\n", path, target)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error for -sl walking the path: %v", err)
	}
	return nil
}

func main() {
	var input parsed_input
	input, err := read_path()
	if err != nil {
		fmt.Printf("error %v", err)
		os.Exit(1)
	}
	if input.d {
		if err := find_d(input.filePath); err != nil {
			fmt.Printf("error %v", err)
			os.Exit(1)
		}
	}
	if input.f {
		if err := find_f(input.filePath, input.ext); err != nil {
			fmt.Printf("error %v", err)
			os.Exit(1)
		}
	}
	if input.sl {
		if err := find_sl(input.filePath); err != nil {
			fmt.Printf("error %v", err)
			os.Exit(1)
		}
	}

}
