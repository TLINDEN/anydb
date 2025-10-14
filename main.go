/*
Copyright © 2024 Thomas von Dein

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/inconshreveable/mousetrap"
	"github.com/tlinden/anydb/cmd"
)

func main() {
	const NoLogsLevel = 100
	slog.SetLogLoggerLevel(NoLogsLevel)

	Main()
}

func init() {
	// if we're running on Windows  AND if the user double clicked the
	// exe  file from explorer, we  tell them and then  wait until any
	// key has been hit, which  will make the cmd window disappear and
	// thus give the user time to read it.
	if runtime.GOOS == "windows" {
		if mousetrap.StartedByExplorer() {
			fmt.Println("Do no double click anydb.exe!")
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

func Main() int {
	cmd.Execute()
	return 0
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
