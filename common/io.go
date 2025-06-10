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
package common

import "os"

func CleanError(file string, err error) error {
	// remove given [backup] file and forward the given error
	return os.Remove(file)
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)

	if err != nil {
		// return false on any error
		return false
	}

	return !info.IsDir()
}
