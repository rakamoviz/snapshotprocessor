package middlewares

import (
	"log"
	"net/http"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/services/auth"
	"github.com/labstack/echo/v4"
)

type ApiKeyCheck struct {
	authService auth.Service
}

func NewApiKeyCheck(authService auth.Service) *ApiKeyCheck {
	return &ApiKeyCheck{authService: authService}
}

func (m *ApiKeyCheck) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		apiClient, ok, err := m.authService.ValidateApiKey(
			ctx.Request().Context(), ctx.Request().Header.Get("X-API-Key"),
		)
		if err != nil {
			log.Println(err)
			return ctx.String(http.StatusInternalServerError, "Server error")
		}

		if !ok {
			return ctx.String(http.StatusUnauthorized, "Unauthorized")
		}

		ctx.Set("ApiClient", apiClient)

		if err := next(ctx); err != nil {
			ctx.Error(err)
		}

		return err
	}
}
