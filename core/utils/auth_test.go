package utils

import (
	"errors"
	"github.com/agiledragon/gomonkey"
	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"
)

type TestGenerateJWTScenario struct {
	description string
	user        models.User
	tokenString string
	tokenErr    error
	expect      string
	noErr       bool
}

func TestGenerateJWT(t *testing.T) {

	token := constructToken()
	var logger *log.Logger
	scenarios := []TestGenerateJWTScenario{
		{
			description: "get token string err",
			user:        models.User{},
			tokenString: "",
			tokenErr:    errors.New("get token string err"),
			expect:      "",
			noErr:       false,
		},
		{
			description: "normal",
			user:        models.User{},
			tokenString: "token",
			tokenErr:    nil,
			expect:      "token",
			noErr:       true,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyFunc(os.Getenv, func(key string) string {
				return "env"
			})
			patches.ApplyFunc(jwt.NewWithClaims, func(method jwt.SigningMethod, claims jwt.Claims) *jwt.Token {
				return token
			})
			patches.ApplyMethod(reflect.TypeOf(token), "SignedString", func(tk *jwt.Token, key interface{}) (string, error) {
				return scenario.tokenString, scenario.tokenErr
			})
			patches.ApplyMethod(reflect.TypeOf(logger), "Println", func(logger *log.Logger, v ...interface{}) {})
			defer patches.Reset()
			res, err := GenerateJWT(scenario.user)
			assert.Equal(t, res == scenario.expect, true)
			assert.Equal(t, err == nil, scenario.noErr)
		})
	}
}

type TestParseJWTScenario struct {
	description string
	tokenString string
	token       *jwt.Token
	tokenErr    error
	expect      *ClaimsData
	noErr       bool
}

func TestParseJWT(t *testing.T) {
	successfulToken := constructToken()
	claim := &ClaimsData{
		Data:           models.User{},
		StandardClaims: jwt.StandardClaims{},
	}
	successfulToken.Claims = claim
	scenarios := []TestParseJWTScenario{
		{
			description: "token successful",
			tokenString: "token",
			token:       successfulToken,
			tokenErr:    nil,
			expect:      claim,
			noErr:       true,
		},
		{
			description: "token unsuccessful",
			tokenString: "token",
			token:       constructToken(),
			tokenErr:    errors.New("token err"),
			expect:      nil,
			noErr:       false,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyFunc(os.Getenv, func(key string) string {
				return "env"
			})
			patches.ApplyFunc(jwt.ParseWithClaims, func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
				return scenario.token, scenario.tokenErr
			})
			res, err := ParseJWT(scenario.tokenString)
			assert.Equal(t, func() bool {
				if res == nil && scenario.expect == nil {
					return true
				} else if res == nil || scenario.expect == nil {
					return false
				}
				return reflect.DeepEqual(*res, *scenario.expect)
			}(), true)
			assert.Equal(t, err == nil, scenario.noErr)
		})
	}
}

type TestParseJWTWithClaimsScenario struct {
	description string
	tokenString string
	claims      *ClaimsData
	token       *jwt.Token
	tokenErr    error
	noErr       bool
}

func TestParseJWTWithClaims(t *testing.T) {
	invalidToken := constructToken()
	invalidToken.Valid = false
	scenarios := []TestParseJWTWithClaimsScenario{
		{
			description: "normal",
			tokenString: "token",
			claims:      nil,
			token:       constructToken(),
			tokenErr:    nil,
			noErr:       true,
		},
		{
			description: "token invalid",
			tokenString: "token",
			claims:      nil,
			token:       invalidToken,
			tokenErr:    errors.New("token err"),
			noErr:       false,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyFunc(os.Getenv, func(key string) string {
				return "env"
			})
			patches.ApplyFunc(jwt.ParseWithClaims, func(tokenString string, claims jwt.Claims, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
				return scenario.token, scenario.tokenErr
			})
			defer patches.Reset()
			err := ParseJWTWithClaims(scenario.tokenString, scenario.claims)
			assert.Equal(t, err == nil, scenario.noErr)
		})
	}
}

type TestHashStringScenario struct {
	description string
	tokenString string
	hash        string
	hashErr     error
	expect      string
}

func TestHashString(t *testing.T) {
	var logger *log.Logger
	scenarios := []TestHashStringScenario{
		{
			description: "create hash err",
			tokenString: "token",
			hash:        "",
			hashErr:     errors.New("create hash err"),
			expect:      "",
		},
		{
			description: "normal",
			tokenString: "token",
			hash:        "hash",
			hashErr:     nil,
			expect:      "hash",
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyFunc(argon2id.CreateHash, func(password string, params *argon2id.Params) (hash string, err error) {
				return scenario.hash, scenario.hashErr
			})
			patches.ApplyMethod(reflect.TypeOf(logger), "Println", func(logger *log.Logger, v ...interface{}) {})
			defer patches.Reset()
			res := HashString(scenario.tokenString)
			assert.Equal(t, res == scenario.expect, true)
		})
	}
}

type TestValidateJWTScenario struct {
	description string
	r           *http.Request
	authHeader  string
	claims      *ClaimsData
	claimsErr   error
	expect      *ClaimsData
	noErr       bool
}

func TestValidateJWT(t *testing.T) {

	r := constructHTTPRequest()
	var header http.Header
	claims := &ClaimsData{
		Data:           models.User{},
		StandardClaims: jwt.StandardClaims{},
	}
	scenarios := []TestValidateJWTScenario{
		{
			description: "unauthorized",
			r:           r,
			authHeader:  "header",
			expect:      nil,
			noErr:       false,
		},
		{
			description: "unauthenticated",
			r:           r,
			authHeader:  "Bearer header",
			claims:      nil,
			claimsErr:   errors.New("claims err"),
			expect:      nil,
			noErr:       false,
		},
		{
			description: "normal",
			r:           constructHTTPRequest(),
			authHeader:  "Bearer header",
			claims:      claims,
			claimsErr:   nil,
			expect:      claims,
			noErr:       true,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(header), "Get", func(header http.Header, key string) string {
				return scenario.authHeader
			})
			patches.ApplyFunc(ParseJWT, func(tokenString string) (*ClaimsData, error) {
				return scenario.claims, scenario.claimsErr
			})
			defer patches.Reset()
			res, err := ValidateJWT(scenario.r)
			assert.Equal(t, func() bool {
				if res == nil && scenario.expect == nil {
					return true
				} else if res == nil || scenario.expect == nil {
					return false
				}
				return reflect.DeepEqual(*res, *scenario.expect)
			}(), true)
			assert.Equal(t, err == nil, scenario.noErr)
		})
	}
}

type TestIsAdminScenario struct {
	description string
	user        models.User
	db          *gorm.DB
	expect      bool
}

func TestIsAdmin(t *testing.T) {
	dbErr := constructGormDB()
	dbErr.Error = gorm.ErrRecordNotFound
	db := constructGormDB()
	scenarios := []TestIsAdminScenario{
		{
			description: "record not found err",
			user:        constructUser(),
			db:          dbErr,
			expect:      false,
		},
		{
			description: "normal",
			user:        constructUser(),
			db:          db,
			expect:      true,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			_ = gomonkey.ApplyMethod(reflect.TypeOf(db), "Take", func(db *gorm.DB, dest interface{}, conds ...interface{}) *gorm.DB {
				return scenario.db
			})
			res := IsAdmin(scenario.user, scenario.db)
			assert.Equal(t, res == scenario.expect, true)
		})
	}
}

type TestIsSupervisorScenario struct {
	description string
	user        models.User
	moduleId    uint
	db          *gorm.DB
	expect      bool
}

func TestIsSupervisor(t *testing.T) {
	dbErr := constructGormDB()
	dbErr.Error = gorm.ErrRecordNotFound
	db := constructGormDB()
	scenarios := []TestIsSupervisorScenario{
		{
			description: "record not found err",
			user:        constructUser(),
			db:          dbErr,
			moduleId:    uint(1),
			expect:      false,
		},
		{
			description: "normal",
			user:        constructUser(),
			moduleId:    uint(1),
			db:          db,
			expect:      true,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			_ = gomonkey.ApplyMethod(reflect.TypeOf(db), "Take", func(db *gorm.DB, dest interface{}, conds ...interface{}) *gorm.DB {
				return scenario.db
			})
			res := IsSupervisor(scenario.user, scenario.moduleId, scenario.db)
			assert.Equal(t, res == scenario.expect, true)
		})
	}
}

type TestIsEnrolledScenario struct {
	description string
	user        models.User
	moduleId    uint
	db          *gorm.DB
	isAdmin     bool
	expect      bool
}

func TestIsEnrolled(t *testing.T) {
	dbErr := constructGormDB()
	dbErr.Error = gorm.ErrRecordNotFound
	db := constructGormDB()
	scenarios := []TestIsEnrolledScenario{
		{
			description: "bypass admin",
			user:        constructUser(),
			db:          dbErr,
			moduleId:    uint(1),
			isAdmin:     true,
			expect:      true,
		},
		{
			description: "record not found err",
			user:        constructUser(),
			db:          dbErr,
			moduleId:    uint(1),
			expect:      false,
		},
		{
			description: "normal",
			user:        constructUser(),
			moduleId:    uint(1),
			db:          db,
			expect:      true,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyFunc(IsAdmin, func(user models.User, db *gorm.DB) bool {
				return scenario.isAdmin
			})
			patches.ApplyMethod(reflect.TypeOf(db), "Take", func(db *gorm.DB, dest interface{}, conds ...interface{}) *gorm.DB {
				return scenario.db
			})
			res := IsEnrolled(scenario.user, scenario.moduleId, scenario.db)
			assert.Equal(t, res == scenario.expect, true)
		})
	}
}

type TestIsMarkerScenario struct {
	description  string
	user         models.User
	assignmentId uint
	studentId    uint
	db           *gorm.DB
	isAdmin      bool
	isSupervisor bool
	expect       bool
}

func TestIsMarker(t *testing.T) {
	dbErr := constructGormDB()
	dbErr.Error = gorm.ErrRecordNotFound
	db := constructGormDB()
	scenarios := []TestIsMarkerScenario{
		{
			description:  "bypass admin",
			user:         constructUser(),
			db:           dbErr,
			assignmentId: uint(1),
			studentId:    uint(1),
			isAdmin:      true,
			expect:       true,
		},
		{
			description:  "bypass supervisor",
			user:         constructUser(),
			db:           dbErr,
			assignmentId: uint(1),
			studentId:    uint(1),
			isSupervisor: true,
			expect:       true,
		},
		{
			description:  "record not found err",
			user:         constructUser(),
			db:           dbErr,
			assignmentId: uint(1),
			studentId:    uint(1),
			expect:       false,
		},
		{
			description:  "normal",
			user:         constructUser(),
			assignmentId: uint(1),
			studentId:    uint(1),
			db:           db,
			expect:       true,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyFunc(IsAdmin, func(user models.User, db *gorm.DB) bool {
				return scenario.isAdmin
			})
			patches.ApplyFunc(IsSupervisor, func(user models.User, moduleId uint, db *gorm.DB) bool {
				return scenario.isSupervisor
			})
			patches.ApplyMethod(reflect.TypeOf(db), "Model", func(db *gorm.DB, value interface{}) *gorm.DB {
				return scenario.db
			})
			patches.ApplyMethod(reflect.TypeOf(db), "Where", func(db *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
				return scenario.db
			})
			patches.ApplyMethod(reflect.TypeOf(db), "Find", func(db *gorm.DB, dest interface{}, conds ...interface{}) *gorm.DB {
				return scenario.db
			})
			patches.ApplyMethod(reflect.TypeOf(db), "Take", func(db *gorm.DB, dest interface{}, conds ...interface{}) *gorm.DB {
				return scenario.db
			})
			res := IsMarker(scenario.user, scenario.assignmentId, scenario.studentId, scenario.db)
			assert.Equal(t, res == scenario.expect, true)
		})
	}
}

type TestIsRevieweeScenario struct {
	description  string
	user         models.User
	assignmentId uint
	markerId     uint
	db           *gorm.DB
	isAdmin      bool
	expect       bool
}

func TestIsReviewee(t *testing.T) {
	dbErr := constructGormDB()
	dbErr.Error = gorm.ErrRecordNotFound
	db := constructGormDB()
	scenarios := []TestIsRevieweeScenario{
		{
			description:  "bypass admin",
			user:         constructUser(),
			db:           dbErr,
			assignmentId: uint(1),
			markerId:     uint(1),
			isAdmin:      true,
			expect:       true,
		},
		{
			description:  "record not found err",
			user:         constructUser(),
			db:           dbErr,
			assignmentId: uint(1),
			markerId:     uint(1),
			expect:       false,
		},
		{
			description:  "normal",
			user:         constructUser(),
			assignmentId: uint(1),
			markerId:     uint(1),
			db:           db,
			expect:       true,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyFunc(IsAdmin, func(user models.User, db *gorm.DB) bool {
				return scenario.isAdmin
			})
			patches.ApplyMethod(reflect.TypeOf(db), "Take", func(db *gorm.DB, dest interface{}, conds ...interface{}) *gorm.DB {
				return scenario.db
			})
			res := IsReviewee(scenario.user, scenario.assignmentId, scenario.markerId, scenario.db)
			assert.Equal(t, res == scenario.expect, true)
		})
	}
}

func TestIsPair(t *testing.T) {
	convey.Convey("Test IsPair", t, func() {
		patches := gomonkey.ApplyFunc(IsMarker, func(user models.User, assignmentId uint, studentId uint, db *gorm.DB) bool {
			return true
		})
		patches.ApplyFunc(IsReviewee, func(user models.User, assignmentId uint, markerId uint, db *gorm.DB) bool {
			return true
		})
		res := IsPair(constructUser(), 1, 1, constructGormDB())
		assert.Equal(t, res == true, true)
	})
}

func TestIsMemberOf(t *testing.T) {
	convey.Convey("Test IsMemberOf", t, func() {
		patches := gomonkey.ApplyFunc(IsEnrolled, func(user models.User, moduleId uint, db *gorm.DB) bool {
			return true
		})
		patches.ApplyFunc(IsSupervisor, func(user models.User, moduleId uint, db *gorm.DB) bool {
			return true
		})
		res := IsMemberOf(constructUser(), 1, constructGormDB())
		assert.Equal(t, res == true, true)
	})
}
