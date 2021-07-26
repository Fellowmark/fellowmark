package admin

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type AdminRoute struct {
	DB *gorm.DB
}

func (ur AdminRoute) CreateRouters(route *mux.Router) {
	ur.CreatePrivilegedRouter(route.NewRoute().Subrouter())
	ur.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
}

func (ur AdminRoute) CreateAuthRouter(route *mux.Router) {
	route.Use(utils.DecodeBodyMiddleware(&models.Admin{}, "user"))

	loginRoute := route.NewRoute().Subrouter()
	loginRoute.Use(ur.EmailCheck)
	loginRoute.Use(ur.PasswordCheck)
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(ur.DB, "Admin", "user")).Methods(http.MethodGet)
}

func (ur AdminRoute) CreatePrivilegedRouter(route *mux.Router) {
	if os.Getenv("RUN_ENV") != "production" {
		route.Use(utils.ValidateJWTMiddleware("Admin", "claims"))
	}

	ur.CreateStaffOptsRouter(route.PathPrefix("/staff").Subrouter())
	ur.CreateModuleOptsRouter(route.PathPrefix("/module").Subrouter())
}

func (ur AdminRoute) CreateStaffOptsRouter(route *mux.Router) {
	createStaffRoute := route.NewRoute().Subrouter()
	createStaffRoute.Use(utils.DecodeBodyMiddleware(&models.Staff{}, "user"))
	createStaffRoute.Use(utils.SanitizeDataMiddleware("user"))
	createStaffRoute.Use(ur.PasswordHash)
	createStaffRoute.HandleFunc("/", utils.DBCreateHandleFunc(ur.DB, "staffs", "user")).Methods(http.MethodPost)
}

func (ur AdminRoute) CreateModuleOptsRouter(route *mux.Router) {
	createModuleRoute := route.NewRoute().Subrouter()
	createModuleRoute.Use(ur.DecodeModuleJson)
	createModuleRoute.Use(ur.SanitizeModuleData)
	createModuleRoute.HandleFunc("/", ur.CreateModule).Methods(http.MethodPost)

	enrollModuleStaffRoute := route.NewRoute().Subrouter()
	enrollModuleStaffRoute.Use(ur.DecodeEnrollmentJson)
	enrollModuleStaffRoute.HandleFunc("/enroll", ur.EnrollModuleForStudent).Methods(http.MethodPost)

	superviseModuleRoute := route.NewRoute().Subrouter()
	superviseModuleRoute.Use(ur.DecodeSupervisesJson)
	superviseModuleRoute.HandleFunc("/supervision", ur.SuperviseModuleWithStaff).Methods(http.MethodPost)
}
