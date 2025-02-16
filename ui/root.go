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
package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/tlinden/anydb/app"
	"github.com/tlinden/anydb/cfg"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type Loader struct {
	items []list.Item
	conf  *cfg.Config
}

func (loader *Loader) Update() error {
	entries, err := loader.conf.DB.List(&app.DbAttr{}, loader.conf.Fulltext)
	if err != nil {
		return err
	}

	loader.items = nil

	for _, entry := range entries {
		loader.items = append(loader.items, item{
			title:       entry.Key,
			description: entry.Preview,
		})
	}

	return nil
}

const (
	ModeDefault = iota
	ModeView
)

type model struct {
	conf         *cfg.Config
	loader       *Loader
	quitting     bool
	err          error
	list         list.Model
	keys         *listKeyMap
	delegateKeys *delegateKeyMap
	mode         int    // mode
	selected     string // current key to be deleted, viewed or edited
	viewport     viewport.Model
}

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
	viewItem         key.Binding
	quit             key.Binding
	close            key.Binding
}

type item struct {
	title       string
	description string
}

type ChoiceMsg string

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
		viewItem: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "view"),
		),
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl-c", "esc"),
			key.WithHelp("q", "exit"),
		),
		close: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "close pager"),
		),
	}
}

func NewModel(config *cfg.Config, entries app.DbEntries) model {
	var (
		delegateKeys = newDelegateKeyMap()
		listKeys     = newListKeyMap()
		loader       = Loader{conf: config}
	)

	// Setup list
	if err := loader.Update(); err != nil {
		panic(err)
	}

	delegate := newItemDelegate(delegateKeys, config)
	dbList := list.New(loader.items, delegate, 0, 0)
	dbList.Title = "DB Entries"
	dbList.Styles.Title = titleStyle

	dbList.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.toggleSpinner,
			listKeys.insertItem,
			listKeys.toggleTitleBar,
			listKeys.toggleStatusBar,
			listKeys.togglePagination,
			listKeys.toggleHelpMenu,
		}
	}

	return model{
		list:         dbList,
		viewport:     viewport.New(20, 40),
		keys:         listKeys,
		delegateKeys: delegateKeys,
	}
}

func (m model) Init() tea.Cmd {
	m.mode = ModeDefault
	return nil
}

// Main update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgt := msg.(type) {
	case tea.WindowSizeMsg:
		// pre-calculate pager size
		verticalMarginHeight, _ := m.pagerMargin()
		m.viewport.Width = msgt.Width
		m.viewport.Height = msgt.Height - verticalMarginHeight
	}

	// Hand off to subs
	switch m.mode {
	case ModeDefault:
		return m.UpdateList(msg)
	case ModeView:
		return m.UpdatePager(msg)
	}

	return nil, nil
}

func (m model) View() string {
	if m.quitting {
		return "\n  See you later!\n\n"
	}

	// Hand off to subs
	switch m.mode {
	case ModeDefault:
		return appStyle.Render(m.list.View())
	case ModeView:
		return m.ViewPager()
	}

	return ""
}
