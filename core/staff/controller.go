package staff

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

type StaffController struct {
	DB *gorm.DB
}

func (controller StaffController) CreateRouters(route *mux.Router) {
	controller.CreateAuthRouter(route.PathPrefix("/auth").Subrouter())
	controller.CreatePrivilegedRouter(route.NewRoute().Subrouter())
}

func (controller StaffController) CreateAuthRouter(route *mux.Router) {
	loginRoute := route.NewRoute().Subrouter()
	loginRoute.Use(utils.DecodeParamsMiddleware(&models.User{}))
	loginRoute.Use(utils.AccountExistCheckMiddleware(controller.DB, &models.PendingStaff{}, utils.DecodeParamsContextKey, false, "Wait for approval, please send an email to fellowmarksystem@gmail.com"))
	loginRoute.HandleFunc("/login", utils.LoginHandleFunc(controller.DB, utils.ModelDBScope(&models.Staff{}))).Methods(http.MethodGet)

	signUpRoute := route.NewRoute().Subrouter()
	signUpRoute.Use(utils.DecodeBodyMiddleware(&models.User{}))
	signUpRoute.Use(utils.SanitizeDataMiddleware())
	signUpRoute.Use(utils.UserPasswordHashMiddleware)
	signUpRoute.Use(utils.AccountExistCheckMiddleware(controller.DB, &models.PendingStaff{}, utils.DecodeBodyContextKey, false, "Wait for approval, please send an email to fellowmarksystem@gmail.com"))
	signUpRoute.Use(utils.AccountExistCheckMiddleware(controller.DB, &models.Staff{}, utils.DecodeBodyContextKey, false, "Staff account already exists"))
	signUpRoute.HandleFunc("/signup", utils.UserCreateHandleFunc(controller.DB, &models.PendingStaff{})).Methods(http.MethodPost)
}

func (controller StaffController) CreatePrivilegedRouter(route *mux.Router) {
	route.Use(utils.AuthenticationMiddleware())
	controller.CreateStaffApproveRoute(route.PathPrefix("/approve").Subrouter())
	controller.GetPairingsRoute(route.PathPrefix("/module/{moduleId}/pairing").Subrouter())
	controller.GetStaffsRoute(route.NewRoute().Subrouter())
}

func (controller StaffController) CreateStaffApproveRoute(route *mux.Router) {
	route.Use(utils.IsAdminMiddleware(controller.DB))
	route.Use(utils.DecodeBodyMiddleware(&models.User{}))
	route.Use(utils.SanitizeDataMiddleware())
	route.Use(utils.AccountExistCheckMiddleware(controller.DB, &models.PendingStaff{}, utils.DecodeBodyContextKey, true, "No such signup request"))
	route.Use(utils.AccountExistCheckMiddleware(controller.DB, &models.Staff{}, utils.DecodeBodyContextKey, false,  "Staff account already exists"))
	route.HandleFunc("", controller.StaffApproveHandleFunc()).Methods(http.MethodPost)
}

func (controller StaffController) GetPairingsRoute(route *mux.Router) {
	route.Use(utils.DecodeParamsMiddleware(&models.Pairing{}))
	route.HandleFunc("", utils.DBGetFromDataParams(controller.DB, &models.Pairing{}, &[]models.Pairing{})).Methods(http.MethodGet)
}

func (controller StaffController) GetStaffsRoute(route *mux.Router) {
	route.Use(utils.IsAdminMiddleware(controller.DB))
	route.Use(utils.DecodeParamsMiddleware(&models.User{}))
	route.Use(utils.SanitizeDataMiddleware())
	route.HandleFunc("/pending", utils.DBGetFromDataParams(controller.DB, &models.PendingStaff{}, &[]models.PendingStaff{})).Methods(http.MethodGet)
	route.HandleFunc("/approve", utils.DBGetFromDataParams(controller.DB, &models.Staff{}, &[]models.Staff{})).Methods(http.MethodGet)
}
