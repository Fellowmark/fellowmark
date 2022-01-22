package utils

import (
	"math/rand"
	"time"

	"github.com/nus-utils/nus-peer-review/models"

	"gorm.io/gorm"
)

func InitializePairings(db *gorm.DB, assignment models.Assignment) *gorm.DB {
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
				SELECT student_id FROM enrollments
				WHERE module_id = ?
			) AND B.id in (
				SELECT student_id FROM enrollments
				WHERE  module_id = ?
			)
		) RETURNING *`,
		assignment.ID,
		false,
		assignment.ModuleID,
		assignment.ModuleID,
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
		assignment.ModuleID,
	).Scan(&enrolledStudents)

	var newPairings []models.Pairing
	if result.Error != nil {
		return newPairings, result
	}

	noOfSmallerGroups := 0
	noOfLargerGroups := len(enrolledStudents) / assignment.GroupSize
	if len(enrolledStudents)%assignment.GroupSize != 0 {
		noOfSmallerGroups = assignment.GroupSize - (len(enrolledStudents) % assignment.GroupSize) // groups of size GroupSize - 1
		noOfLargerGroups = (len(enrolledStudents) / assignment.GroupSize) + 1 - noOfSmallerGroups // groups of size GroupSize
	}

	shuffleStudents(enrolledStudents)
	index := 0

	newPairings = append(newPairings, generateGroupPairings(db, noOfLargerGroups, assignment.GroupSize, assignment, &index, enrolledStudents)...)
	newPairings = append(newPairings, generateGroupPairings(db, noOfSmallerGroups, assignment.GroupSize-1, assignment, &index, enrolledStudents)...)

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

func generateGroupPairings(db *gorm.DB, noOfGroups, groupSize int, assignment models.Assignment, index *int, students []models.Student) []models.Pairing {
	var pairings []models.Pairing
	for i := 0; i < noOfGroups; i++ {
		var group []models.Student
		for j := 0; j < groupSize; j++ {
			group = append(group, students[*index])
			*index++
		}
		for _, student := range group {
			for _, marker := range group {
				if student.ID != marker.ID {
					var newPairingForStudent models.Pairing
					db.Where("student_id = ? AND marker_id = ? AND assignment_id = ?", student.ID, marker.ID, assignment.ID).Find(&newPairingForStudent)
					pairings = append(pairings, newPairingForStudent)
				}
			}
		}
	}
	return pairings
}

func shuffleStudents(students []models.Student) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := len(students); i > 0; i-- {
		randomIndex := r.Intn(i)
		students[i-1], students[randomIndex] = students[randomIndex], students[i-1]
	}
}
