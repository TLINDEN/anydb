package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
		ShowVersion    bool
		ShowCompletion string
	)

	var rootCmd = &cobra.Command{
		Use:   "anydb <command> [options]",
		Short: "anydb",
		Long:  `A personal key value store`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			db, err := app.New(conf.Dbfile, conf.Debug)
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

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// options
	rootCmd.PersistentFlags().BoolVarP(&ShowVersion, "version", "v", false, "Print program version")
	rootCmd.PersistentFlags().BoolVarP(&conf.Debug, "debug", "d", false, "Enable debugging")
	rootCmd.PersistentFlags().StringVarP(&conf.Dbfile, "dbfile", "f",
		filepath.Join(home, ".config", "anydb", "default.db"), "DB file to use")

	rootCmd.AddCommand(Set(&conf))
	rootCmd.AddCommand(List(&conf))
	rootCmd.AddCommand(Get(&conf))
	rootCmd.AddCommand(Del(&conf))
	rootCmd.AddCommand(Export(&conf))
	rootCmd.AddCommand(Import(&conf))
	rootCmd.AddCommand(Man(&conf))

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
