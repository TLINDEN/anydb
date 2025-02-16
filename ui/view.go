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
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) UpdatePager(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		verticalMarginHeight, headerHeight := m.pagerMargin()

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
		m.viewport.YPosition = headerHeight

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.close):
			m.mode = ModeDefault
		case key.Matches(msg, m.keys.quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) ViewPager() string {
	return fmt.Sprintf("%s\n%s\n%s", m.pagerHeaderView(),
		m.viewport.View(), m.pagerFooterView())
}

func (m model) pagerHeaderView() string {
	title := titleStyle.Render("View entry")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) pagerFooterView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m model) pagerMargin() (int, int) {
	headerHeight := lipgloss.Height(m.pagerHeaderView())
	footerHeight := lipgloss.Height(m.pagerFooterView())
	return headerHeight + footerHeight, headerHeight
}
