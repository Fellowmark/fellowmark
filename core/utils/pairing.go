package utils

import (
	"github.com/nus-utils/nus-peer-review/models"

	"gorm.io/gorm"
)

func ResetPairings(db *gorm.DB, assignment models.Assignment) {
	db.Exec("DELETE FROM pairings WHERE assignment_id = ?", assignment.ID)
	db.Exec(
		`INSERT INTO pairings(assignment_id, student_id, marker_id, active) (
			SELECT ?, A.id, B.id, ?
			FROM students AS A
			CROSS JOIN students AS B
			WHERE A.id != B.id
			AND (
				(SELECT COUNT(id) FROM enrollments
				WHERE student_id = A.id AND module_id = ?) = 1
			) AND (
				(SELECT COUNT(id) FROM enrollments
				WHERE student_id = B.id AND module_id = ?) = 1
			)
		) RETURNING *`,
		assignment.ID,
		false,
		assignment.Module.ID,
		assignment.Module.ID,
	)
}

func SetNewPairings(db *gorm.DB, assignment models.Assignment) {
	newPairings := getNewPairings(db, assignment)
	deactivateOldPairings(db)
	activateNewPairings(db, newPairings)
}

func getNewPairings(db *gorm.DB, assignment models.Assignment) []models.Pairing {
	var enrolledStudents []models.Student
	db.Raw(
		`SELECT * FROM students 
		WHERE (
			(SELECT COUNT(student_id) FROM enrollments
			WHERE student_id = students.id AND module_id = ?)
 			= 1
		)`,
		assignment.Module.ID,
	).Scan(&enrolledStudents)

	var newPairings []models.Pairing
	for _, student := range enrolledStudents {
		var newPairingForStudent []models.Pairing
		db.Raw(
			`SELECT * FROM pairings
			WHERE active = ? AND student_id = ?
			ORDER BY random()
			LIMIT ?`,
			false,
			student.ID,
			assignment.GroupSize,
		).Scan(&newPairingForStudent)
		newPairings = append(newPairings, newPairingForStudent...)
	}

	return newPairings
}

func deactivateOldPairings(db *gorm.DB) {
	db.Exec(
		"UPDATE pairings SET active = ? WHERE active = ?",
		false,
		true,
	)
}

func activateNewPairings(db *gorm.DB, pairings []models.Pairing) {
	for _, pairing := range pairings {
		db.Exec(
			"UPDATE pairings SET active = ? WHERE id = ?",
			true,
			pairing.ID,
		)
	}
}
