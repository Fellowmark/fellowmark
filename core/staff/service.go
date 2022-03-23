package staff

import (
	"net/http"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

func (controller StaffController) StaffApproveHandleFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(utils.DecodeBodyContextKey).(*models.User)
		txError := controller.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(&models.PendingStaff{}, data.ID).Error; err != nil {
				return err
			}
			data.ID = 0
			if err := tx.Model(&models.Staff{}).Create(data).Error; err != nil {
				return err
			}
			return nil
		})
		if txError != nil {
			utils.HandleResponse(w, txError.Error(), http.StatusBadRequest)
			return
		}
		utils.HandleResponseWithObject(w, data, http.StatusOK)
	}
}