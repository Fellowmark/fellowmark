export interface User {
  ID?: number,
  Name?: string,
  Email?: string,
  Password?: string
}

export interface Enrollment {
  ModuleId?: number,
  Student?: User
}

export interface Assignment {
  ID?: number,
  Name?: string,
  ModuleID?: number,
  GroupSize?: number,
  Deadline?: number,
}

export interface Question {
  ID?: number,
  QuestionNumber?: number,
  QuestionText?: string,
  AssignmentID?: number
}

export interface Pairing {
  ID?: number,
  AssignmentID?: number,
  Student?: User,
  Marker?: User,
  Active?: Boolean
}

export interface Rubric {
  ID?: number,
  QuestionID?: number,
  Criteria?: string
  Description?: string
  MinMark?: number
  MaxMark?: number
}

export interface Grade {
  ID?: number,
  RubricID?: number,
  PairingID?: number,
  Grade?: number
  Comment?: string
}

export interface Submission {
  ID?: number
  StudentID?: number
  QuestionID?: number
  ContentFile?: string
  Content?: string
}
