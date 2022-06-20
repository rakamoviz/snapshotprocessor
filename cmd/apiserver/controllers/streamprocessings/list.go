package streamprocessings

import (
	//"errors"
	"net/http"

	//"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/misc"
	"github.com/labstack/echo/v4"
	//"gorm.io/gorm"
)

func (c *controller) list(ctx echo.Context) error {
	listQueryParams, err := misc.BindListQueryParams(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Bad Request")
	}

	return ctx.JSON(http.StatusOK, listQueryParams)
}
