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
