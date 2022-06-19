package streamprocessings

import (
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (c *controller) getByID(ctx echo.Context) error {
	id := ctx.Param("id")

	fmt.Println("DKDKDKDKDKDKD")

	var report entities.StreamProcessingReport
	result := c.gormDB.First(&report, id)

	if result.Error == nil {
		ctx.JSON(http.StatusOK, report)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.String(http.StatusNotFound, fmt.Sprintf("Stream processing with id %s not found", id))
	}

	return result.Error
}
