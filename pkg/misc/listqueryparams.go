package misc

import (
	//"errors"

	//"github.com/rakamoviz/snapshotprocessor/pkg/entities"
	"fmt"

	"encoding/json"
	"strings"

	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
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

func ListQueryParamsToQuery(ctx context.Context, gormDB *gorm.DB, queryParams ListQueryParams) *gorm.DB {
	query := gormDB.Where(queryParams.Filter)

	if len(queryParams.Sort) == 2 {
		query = query.Order(strings.Join(queryParams.Sort, " "))
	}
	if len(queryParams.Range) == 2 {
		query = query.Offset(queryParams.Range[0]).Limit(queryParams.Range[1] - queryParams.Range[0])
	}

	return query
}
