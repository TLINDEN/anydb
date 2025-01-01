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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/tlinden/anydb/cfg"
)

// used to return to the api client
type Result struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Runserver(conf *cfg.Config, args []string) error {
	// setup api server
	router := SetupServer(conf)

	// public rest api routes
	api := router.Group("/anydb/v1")
	{
		api.Get("/", func(c *fiber.Ctx) error {
			return RestList(c, conf)
		})

		api.Post("/", func(c *fiber.Ctx) error {
			// same thing as above but allows to supply parameters, see app.Dbattr{}
			return RestList(c, conf)
		})

		api.Get("/:key", func(c *fiber.Ctx) error {
			return RestGet(c, conf)
		})

		api.Delete("/:key", func(c *fiber.Ctx) error {
			return RestDelete(c, conf)
		})

		api.Put("/", func(c *fiber.Ctx) error {
			return RestSet(c, conf)
		})
	}

	// public routes
	{
		router.Get("/", func(c *fiber.Ctx) error {
			return c.Send([]byte("Use the REST API"))
		})
	}

	return router.Listen(conf.Listen)
}

func SetupServer(conf *cfg.Config) *fiber.App {
	// disable colors
	fiber.DefaultColors = fiber.Colors{}

	router := fiber.New(fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		Immutable:     true,
		ServerHeader:  "anydb serve",
		AppName:       "anydb",
	})

	router.Use(logger.New(logger.Config{
		Format:        "${pid} ${ip}:${port} ${status} - ${method} ${path}\n",
		DisableColors: true,
	}))

	router.Use(cors.New(cors.Config{
		AllowMethods:  "GET,PUT,POST,DELETE",
		ExposeHeaders: "Content-Type,Accept",
	}))

	router.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	return router
}

/*
Wrapper to respond with proper json status, message and code,
shall be prepared and called by the handlers directly.
*/
func JsonStatus(c *fiber.Ctx, code int, msg string) error {
	success := true

	if code != fiber.StatusOK {
		success = false
	}

	return c.Status(code).JSON(Result{
		Code:    code,
		Message: msg,
		Success: success,
	})
}
