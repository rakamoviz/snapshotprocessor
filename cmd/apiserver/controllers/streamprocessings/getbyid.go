package streamprocessings

import (
	"errors"
	"fmt"
	"net/http"

	"log"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (c *controller) getByID(ctx echo.Context) error {
	id := ctx.Param("id")

	var report entities.StreamProcessingReport
	result := c.gormDB.First(&report, id)

	if result.Error == nil {
		return ctx.JSON(http.StatusOK, report)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("Stream processing with id %s not found", id))
	}

	log.Println(result.Error)
	return ctx.String(http.StatusInternalServerError, "Server error")
}
