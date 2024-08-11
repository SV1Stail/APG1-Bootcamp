package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR!!! too few arguments")
		os.Exit(1)
	}
	fmt.Println(len(os.Args))
	command := os.Args[1]
	command_args := os.Args[2:]

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command_args = append(command_args, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("ERROR!!! when scanning")
		os.Exit(1)
	}

	command_exec := exec.Command(command, command_args...)
	command_exec.Stdout = os.Stdout
	command_exec.Stderr = os.Stderr

	if err := command_exec.Run(); err != nil {
		fmt.Println("ERROR!!! when starting command_exec")
		os.Exit(1)
	}
}

//./myFind -f -ext 'log' /home/svs/work/Go_Day02-1/materials/ | ./myXargs ./myWc -l
