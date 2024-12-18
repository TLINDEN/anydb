package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"

	"github.com/inconshreveable/mousetrap"
	"github.com/tlinden/anydb/cmd"
)

func main() {
	cmd.Execute()
}

func init() {
	// if we're running on Windows  AND if the user double clicked the
	// exe  file from explorer, we  tell them and then  wait until any
	// key has been hit, which  will make the cmd window disappear and
	// thus give the user time to read it.
	if runtime.GOOS == "windows" {
		if mousetrap.StartedByExplorer() {
			fmt.Println("Please do no double click anydb.exe!")
			fmt.Println("Please open a command shell and run it from there.")
			fmt.Println()
			fmt.Print("Press any key to quit: ")
			_, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				panic(err)
			}
		}
	}
}
