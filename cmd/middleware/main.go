package main

import (
	"net/http"

	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers/notifications"
	"fmt"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	notificationsGroup := e.Group("/notifications")
	notifications.Bind(notificationsGroup)

	fmt.Println(notificationsGroup)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
