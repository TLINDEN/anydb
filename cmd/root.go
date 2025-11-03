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
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/spf13/cobra"
	"codeberg.org/scip/anydb/app"
	"codeberg.org/scip/anydb/cfg"
	"github.com/tlinden/yadu"
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
				buildInfo, _ := debug.ReadBuildInfo()
				opts := &yadu.Options{
					Level:     slog.LevelDebug,
					AddSource: true,
				}

				slog.SetLogLoggerLevel(slog.LevelDebug)

				handler := yadu.NewHandler(os.Stdout, opts)
				debuglogger := slog.New(handler).With(
					slog.Group("program_info",
						slog.Int("pid", os.Getpid()),
						slog.String("go_version", buildInfo.GoVersion),
					),
				)
				slog.SetDefault(debuglogger)

				slog.Debug("parsed config", "conf", conf)
			}

			dbfile := app.GetDbFile(conf.Dbfile)

			db, err := app.New(dbfile, conf.Dbbucket, conf.Debug)
			if err != nil {
				return err
			}

			conf.DB = db
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
		"", "DB file to use (default: ~/.config/anydb/default.db)")
	rootCmd.PersistentFlags().StringVarP(&conf.Dbbucket, "bucket", "b",
		app.BucketData, "use other bucket (default: "+app.BucketData+")")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", "", "toml config file")

	// CRUD
	rootCmd.AddCommand(Set(&conf))
	rootCmd.AddCommand(List(&conf))
	rootCmd.AddCommand(Get(&conf))
	rootCmd.AddCommand(Del(&conf))

	// backup
	rootCmd.AddCommand(Export(&conf))
	rootCmd.AddCommand(Import(&conf))

	// REST API
	rootCmd.AddCommand(Serve(&conf))

	// auxiliary
	rootCmd.AddCommand(Man(&conf))
	rootCmd.AddCommand(Info(&conf))
	rootCmd.AddCommand(Edit(&conf))

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
