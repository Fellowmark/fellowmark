package utils

import (
	"math"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,omitempty;query:page"`
	Sort       string      `json:"sort,omitempty;query:sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Page > p.TotalPages {
		p.Page = p.TotalPages
	}
	if p.Page < 1 {
		p.Page = 1
	}
	return p.Page
}

func (p Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "id desc"
	}
	return p.Sort
}

func GetPagination(r *http.Request) Pagination {
	urlVars := r.URL.Query()
	limit, _ := strconv.Atoi(urlVars.Get("limit"))
	page, _ := strconv.Atoi(urlVars.Get("page"))
	pagination := Pagination{
		Limit: limit,
		Page:  page,
		Sort:  urlVars.Get("sort"),
	}
	return pagination
}

func Paginate(tx *gorm.DB, r *http.Request, pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	tx.Count(&pagination.TotalRows)
	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.GetLimit())))
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
