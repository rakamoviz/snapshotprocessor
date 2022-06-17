package notifications

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func receive(c echo.Context) error {
	pCluster := &models.Cluster{Code: "clust_1"}
	return c.JSON(http.StatusOK, pCluster)
}
