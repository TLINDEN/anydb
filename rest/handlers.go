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
package rest

import (
	//"github.com/alecthomas/repr"

	"github.com/gofiber/fiber/v2"
	"codeberg.org/scip/anydb/app"
	"codeberg.org/scip/anydb/cfg"
)

type SetContext struct {
	Query string `json:"query" form:"query"`
}

type ListResponse struct {
	Success bool
	Code    int
	Entries app.DbEntries
}

type SingleResponse struct {
	Success bool
	Code    int
	Entry   *app.DbEntry
}

func RestList(c *fiber.Ctx, conf *cfg.Config) error {
	attr := new(app.DbAttr)

	if len(c.Body()) > 0 {
		if err := c.BodyParser(attr); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"errors": err.Error(),
			})

		}
	}

	// get list
	entries, err := conf.DB.List(attr, attr.Fulltext)
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"Unable to list keys: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(
		ListResponse{
			Success: true,
			Code:    fiber.StatusOK,
			Entries: entries,
		},
	)
}

func RestGet(c *fiber.Ctx, conf *cfg.Config) error {
	if c.Params("key") == "" {
		return JsonStatus(c, fiber.StatusForbidden,
			"key not provided")
	}

	// get list
	entry, err := conf.DB.Get(&app.DbAttr{Key: c.Params("key")})
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"Unable to get key: "+err.Error())
	}
	if entry.Key == "" {
		return JsonStatus(c, fiber.StatusForbidden,
			"Key does not exist")
	}

	return c.Status(fiber.StatusOK).JSON(
		SingleResponse{
			Success: true,
			Code:    fiber.StatusOK,
			Entry:   entry,
		},
	)
}

func RestDelete(c *fiber.Ctx, conf *cfg.Config) error {
	if c.Params("key") == "" {
		return JsonStatus(c, fiber.StatusForbidden,
			"key not provided")
	}

	// get list
	err := conf.DB.Del(&app.DbAttr{Key: c.Params("key")})
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"Unable to delete key: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(
		Result{
			Success: true,
			Code:    fiber.StatusOK,
			Message: "key deleted",
		},
	)
}

func RestSet(c *fiber.Ctx, conf *cfg.Config) error {
	attr := new(app.DbAttr)
	if err := c.BodyParser(attr); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": err.Error(),
		})

	}

	err := conf.DB.Set(attr)
	if err != nil {
		return JsonStatus(c, fiber.StatusForbidden,
			"Unable to set key: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(
		Result{
			Success: true,
			Code:    fiber.StatusOK,
		},
	)
}
