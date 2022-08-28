package utils

import (
	"github.com/agiledragon/gomonkey"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"math/rand"
	"reflect"
	"testing"
)

type TestInitializeParingsScenario struct {
	description string
	db          *gorm.DB
	assignment  models.Assignment
	expect      *gorm.DB
}

func TestInitializePairings(t *testing.T) {
	dbErr := constructGormDB()
	dbErr.Error = gorm.ErrRecordNotFound
	db := constructGormDB()
	scenarios := []TestInitializeParingsScenario{
		{
			description: "result err",
			db:          dbErr,
			assignment:  constructAssignment(),
			expect:      dbErr,
		},
		{
			description: "normal",
			db:          db,
			assignment:  constructAssignment(),
			expect:      db,
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(db), "Delete", func(db *gorm.DB, value interface{}, conds ...interface{}) *gorm.DB {
				return scenario.db
			})
			patches.ApplyMethod(reflect.TypeOf(db), "Exec", func(db *gorm.DB, sql string, values ...interface{}) *gorm.DB {
				return scenario.db
			})
			res := InitializePairings(scenario.db, scenario.assignment)
			assert.Equal(t, reflect.DeepEqual(*res, *scenario.expect), true)
		})
	}
}

type TestSetNewPairingsScenario struct {
	description                 string
	db                          *gorm.DB
	assignment                  models.Assignment
	newPairings                 []models.Pairing
	getNewPairingsResult        *gorm.DB
	deactivateOldPairingsResult *gorm.DB
	activateNewPairingsResult   *gorm.DB
	expect                      *gorm.DB
}

func TestSetNewPairings(t *testing.T) {
	dbErr := constructGormDB()
	dbErr.Error = gorm.ErrRecordNotFound
	scenarios := []TestSetNewPairingsScenario{
		{
			description:          "get new pairing err",
			db:                   constructGormDB(),
			assignment:           constructAssignment(),
			newPairings:          []models.Pairing{},
			getNewPairingsResult: dbErr,
			expect:               dbErr,
		},
		{
			description:                 "deactivate old parings err",
			db:                          constructGormDB(),
			assignment:                  constructAssignment(),
			newPairings:                 []models.Pairing{},
			getNewPairingsResult:        constructGormDB(),
			deactivateOldPairingsResult: dbErr,
			expect:                      dbErr,
		},
		{
			description:                 "normal",
			db:                          constructGormDB(),
			assignment:                  constructAssignment(),
			newPairings:                 []models.Pairing{},
			getNewPairingsResult:        constructGormDB(),
			deactivateOldPairingsResult: constructGormDB(),
			activateNewPairingsResult:   constructGormDB(),
			expect:                      constructGormDB(),
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyFunc(getNewPairings, func(db *gorm.DB, assignment models.Assignment) ([]models.Pairing, *gorm.DB) {
				return scenario.newPairings, scenario.getNewPairingsResult
			})
			patches.ApplyFunc(deactivateOldPairings, func(db *gorm.DB) *gorm.DB {
				return scenario.deactivateOldPairingsResult
			})
			patches.ApplyFunc(activateNewPairings, func(db *gorm.DB, pairings []models.Pairing) *gorm.DB {
				return scenario.activateNewPairingsResult
			})
			res := SetNewPairings(scenario.db, scenario.assignment)
			assert.Equal(t, reflect.DeepEqual(*res, *scenario.expect), true)
		})
	}
}

type TestGetNewPairingsScenario struct {
	description   string
	db            *gorm.DB
	assignment    models.Assignment
	dbRawResult   *gorm.DB
	expectParings []models.Pairing
	expect        *gorm.DB
}

//todo: fix problems

func TestGetNewPairings(t *testing.T) {
	dbErr := constructGormDB()
	dbErr.Error = gorm.ErrRecordNotFound
	scenarios := []TestGetNewPairingsScenario{
		{
			description:   "db raw err",
			db:            dbErr,
			assignment:    constructAssignment(),
			dbRawResult:   dbErr,
			expectParings: nil,
			expect:        dbErr,
		},
		{
			description:   "normal",
			db:            constructGormDB(),
			assignment:    constructAssignment(),
			dbRawResult:   constructGormDB(),
			expectParings: []models.Pairing{constructParing(), constructParing()},
			expect:        constructGormDB(),
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(dbErr), "Raw", func(db *gorm.DB, sql string, values ...interface{}) *gorm.DB {
				return scenario.dbRawResult
			})
			patches.ApplyMethod(reflect.TypeOf(dbErr), "Scan", func(db *gorm.DB, dest interface{}) *gorm.DB {
				return scenario.dbRawResult
			})
			patches.ApplyFunc(shuffleStudents, func(students []models.Student) {})
			patches.ApplyFunc(generateGroupPairings, func(db *gorm.DB, noOfGroups, groupSize int, assignment models.Assignment, index *int, students []models.Student) []models.Pairing {
				return []models.Pairing{constructParing()}
			})
			pairings, res := getNewPairings(scenario.db, scenario.assignment)
			assert.Equal(t, reflect.DeepEqual(*res, *scenario.expect), true)
			assert.Equal(t, reflect.DeepEqual(pairings, scenario.expectParings), true)
		})
	}
}

func TestDeactivateOldPairings(t *testing.T) {
	db := constructGormDB()
	convey.Convey("test DeactivateOldPairings", t, func() {
		patches := gomonkey.ApplyMethod(reflect.TypeOf(db), "Model", func(db *gorm.DB, value interface{}) *gorm.DB {
			return db
		})
		patches.ApplyMethod(reflect.TypeOf(db), "Where", func(db *gorm.DB, query interface{}, args ...interface{}) *gorm.DB {
			return db
		})
		patches.ApplyMethod(reflect.TypeOf(db), "Update", func(db *gorm.DB, column string, value interface{}) *gorm.DB {
			return db
		})
		res := deactivateOldPairings(db)
		assert.Equal(t, reflect.DeepEqual(*res, *db), true)
	})
}

type TestActivateNewPairingsScenarios struct {
	description string
	db          *gorm.DB
	pairings    []models.Pairing
	modelResult *gorm.DB
	expect      *gorm.DB
}

func TestActivateNewPairings(t *testing.T) {
	dbErr := constructGormDB()
	dbErr.Error = gorm.ErrRecordNotFound
	scenarios := []TestActivateNewPairingsScenarios{
		{
			description: "db model err",
			db:          constructGormDB(),
			pairings:    []models.Pairing{constructParing()},
			modelResult: dbErr,
			expect:      dbErr,
		},
		{
			description: "normal",
			db:          constructGormDB(),
			pairings:    []models.Pairing{constructParing()},
			modelResult: constructGormDB(),
			expect:      constructGormDB(),
		},
	}
	for _, scenario := range scenarios {
		convey.Convey(scenario.description, t, func() {
			patches := gomonkey.ApplyMethod(reflect.TypeOf(dbErr), "Model", func(db *gorm.DB, value interface{}) *gorm.DB {
				return scenario.modelResult
			})
			patches.ApplyMethod(reflect.TypeOf(dbErr), "Update", func(db *gorm.DB, column string, value interface{}) *gorm.DB {
				return scenario.modelResult
			})
			res := activateNewPairings(scenario.db, scenario.pairings)
			assert.Equal(t, reflect.DeepEqual(*res, *scenario.expect), true)
		})
	}
}

func TestGenerateGroupPairings(t *testing.T) {

}

func TestShuffleStudents(t *testing.T) {
	r := &rand.Rand{}
	convey.Convey("test ShuffleStudents", t, func() {
		patches := gomonkey.ApplyFunc(rand.New, func(src rand.Source) *rand.Rand {
			return r
		})
		patches.ApplyMethod(reflect.TypeOf(r), "Intn", func(r *rand.Rand, n int) int {
			return 0
		})
		input := []models.Student{constructStudent("a"), constructStudent("b"), constructStudent("b")}
		expect := []models.Student{constructStudent("b"), constructStudent("b"), constructStudent("a")}
		shuffleStudents(input)
		assert.Equal(t, reflect.DeepEqual(input, expect), true)
	})
}
