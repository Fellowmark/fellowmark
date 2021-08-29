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
	TotalRows  int64       `json:"totalRows"`
	TotalPages int         `json:"totalPages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

func (p *Pagination) SetTotalPages() {
	if p.Limit == 0 {
		p.TotalPages = 0
	} else {
		p.TotalPages = int(math.Ceil(float64(p.TotalRows) / float64(p.Limit)))
	}
}

func (p *Pagination) SetTotalRows(totalRows int64) {
	p.TotalRows = totalRows
	if p.Limit <= 0 || p.Limit > int(p.TotalRows) {
		p.Limit = int(p.TotalRows)
	}

	// p.TotalPages = int(math.Ceil(float64(p.TotalRows) / float64(p.Limit)))
	p.SetTotalPages()

	if p.Page <= 0 || p.Page > int(p.TotalPages) {
		p.Page = int(p.TotalPages)
	}
}

func GetPagination(r *http.Request) Pagination {
	urlVars := r.URL.Query()
	limit, _ := strconv.Atoi(urlVars.Get("limit"))
	page, _ := strconv.Atoi(urlVars.Get("page"))
	sort := urlVars.Get("sort")

	if sort == "" {
		sort = "id asc"
	}

	return Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}
}

func Paginate(db *gorm.DB, scope func(db *gorm.DB) *gorm.DB, r *http.Request, pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Scopes(scope).Count(&totalRows)
	pagination.SetTotalRows(totalRows)
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.Limit).Order(pagination.Sort)
	}
}
