package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Email    string `gorm:"type:varchar(20);column:email;unique"`
	Name     string `gorm:"type:varchar(50);column:name;not null"`
	Password string `gorm:"column:password;not null"`
}

type Staff struct {
	gorm.Model
	Email    string `gorm:"type:varchar(20);column:email;unique"`
	Name     string `gorm:"type:varchar(50);column:name;not null"`
	Password string `gorm:"column:password;not null"`
}

type Module struct {
	gorm.Model
	Code     string `gorm:"type:varchar(8);column:code;not null"`
	Semester string `gorm:"type:varchar(6);column:semester;not null"`
	Name     string `gorm:"column:name;not null"`
}

type Enrollment struct {
	gorm.Model
	Module    Module  `gorm:"foreignKey:ModuleID;references:ID"`
	ModuleID  uint    `gorm:"column:module_id;not null"`
	Student   Student `gorm:"foreignKey:StudentID"`
	StudentID uint    `gorm:"column:student_id;not null"`
}

type Supervision struct {
	gorm.Model
	Module   Module `gorm:"foreignKey:ModuleID"`
	ModuleID uint   `gorm:"column:module_id;not null"`
	Staff    Staff  `gorm:"foreignKey:StaffID"`
	StaffID  uint   `gorm:"column:staff_id;not null"`
}

type Assignment struct {
	gorm.Model
	Name      string `gorm:"column:name;not null"`
	Module    Module `gorm:"foreignKey:ModuleID"`
	ModuleID  uint   `gorm:"column:module_id;not null"`
	GroupSize int    `gorm:"column:group_size;not null;check:group_size > 0"`
}

type Question struct {
	gorm.Model
	QuestionNumber uint       `gorm:"column:question_number;not null"`
	QuestionText   string     `gorm:"column:question_text;not null"`
	Assignment     Assignment `gorm:"foreignKey:AssignmentID"`
	AssignmentID   uint       `gorm:"column:assignment_id;not null"`
}

type Rubric struct {
	gorm.Model
	Question    Question `gorm:"foreignKey:QuestionID"`
	QuestionID  uint     `gorm:"column:question_id;not null"`
	Criteria    string   `gorm:"not null"`
	Description string   `gorm:"not null"`
	MinMark     int      `gorm:"column:min_mark;default:0"`
	MaxMark     int      `gorm:"column:max_mark;default:10"`
}

type Pairing struct {
	gorm.Model
	Assignment   Assignment `gorm:"foreignKey:AssignmentID"`
	AssignmentID uint       `gorm:"column:assignment_id;not null"`
	Student      Student    `gorm:"foreignKey:StudentID"`
	StudentID    uint       `gorm:"column:student_id;not null"`
	Marker       Student    `gorm:"foreignKey:MarkerID"`
	MarkerID     uint       `gorm:"column:marker_id;not null"`
	Active       bool       `gorm:"not null"`
}

type Submission struct {
	gorm.Model
	SubmittedBy Student  `gorm:"foreignKey:StudentID"`
	StudentID   uint     `gorm:"column:submitted_by;not null"`
	Question    Question `gorm:"foreignKey:QuestionID"`
	QuestionID  uint     `gorm:"column:question_id;not null"`
	ContentFile string   `gorm:"column:content_file_location"`
	Content     string   `gorm:"column:content"`
}

type Grade struct {
	gorm.Model
	Pairing   Pairing `gorm:"foreignKey:PairingID"`
	PairingID uint    `gorm:"column:pairing_id;not null"`
	Rubric    Rubric  `gorm:"foreignKey:RubricID"`
	RubricID  uint    `gorm:"column:rubric_id;not null"`
	Grade     int     `gorm:"column:grade;not null"`
}
