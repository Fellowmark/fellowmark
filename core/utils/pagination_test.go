package utils

import (
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestPagination_SetTotalPages(t *testing.T) {
	inputs := []*Pagination{
		{
			Limit:      0,
			Page:       0,
			Sort:       "",
			TotalRows:  0,
			TotalPages: 10,
			Rows:       nil,
		},
		{
			Limit:      10,
			Page:       0,
			Sort:       "",
			TotalRows:  100,
			TotalPages: 0,
			Rows:       nil,
		},
	}
	expects := []*Pagination{
		{
			Limit:      0,
			Page:       0,
			Sort:       "",
			TotalRows:  0,
			TotalPages: 0,
			Rows:       nil,
		},
		{
			Limit:      10,
			Page:       0,
			Sort:       "",
			TotalRows:  100,
			TotalPages: 10,
			Rows:       nil,
		},
	}
	for i, input := range inputs {
		input.SetTotalPages()
		assert.Equal(t, reflect.DeepEqual(*input, *expects[i]), true)
	}
}

func TestPagination_SetTotalRows(t *testing.T) {
	input := &Pagination{
		Limit:      0,
		Page:       0,
		Sort:       "",
		TotalRows:  0,
		TotalPages: 10,
		Rows:       nil,
	}
	expect := &Pagination{
		Limit:      5,
		Page:       10,
		Sort:       "",
		TotalRows:  5,
		TotalPages: 10,
		Rows:       nil,
	}
	_ = gomonkey.ApplyMethod(reflect.TypeOf(input), "SetTotalPages", func(p *Pagination) {})
	input.SetTotalRows(5)
	assert.Equal(t, reflect.DeepEqual(*input, *expect), true)
}

func TestGetPagination(t *testing.T) {
	r := http.Request{
		URL: &url.URL{},
	}
	urlVals := url.Values{
		"limit": []string{"1"},
		"page":  []string{"2"},
		"sort":  []string{},
	}
	_ = gomonkey.ApplyMethod(reflect.TypeOf(r.URL), "Query", func(url *url.URL) url.Values {
		return urlVals
	})
	pagination := Pagination{
		Limit: 1,
		Page:  2,
		Sort:  "id asc",
	}
	res := GetPagination(&r)
	assert.Equal(t, reflect.DeepEqual(res, pagination), true)
}
