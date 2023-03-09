package provider

import (
	"gorm.io/gorm"
	"math"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func isSqlError(err error) bool {
	if err == nil {
		return false
	}
	if err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

type Pagination[T any] struct {
	Limit      int    `json:"limit,omitempty;query:limit"`
	Page       int    `json:"page,omitempty;query:page"`
	Sort       string `json:"sort,omitempty;query:sort"`
	TotalRows  int64  `json:"totalRows"`
	TotalPages int    `json:"totalPages"`
	Rows       []T    `json:"rows"`
	Search     string `json:"search"`
}

func (p *Pagination[any]) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination[any]) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 20
	}
	if p.Limit > 200 {
		p.Limit = 20
	}
	return p.Limit
}

func (p *Pagination[any]) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination[any]) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at asc"
	}
	return p.toSnakeCase(p.Sort)
}

func (p *Pagination[any]) toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func paginate[T any](value T, searchField string, pagination *Pagination[T], db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Where(searchField+" like ?", "%"+pagination.Search+"%").Count(&totalRows)
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.GetLimit())))

	return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Where(searchField+" like ?", "%"+pagination.Search+"%").Order(pagination.GetSort())
}
