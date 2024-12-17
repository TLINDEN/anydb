package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
	"github.com/tlinden/anydb/output"
)

func Set(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
	)

	var cmd = &cobra.Command{
		Use:   "set <key> [<value> | -r <file>] [-t <tag>]",
		Short: "Insert key/value pair",
		Long:  `Insert key/value pair`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No key/value pair specified")
			}

			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			if len(args) > 0 {
				attr.Args = args
			}

			if err := conf.DB.Set(&attr); err != nil {
				return err
			}

			return conf.DB.Close()
		},
	}

	cmd.PersistentFlags().StringVarP(&attr.File, "file", "r", "", "Filename or - for STDIN")
	cmd.PersistentFlags().StringArrayVarP(&attr.Tags, "tags", "t", nil, "tags, multiple allowed")

	return cmd
}

func Get(conf *cfg.Config) *cobra.Command {
	return nil
}

func Del(conf *cfg.Config) *cobra.Command {
	return nil
}

func List(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
	)

	var cmd = &cobra.Command{
		Use:   "list [-t <tag>] [-o <mode>] [<filter-regex>]",
		Short: "List database contents",
		Long:  `List database contents`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			if len(args) > 0 {
				attr.Args = args
			}

			entries, err := conf.DB.List(&attr)
			if err != nil {
				return err
			}

			output.List(os.Stdout, conf, entries)

			return conf.DB.Close()
		},
	}

	cmd.PersistentFlags().StringVarP(&conf.Mode, "output-mode", "o", "", "output mode: wide, yaml, json, table")
	cmd.PersistentFlags().StringArrayVarP(&attr.Tags, "tags", "t", nil, "tags, multiple allowed")

	return cmd
}

func Find(conf *cfg.Config) *cobra.Command {
	return nil
}

func Help(conf *cfg.Config) *cobra.Command {
	return nil
}

func Man(conf *cfg.Config) *cobra.Command {
	return nil
}
