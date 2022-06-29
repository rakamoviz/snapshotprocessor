package streamprocessings

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/rakamoviz/snapshotprocessor/pkg/entities/reads"
	"gorm.io/gorm"
)

func (c *controller) getByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Incorrect format of id query parameter")
	}

	report, err := c.streamProcessingReportRepository.ExecuteOne(
		ctx.Request().Context(),
		reads.StreamProcessingReportById(uint(id)),
	)
	if err == nil {
		return ctx.JSON(http.StatusOK, report)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("Stream processing with id %d not found", id))
	}

	log.Println(err)
	return ctx.String(http.StatusInternalServerError, "Server error")
}
