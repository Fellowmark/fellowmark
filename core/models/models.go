package models

import (
	"time"
)

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	Model
	Email    string
	Name     string
	Password string
}

type Admin struct {
	Model
	Email    string `gorm:"type:varchar(100);column:email;unique" validate:"nonzero"`
	Name     string `gorm:"type:varchar(255);column:name;not null"`
	Password string `gorm:"column:password;not null" validate:"min=8"`
}

type Student struct {
	Model
	Email    string `gorm:"type:varchar(100);column:email;unique" validate:"nonzero"`
	Name     string `gorm:"type:varchar(255);column:name;not null"`
	Password string `gorm:"column:password;not null" validate:"min=8"`
}

type Staff struct {
	Model
	Email    string `gorm:"type:varchar(100);column:email;unique" validate:"nonzero"`
	Name     string `gorm:"type:varchar(255);column:name;not null"`
	Password string `gorm:"column:password;not null" validate:"min=8"`
}

type PendingStaff struct {
	Model
	Email    string `gorm:"type:varchar(100);column:email;unique" validate:"nonzero"`
	Name     string `gorm:"type:varchar(255);column:name;not null"`
	Password string `gorm:"column:password;not null" validate:"min=8"`
}

type Module struct {
	Model
	Code     string `gorm:"index:moduleIdx,unique;type:varchar(8);column:code;not null"`
	Semester string `gorm:"index:moduleIdx,unique;type:varchar(40);column:semester;not null"`
	Name     string `gorm:"column:name;not null"`
}

type Enrollment struct {
	Model
	Module    Module  `gorm:"foreignKey:ModuleID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ModuleID  uint    `gorm:"uniqueIndex:enrollmentIdx;column:module_id;not null"`
	Student   Student `gorm:"foreignKey:StudentID"`
	StudentID uint    `gorm:"uniqueIndex:enrollmentIdx;column:student_id;not null"`
}

type Supervision struct {
	Model
	ModuleID uint  `gorm:"uniqueIndex:supervisionIdx;column:module_id;not null"`
	Staff    Staff `gorm:"foreignKey:StaffID"`
	StaffID  uint  `gorm:"uniqueIndex:supervisionIdx;column:staff_id;not null"`
}

type Assistance struct {
	Model
	Module    Module  `gorm:"foreignKey:ModuleID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	ModuleID  uint    `gorm:"uniqueIndex:assistanceIdx;column:module_id;not null"`
	Student   Student `gorm:"foreignKey:StudentID"`
	StudentID uint    `gorm:"uniqueIndex:assistanceIdx;column:student_id;not null"`
}

type Assignment struct {
	Model
	Name      string `gorm:"uniqueIndex:assignmentIdx;column:name;not null"`
	Module    Module `gorm:"foreignKey:ModuleID;constraint:OnDelete:CASCADE" json:"-"`
	ModuleID  uint   `gorm:"column:module_id;not null"`
	GroupSize int    `gorm:"uniqueIndex:assignmentIdx;column:group_size;not null;check:group_size > 0"`
	Deadline  int64  `gorm:"not null"`
}

type Question struct {
	Model
	QuestionNumber uint       `gorm:"uniqueIndex:questionIdx;column:question_number;not null"`
	QuestionText   string     `gorm:"column:question_text;not null"`
	Assignment     Assignment `gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE" json:"-"`
	AssignmentID   uint       `gorm:"uniqueIndex:questionIdx;column:assignment_id;not null"`
	StartDate      time.Time  `gorm:"column:start_date"`
	EndDate        time.Time  `gorm:"column:end_date"`
}

type Rubric struct {
	Model
	Question    Question `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"-"`
	QuestionID  uint     `gorm:"uniqueIndex:rubricIdx;column:question_id;not null"`
	Criteria    string   `gorm:"uniqueIndex:rubricIdx;not null"`
	Description string   `gorm:"uniqueIndex:rubricIdx;not null"`
	MinMark     int      `gorm:"column:min_mark;default:0"`
	MaxMark     int      `gorm:"column:max_mark;default:10"`
}

type Pairing struct {
	Model
	Assignment   Assignment `gorm:"foreignKey:AssignmentID;constraint:OnDelete:CASCADE" json:"-"`
	AssignmentID uint       `gorm:"column:assignment_id;not null"`
	Student      Student    `gorm:"foreignKey:StudentID"`
	StudentID    uint       `gorm:"column:student_id;not null"`
	Marker       Student    `gorm:"foreignKey:MarkerID"`
	MarkerID     uint       `gorm:"column:marker_id;not null"`
	Active       bool       `gorm:"not null"`
}

type Submission struct {
	Model
	SubmittedBy Student  `gorm:"foreignKey:StudentID"`
	StudentID   uint     `gorm:"uniqueIndex:submissionIdx;column:submitted_by;not null"`
	Question    Question `gorm:"foreignKey:QuestionID" json:"-"`
	QuestionID  uint     `gorm:"uniqueIndex:submissionIdx;column:question_id;not null"`
	ContentFile string   `gorm:"column:content_file_location" json:"-"`
	Content     string   `gorm:"column:content"`
}

type Grade struct {
	Model
	Pairing   Pairing `gorm:"foreignKey:PairingID"`
	PairingID uint    `gorm:"column:pairing_id;not null"`
	Rubric    Rubric  `gorm:"foreignKey:RubricID"`
	RubricID  uint    `gorm:"column:rubric_id;not null"`
	Grade     int     `gorm:"column:grade;not null"`
	Comment   string  `gorm:"column:comment"`
}
