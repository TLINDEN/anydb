package cmd

import (
	"errors"
	"os"
	"strings"

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
				return errors.New("no key/value pair specified")
			}

			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			if len(args) > 0 {
				attr.Args = args
			}

			// turn comma list into slice, if needed
			if len(attr.Tags) == 1 && strings.Contains(attr.Tags[0], ",") {
				attr.Tags = strings.Split(attr.Tags[0], ",")
			}

			return conf.DB.Set(&attr)
		},
	}

	cmd.PersistentFlags().StringVarP(&attr.File, "file", "r", "", "Filename or - for STDIN")
	cmd.PersistentFlags().StringArrayVarP(&attr.Tags, "tags", "t", nil, "tags, multiple allowed")

	cmd.Aliases = append(cmd.Aliases, "add")
	cmd.Aliases = append(cmd.Aliases, "s")
	cmd.Aliases = append(cmd.Aliases, "+")

	return cmd
}

func Get(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
	)

	var cmd = &cobra.Command{
		Use:   "get <key> [-o <file>]",
		Short: "Retrieve value for a key",
		Long:  `Retrieve value for a key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("no key specified")
			}

			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			if len(args) > 0 {
				attr.Key = args[0]
			}

			entry, err := conf.DB.Get(&attr)
			if err != nil {
				return err
			}

			return output.Print(os.Stdout, conf, entry)
		},
	}

	cmd.PersistentFlags().StringVarP(&conf.Mode, "output", "o", "", "output to file")

	cmd.Aliases = append(cmd.Aliases, "show")
	cmd.Aliases = append(cmd.Aliases, "g")
	cmd.Aliases = append(cmd.Aliases, ".")

	return cmd
}

func Del(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
	)

	var cmd = &cobra.Command{
		Use:   "del <key>",
		Short: "Delete key",
		Long:  `Delete key and value matching key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("No key specified")
			}

			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			if len(args) > 0 {
				attr.Key = args[0]
			}

			return conf.DB.Del(&attr)
		},
	}

	cmd.PersistentFlags().StringVarP(&conf.Mode, "output", "o", "", "output to file")

	cmd.Aliases = append(cmd.Aliases, "d")
	cmd.Aliases = append(cmd.Aliases, "rm")

	return cmd
}

func Export(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
	)

	var cmd = &cobra.Command{
		Use:   "export [<json filename>]",
		Short: "Export database to json",
		Long:  `Export database to json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			if len(args) == 0 {
				attr.File = "-"
			} else {
				attr.File = args[0]
			}

			conf.Mode = "json"

			entries, err := conf.DB.List(&attr)
			if err != nil {
				return err
			}

			return output.WriteFile(&attr, conf, entries)
		},
	}

	cmd.Aliases = append(cmd.Aliases, "dump")
	cmd.Aliases = append(cmd.Aliases, "backup")

	return cmd
}

func List(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
		wide bool
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

			// turn comma list into slice, if needed
			if len(attr.Tags) == 1 && strings.Contains(attr.Tags[0], ",") {
				attr.Tags = strings.Split(attr.Tags[0], ",")
			}

			if wide {
				conf.Mode = "wide"
			}

			entries, err := conf.DB.List(&attr)
			if err != nil {
				return err
			}

			return output.List(os.Stdout, conf, entries)
		},
	}

	cmd.PersistentFlags().StringVarP(&conf.Mode, "output-mode", "o", "", "output format (table|wide|json), wide is a verbose table. (default 'table')")
	cmd.PersistentFlags().BoolVarP(&wide, "wide-output", "l", false, "output mode: wide")
	cmd.PersistentFlags().BoolVarP(&conf.NoHeaders, "no-headers", "n", false, "omit headers in tables")
	cmd.PersistentFlags().StringArrayVarP(&attr.Tags, "tags", "t", nil, "tags, multiple allowed")

	cmd.Aliases = append(cmd.Aliases, "/")
	cmd.Aliases = append(cmd.Aliases, "ls")

	return cmd
}

func Import(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
	)

	var cmd = &cobra.Command{
		Use:   "import [<json file>]",
		Short: "Import database dump",
		Long:  `Import database dump`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			if len(args) == 0 {
				attr.File = "-"
			} else {
				attr.File = args[0]
			}

			return conf.DB.Import(&attr)
		},
	}

	cmd.PersistentFlags().StringVarP(&attr.File, "file", "r", "", "Filename or - for STDIN")
	cmd.PersistentFlags().StringArrayVarP(&attr.Tags, "tags", "t", nil, "tags, multiple allowed")

	cmd.Aliases = append(cmd.Aliases, "add")
	cmd.Aliases = append(cmd.Aliases, "s")
	cmd.Aliases = append(cmd.Aliases, "+")

	return cmd
}

func Help(conf *cfg.Config) *cobra.Command {
	return nil
}

func Man(conf *cfg.Config) *cobra.Command {
	return nil
}
