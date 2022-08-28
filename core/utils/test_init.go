package utils

import (
	"github.com/golang-jwt/jwt"
	"github.com/nus-utils/nus-peer-review/models"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

func constructToken() *jwt.Token {
	return &jwt.Token{
		Raw:       "",
		Method:    nil,
		Header:    nil,
		Claims:    nil,
		Signature: "",
		Valid:     true,
	}
}

func constructHTTPRequest() *http.Request {
	return &http.Request{
		Method:           "",
		URL:              nil,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           http.Header{},
		Body:             nil,
		GetBody:          nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Host:             "",
		Form:             nil,
		PostForm:         nil,
		MultipartForm:    nil,
		Trailer:          nil,
		RemoteAddr:       "",
		RequestURI:       "",
		TLS:              nil,
		Cancel:           nil,
		Response:         nil,
	}
}

func constructUser() models.User {
	return models.User{
		Model: models.Model{
			ID:        0,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
		},
		Email:    "",
		Name:     "",
		Password: "",
	}
}

func constructGormDB() *gorm.DB {
	return &gorm.DB{
		Config:       nil,
		Error:        nil,
		RowsAffected: 0,
		Statement: &gorm.Statement{
			DB:                   nil,
			TableExpr:            nil,
			Table:                "",
			Model:                nil,
			Unscoped:             false,
			Dest:                 nil,
			ReflectValue:         reflect.Value{},
			Clauses:              nil,
			BuildClauses:         nil,
			Distinct:             false,
			Selects:              nil,
			Omits:                nil,
			Joins:                nil,
			Preloads:             nil,
			Settings:             sync.Map{},
			ConnPool:             nil,
			Schema:               nil,
			Context:              nil,
			RaiseErrorOnNotFound: false,
			SkipHooks:            false,
			SQL:                  strings.Builder{},
			Vars:                 nil,
			CurDestIndex:         0,
		},
	}
}

func constructAssignment() models.Assignment {
	return models.Assignment{
		Model:     models.Model{},
		Name:      "",
		Module:    models.Module{},
		ModuleID:  0,
		GroupSize: 1,
		Deadline:  0,
	}
}

func constructParing() models.Pairing {
	return models.Pairing{
		Model:        models.Model{},
		Assignment:   models.Assignment{},
		AssignmentID: 0,
		Student:      models.Student{},
		StudentID:    0,
		Marker:       models.Student{},
		MarkerID:     0,
		Active:       false,
	}
}

func constructStudent(name string) models.Student {
	return models.Student{
		Model:    models.Model{},
		Email:    "",
		Name:     name,
		Password: "",
	}
}
