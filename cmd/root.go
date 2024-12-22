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
package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/repr"
	"github.com/spf13/cobra"
	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
)

func completion(cmd *cobra.Command, mode string) error {
	switch mode {
	case "bash":
		return cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return errors.New("invalid shell parameter! Valid ones: bash|zsh|fish|powershell")
	}
}

func Execute() {
	var (
		conf           cfg.Config
		configfile     string
		ShowVersion    bool
		ShowCompletion string
	)

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	SearchConfigs := []string{
		filepath.Join(home, ".config", "anydb", "anydb.toml"),
		filepath.Join(home, ".anydb.toml"),
		"anydb.toml",
	}

	var rootCmd = &cobra.Command{
		Use:   "anydb <command> [options]",
		Short: "anydb",
		Long:  `A personal key value store`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			db, err := app.New(conf.Dbfile, conf.Dbbucket, conf.Debug)
			if err != nil {
				return err
			}

			conf.DB = db

			var configs []string
			if configfile != "" {
				configs = []string{configfile}
			} else {
				configs = SearchConfigs
			}

			if err := conf.GetConfig(configs); err != nil {
				return err
			}

			if conf.Debug {
				repr.Println(conf)
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if ShowVersion {
				fmt.Printf("This is anydb version %s\n", cfg.Version)
				return nil
			}

			if len(ShowCompletion) > 0 {
				return completion(cmd, ShowCompletion)
			}

			if len(args) == 0 {
				return errors.New("no command specified")
			}

			return nil
		},
	}

	// options
	rootCmd.PersistentFlags().BoolVarP(&ShowVersion, "version", "v", false, "Print program version")
	rootCmd.PersistentFlags().BoolVarP(&conf.Debug, "debug", "d", false, "Enable debugging")
	rootCmd.PersistentFlags().StringVarP(&conf.Dbfile, "dbfile", "f",
		filepath.Join(home, ".config", "anydb", "default.db"), "DB file to use")
	rootCmd.PersistentFlags().StringVarP(&conf.Dbbucket, "bucket", "b",
		app.BucketData, "use other bucket (default: "+app.BucketData+")")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", "", "toml config file")

	rootCmd.AddCommand(Set(&conf))
	rootCmd.AddCommand(List(&conf))
	rootCmd.AddCommand(Get(&conf))
	rootCmd.AddCommand(Del(&conf))
	rootCmd.AddCommand(Export(&conf))
	rootCmd.AddCommand(Import(&conf))
	rootCmd.AddCommand(Serve(&conf))
	rootCmd.AddCommand(Man(&conf))
	rootCmd.AddCommand(Info(&conf))

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
