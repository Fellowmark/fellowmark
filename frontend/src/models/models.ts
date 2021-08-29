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
