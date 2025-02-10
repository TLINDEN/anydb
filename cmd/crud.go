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

			// check if value given as file or via stdin and fill attr accordingly
			if err := attr.ParseKV(); err != nil {
				return err
			}

			// encrypt?
			if conf.Encrypt {
				pass, err := getPassword()
				if err != nil {
					return err
				}

				err = app.Encrypt(pass, &attr)
				if err != nil {
					return err
				}
			}

			return conf.DB.Set(&attr)
		},
	}

	cmd.PersistentFlags().BoolVarP(&conf.Encrypt, "encrypt", "e", false, "encrypt value")
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
		Use:   "get  <key> [-o <file>] [-m <mode>] [-n -N] [-T <tpl>]",
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

			if entry.Encrypted {
				pass, err := getPassword()
				if err != nil {
					return err
				}

				clear, err := app.Decrypt(pass, []byte(entry.Value))
				if err != nil {
					return err
				}

				entry.Value = string(clear)
				entry.Size = uint64(len(entry.Value))
				entry.Encrypted = false
			}

			return output.Print(os.Stdout, conf, &attr, entry)
		},
	}

	cmd.PersistentFlags().StringVarP(&attr.File, "output", "o", "", "output value to file (ignores -m)")
	cmd.PersistentFlags().StringVarP(&conf.Mode, "mode", "m", "", "output format (simple|wide|json|template) (default 'simple')")
	cmd.PersistentFlags().BoolVarP(&conf.NoHeaders, "no-headers", "n", false, "omit headers in tables")
	cmd.PersistentFlags().BoolVarP(&conf.NoHumanize, "no-human", "N", false, "do not translate to human readable values")
	cmd.PersistentFlags().StringVarP(&conf.Template, "template", "T", "", "go template for '-m template'")

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

	cmd.Aliases = append(cmd.Aliases, "d")
	cmd.Aliases = append(cmd.Aliases, "rm")

	return cmd
}

func List(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
		wide bool
	)

	var cmd = &cobra.Command{
		Use:   "list  [<filter-regex> | -t <tag> ] [-m <mode>] [-nNis] [-T <tpl>]",
		Short: "List database contents",
		Long:  `List database contents`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			if len(args) > 0 {
				if conf.CaseInsensitive {
					attr.Args = []string{"(?i)" + args[0]}
				} else {
					attr.Args = args
				}
			}

			// turn comma list into slice, if needed
			if len(attr.Tags) == 1 && strings.Contains(attr.Tags[0], ",") {
				attr.Tags = strings.Split(attr.Tags[0], ",")
			}

			if wide {
				conf.Mode = "wide"
			}

			entries, err := conf.DB.List(&attr, conf.Fulltext)
			if err != nil {
				return err
			}

			return output.List(os.Stdout, conf, entries)
		},
	}

	cmd.PersistentFlags().StringVarP(&conf.Mode, "mode", "m", "", "output format (table|wide|json|template), wide is a verbose table. (default 'table')")
	cmd.PersistentFlags().StringVarP(&conf.Template, "template", "T", "", "go template for '-m template'")
	cmd.PersistentFlags().BoolVarP(&wide, "wide-output", "l", false, "output mode: wide")
	cmd.PersistentFlags().BoolVarP(&conf.NoHeaders, "no-headers", "n", false, "omit headers in tables")
	cmd.PersistentFlags().BoolVarP(&conf.NoHumanize, "no-human", "N", false, "do not translate to human readable values")
	cmd.PersistentFlags().BoolVarP(&conf.CaseInsensitive, "case-insensitive", "i", false, "filter case insensitive")
	cmd.PersistentFlags().BoolVarP(&conf.Fulltext, "search-fulltext", "s", false, "perform a full text search")
	cmd.PersistentFlags().StringArrayVarP(&attr.Tags, "tags", "t", nil, "tags, multiple allowed")

	cmd.Aliases = append(cmd.Aliases, "ls")
	cmd.Aliases = append(cmd.Aliases, "/")
	cmd.Aliases = append(cmd.Aliases, "find")
	cmd.Aliases = append(cmd.Aliases, "search")

	return cmd
}

func getPassword() ([]byte, error) {
	var pass []byte

	envpass := os.Getenv("ANYDB_PASSWORD")

	if envpass == "" {
		readpass, err := app.AskForPassword()
		if err != nil {
			return nil, err
		}

		pass = readpass
	} else {
		pass = []byte(envpass)
	}

	return pass, nil
}
