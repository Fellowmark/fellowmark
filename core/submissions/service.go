package submissions

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/nus-utils/nus-peer-review/models"
	"github.com/nus-utils/nus-peer-review/utils"
	"gorm.io/gorm"
)

const FilePathContextKey = "file"

func (controller FileserverController) UploadPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			submission := r.Context().Value(utils.DecodeParamsContextKey).(*models.Submission)

			var question models.Question
			controller.DB.Model(&models.Assignment{}).Where("id = ?", submission.QuestionID).Find(&question)

			var assignment models.Assignment
			controller.DB.Model(&models.Assignment{}).Where("id = ?", question.AssignmentID).Find(&assignment)

			if pass := utils.IsEnrolled(*claims, assignment.ModuleID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}

func (controller FileserverController) DownloadPermissionCheck() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := r.Context().Value(utils.JWTClaimContextKey).(*models.User)
			submission := r.Context().Value(utils.DecodeParamsContextKey).(*models.Submission)

			var question models.Question
			controller.DB.Model(&models.Assignment{}).Where("id = ?", submission.QuestionID).Find(&question)

			if pass := utils.IsPair(*claims, question.AssignmentID, submission.StudentID, controller.DB); pass {
				next.ServeHTTP(w, r)
			} else {
				utils.HandleResponse(w, "Unauthorized", http.StatusUnauthorized)
			}
		})
	}
}

func (fr FileserverController) UpdateFilePathMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(utils.DecodeBodyContextKey).(*models.Submission)
			fr.DB.Where(data).First(&data)
			ctxWithPath := context.WithValue(r.Context(), FilePathContextKey, data.ContentFile)
			next.ServeHTTP(w, r.WithContext(ctxWithPath))
		})
	}
}

func (fr FileserverController) UpdateSubmissionContentFile() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(utils.DecodeBodyContextKey).(*models.Submission)
			filePath := r.Context().Value(FilePathContextKey).(string)
			data.ContentFile = filePath
			ctxWithPath := context.WithValue(r.Context(), utils.DecodeBodyContextKey, data)
			next.ServeHTTP(w, r.WithContext(ctxWithPath))
		})
	}
}

func (fr FileserverController) StoreUploadLocationInDB(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := r.Context().Value(utils.DecodeBodyContextKey).(*models.Submission)
		result := db.Model(&models.Submission{}).Omit("ID").Create(data)
		if result.Error != nil {
			result = db.Model(&models.Submission{}).Omit("ID").Where("submitted_by = ? AND question_id = ?", data.StudentID, data.QuestionID).Updates(data)
			if result.Error != nil {
				utils.HandleResponse(w, "Failed", http.StatusInternalServerError)
			}
		} else {
			utils.HandleResponseWithObject(w, data, http.StatusOK)
		}
	}
}

func (fr FileserverController) GetSubmissionMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data := r.Context().Value(utils.DecodeParamsContextKey).(*models.Submission)
			fr.DB.Model(&models.Submission{}).Where("submitted_by = ? AND question_id = ?", data.StudentID, data.QuestionID).FirstOrCreate(data)
			ctxWithPath := context.WithValue(r.Context(), utils.DecodeBodyContextKey, data)
			next.ServeHTTP(w, r.WithContext(ctxWithPath))
		})
	}
}

func (fr FileserverController) UploadMiddleware(uploadPath string, maxUploadSize int64) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseMultipartForm(maxUploadSize); err != nil {
				utils.HandleResponse(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
				return
			}

			// parse and validate file and post parameters
			file, _, err := r.FormFile("uploadFile")
			if err != nil {
				utils.HandleResponse(w, "INVALID_FILE", http.StatusBadRequest)
				return
			}

			defer file.Close()
			// Get and print out file size
			// fileSize := fileHeader.Size
			// validate file size
			// if fileSize > maxUploadSize {
			// 	utils.HandleResponse(w, "FILE_TOO_BIG", http.StatusBadRequest)
			// 	return
			// }
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				utils.HandleResponse(w, "INVALID_FILE", http.StatusBadRequest)
				return
			}

			// check file type, detectcontenttype only needs the first 512 bytes
			detectedFileType := http.DetectContentType(fileBytes)
			switch detectedFileType {
			case "image/jpeg", "image/jpg":
			case "image/gif", "image/png":
			case "application/pdf":
				break
			default:
				utils.HandleResponse(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
				return
			}
			fileName := utils.RandToken(20)
			fileEndings, err := mime.ExtensionsByType(detectedFileType)
			if err != nil {
				utils.HandleResponse(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
				return
			}
			newPath := filepath.Join(uploadPath, fileName+fileEndings[0])

			// write file
			newFile, err := os.Create(newPath)
			if err != nil {
				utils.HandleResponse(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}
			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				utils.HandleResponse(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}
			ctxWithPath := context.WithValue(r.Context(), FilePathContextKey, newPath)
			next.ServeHTTP(w, r.WithContext(ctxWithPath))
		})
	}
}

func (fr FileserverController) DownloadHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := r.Context().Value(FilePathContextKey).(string)
		data := r.Context().Value(utils.DecodeBodyContextKey)
		fn := filepath.Base(filePath)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fn))
		file, err := os.Open(filePath)
		if err != nil {
			utils.HandleResponse(w, "FILE_NOT_FOUND", http.StatusNotFound)
			return
		}
		defer file.Close()
		io.Copy(w, file)
		utils.HandleResponseWithObject(w, data, http.StatusOK)
	})
}
