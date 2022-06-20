package misc

import (
	//"errors"

	//"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"fmt"

	"github.com/labstack/echo/v4"
	//"gorm.io/gorm"
	"encoding/json"
	"strings"
)

type ListQueryParams struct {
	Sort   []string          `json:"sort" query:"sort"`
	Range  []int             `json:"range" query:"range"`
	Filter map[string]string `json:"filter" query:"filter"`
}

func BindListQueryParams(ctx echo.Context) (ListQueryParams, error) {
	filterParam := ctx.QueryParam("filter")
	if strings.Trim(filterParam, " ") == "" {
		filterParam = "{}"
	}

	rangeParam := ctx.QueryParam("range")
	if strings.Trim(rangeParam, " ") == "" {
		rangeParam = "[]"
	}

	sortParam := ctx.QueryParam("sort")
	if strings.Trim(sortParam, " ") == "" {
		sortParam = "[]"
	}

	s := fmt.Sprintf(`{"filter":%s,"range":%s,"sort":%s}`, filterParam, rangeParam, sortParam)
	var listQueryParams ListQueryParams
	return listQueryParams, json.Unmarshal([]byte(s), &listQueryParams)
}
