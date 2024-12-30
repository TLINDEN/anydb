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
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
	"github.com/tlinden/anydb/output"
	"github.com/tlinden/anydb/rest"
)

func Export(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
	)

	var cmd = &cobra.Command{
		Use:   "export -o <json filename>",
		Short: "Export database to json file",
		Long:  `Export database to json file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			conf.Mode = "json"

			entries, err := conf.DB.Getall(&attr)
			if err != nil {
				return err
			}

			return output.WriteJSON(&attr, conf, entries)
		},
	}

	cmd.PersistentFlags().StringVarP(&attr.File, "output-file", "o", "", "filename or - for STDIN")
	if err := cmd.MarkPersistentFlagRequired("output-file"); err != nil {
		panic(err)
	}

	cmd.Aliases = append(cmd.Aliases, "dump")
	cmd.Aliases = append(cmd.Aliases, "backup")

	return cmd
}

func Import(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
	)

	var cmd = &cobra.Command{
		Use:   "import -i <json file>",
		Short: "Import database dump",
		Long:  `Import database dump`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			out, err := conf.DB.Import(&attr)
			if err != nil {
				return err
			}

			fmt.Print(out)
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&attr.File, "import-file", "i", "", "filename or - for STDIN")
	cmd.PersistentFlags().StringArrayVarP(&attr.Tags, "tags", "t", nil, "tags, multiple allowed")
	if err := cmd.MarkPersistentFlagRequired("import-file"); err != nil {
		panic(err)
	}

	cmd.Aliases = append(cmd.Aliases, "restore")

	return cmd
}

func Help(conf *cfg.Config) *cobra.Command {
	return nil
}
func Man(conf *cfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "man",
		Short: "show manual page",
		Long:  `show manual page`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			man := exec.Command("less", "-")

			var b bytes.Buffer

			b.WriteString(manpage)

			man.Stdout = os.Stdout
			man.Stdin = &b
			man.Stderr = os.Stderr

			err := man.Run()

			if err != nil {
				return fmt.Errorf("failed to execute 'less': %w", err)
			}

			return nil
		},
	}

	return cmd
}

func Serve(conf *cfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "serve [-l host:port]",
		Short: "run REST API listener",
		Long:  `run REST API listener`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			return rest.Runserver(conf, nil)
		},
	}

	cmd.PersistentFlags().StringVarP(&conf.Listen, "listen", "l", "localhost:8787", "host:port")

	return cmd
}

func Info(conf *cfg.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "info",
		Short: "info",
		Long:  `show info about database`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true

			info, err := conf.DB.Info()
			if err != nil {
				return err
			}

			return output.Info(os.Stdout, conf, info)
		},
	}

	cmd.PersistentFlags().BoolVarP(&conf.NoHumanize, "no-human", "N", false, "do not translate to human readable values")

	return cmd
}

func Edit(conf *cfg.Config) *cobra.Command {
	var (
		attr app.DbAttr
	)

	var cmd = &cobra.Command{
		Use:   "edit <key>",
		Short: "Edit a key",
		Long:  `Edit a key`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("no key specified")
			}

			// errors at this stage do not cause the usage to be shown
			cmd.SilenceUsage = true
			password := []byte{}

			if len(args) > 0 {
				attr.Key = args[0]
			}

			// fetch entry
			entry, err := conf.DB.Get(&attr)
			if err != nil {
				return err
			}

			if len(entry.Value) == 0 && entry.Binary {
				return errors.New("key contains binary uneditable content")
			}

			// decrypt if needed
			if entry.Encrypted {
				pass, err := getPassword()
				if err != nil {
					return err
				}
				password = pass

				clear, err := app.Decrypt(pass, entry.Value)
				if err != nil {
					return err
				}

				entry.Value = clear
				entry.Encrypted = false
			}

			// determine editor, vi is default
			editor := getEditor()

			// save file to a temp file, call the editor with it, read
			// it  back in and  compare the content with  the original
			// one
			newcontent, err := editContent(editor, string(entry.Value))
			if err != nil {
				return err
			}

			// all is valid, fill our DB feeder
			newattr := app.DbAttr{
				Key:       attr.Key,
				Tags:      attr.Tags,
				Encrypted: attr.Encrypted,
				Val:       []byte(newcontent),
			}

			// encrypt if needed
			if conf.Encrypt {
				err = app.Encrypt(password, &attr)
				if err != nil {
					return err
				}
			}

			// done
			return conf.DB.Set(&newattr)
		},
	}

	cmd.Aliases = append(cmd.Aliases, "modify")
	cmd.Aliases = append(cmd.Aliases, "mod")
	cmd.Aliases = append(cmd.Aliases, "ed")
	cmd.Aliases = append(cmd.Aliases, "vi")

	return cmd
}

func getEditor() string {
	editor := "vi"

	enveditor, present := os.LookupEnv("EDITOR")
	if present {
		if editor != "" {
			editor = enveditor
		}
	}

	return editor
}

// taken from github.com/tlinden/rpn/ (my own program)
func editContent(editor string, content string) (string, error) {
	// create a temp file
	tmp, err := os.CreateTemp("", "stack")
	if err != nil {
		return "", fmt.Errorf("failed to create templ file: %w", err)
	}
	defer os.Remove(tmp.Name())

	// put the content into a tmp file
	_, err = tmp.WriteString(content)
	if err != nil {
		return "", fmt.Errorf("failed to write value to temp file: %w", err)
	}

	// execute editor with our tmp file containing current stack
	cmd := exec.Command(editor, tmp.Name())

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run editor command %s: %w", editor, err)
	}

	// read the file back in
	modified, err := os.Open(tmp.Name())
	if err != nil {
		return "", fmt.Errorf("failed to open temp file: %w", err)
	}
	defer modified.Close()

	newcontent, err := io.ReadAll(modified)
	if err != nil {
		return "", fmt.Errorf("failed to read from temp file: %w", err)
	}

	newcontentstr := string(newcontent)
	if content == newcontentstr {
		return "", fmt.Errorf("content not modified, aborting")
	}

	return newcontentstr, nil
}
