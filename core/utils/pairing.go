package utils

import (
	"github.com/nus-utils/nus-peer-review/models"

	"gorm.io/gorm"
)

func ResetPairings(db *gorm.DB, assignment models.Assignment) *gorm.DB {
	result := db.Delete(&models.Pairing{}, assignment.ID)
	if result.Error != nil {
		return result
	}

	result = db.Exec(
		`INSERT INTO pairings(assignment_id, student_id, marker_id, active) (
			SELECT ?, A.id, B.id, ?
			FROM students AS A
			CROSS JOIN students AS B
			WHERE A.id != B.id
			AND A.id in (
				SELECT id FROM enrollments
				WHERE module_id = ? AND deleted_at IS NULL
			) AND B.id in (
				SELECT id FROM enrollments
				WHERE  module_id = ? AND deleted_at IS NULL
			)
		) RETURNING *`,
		assignment.ID,
		false,
		assignment.Module.ID,
		assignment.Module.ID,
	)
	return result
}

func SetNewPairings(db *gorm.DB, assignment models.Assignment) *gorm.DB {
	newPairings, result := getNewPairings(db, assignment)
	if result.Error != nil {
		return result
	}

	result = deactivateOldPairings(db)
	if result.Error != nil {
		return result
	}

	result = activateNewPairings(db, newPairings)
	return result
}

func getNewPairings(db *gorm.DB, assignment models.Assignment) ([]models.Pairing, *gorm.DB) {
	var enrolledStudents []models.Student
	var result *gorm.DB
	result = db.Raw(
		`SELECT students.id, email, name, password FROM students
		INNER JOIN enrollments
		ON enrollments.student_id = students.id
		WHERE module_id = ?`,
		assignment.Module.ID,
	).Scan(&enrolledStudents)

	var newPairings []models.Pairing
	if result.Error != nil {
		return newPairings, result
	}

	for _, student := range enrolledStudents {
		var newPairingForStudent []models.Pairing
		db.Where("active = ? AND student_id = ?", false, student.ID).Order("random()").Limit(assignment.GroupSize).Find(&newPairingForStudent)
		newPairings = append(newPairings, newPairingForStudent...)
	}

	return newPairings, result
}

func deactivateOldPairings(db *gorm.DB) *gorm.DB {
	return db.Model(&models.Pairing{}).Where("active = ?", true).Update("active", false)
}

func activateNewPairings(db *gorm.DB, pairings []models.Pairing) *gorm.DB {
	var result *gorm.DB
	for _, pairing := range pairings {
		result = db.Model(&pairing).Update("active", true)
		if result.Error != nil {
			return result
		}
	}

	return result
}
