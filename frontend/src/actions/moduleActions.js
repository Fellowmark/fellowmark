import axios from "axios";

/**
 * Creates module using given data
 * 
 * @param {string} moduleCode e.g. "CS2113T"
 * @param {string} semester e.g. "2122-1" for AY2021/2022 Semester 1
 * @param {string} name  e.g. "Software Engineering & Object-Oriented Programming"
 */
export const createModule = (moduleCode, semester, name) => {
  axios.post(`/module`, { Code: moduleCode, Semester: semester, Name: name }).catch((err) => {
    console.error(err);
  });
};

/**
 * Adds a student to a module
 * 
 * @param {int} moduleId ID of module
 * @param {int} studentId ID of student
 */
export const createEnrollment = (moduleId, studentId) => {
  axios.post(`/module/enroll`, { ModuleID: moduleId, StudentID: studentId }).catch((err) => {
    console.error(err);
  });
};

/**
 * Adds a staff member to a module
 * 
 * @param {int} moduleId 
 * @param {int} staffId 
 */
export const createSupervision = (moduleId, staffId) => {
  axios.post(`/module/supervise`, { ModuleID: moduleId, StaffID: staffId }).catch((err) => {
    console.error(err);
  });
};

/**
 * Activates pairings which originate from a complete graph of each group, whose size is at least assignment.GroupSize - 1
 * 
 * @param {int} moduleId ID of module the assignment belongs to
 * @param {Object} assignmentData must be (Name, ModuleID) or (ID)
 */
export const assignPairings = (moduleId, assignmentData) => {
  axios.post(`/module/${moduleId}/pairing/assign`, assignmentData).catch((err) => {
    console.error(err);
  });
};

/**
 * Deletes old pairings of assignment (if any) and reinserts all possible pairings, marking them as inactive
 * 
 * @param {int} moduleId ID of module the assignment belongs to
 * @param {Object} assignmentData must be (Name + ModuleID) or (ID)
 */
export const initializePairings = (moduleId, assignmentData) => {
  axios.post(`/module/${moduleId}/pairing/initialize`, assignmentData).catch((err) => {
    console.error(err);
  });
};

/**
 * Returns Modules data that matches given data
 * 
 * @param {Object} moduleData can consist of attributes ID, Code, Name, and/or Semester
 */
export const getModules = (moduleData) => {
  axios.get(`/module`, {
    method: 'GET',
    body: JSON.stringify(moduleData)
  }).catch((err) => {
    console.log(err);
  });
};

/**
 * Returns Enrollments data that matches given data
 * 
 * @param {Object} enrollmentData can consist of ID, ModuleID, and/or StudentID
 */
export const getEnrollments = (enrollmentData) => {
  axios.get(`/module`, {
    method: 'GET',
    body: JSON.stringify(enrollmentData)
  }).catch((err) => {
    console.log(err);
  });
};

/**
 * Returns Supervisions data that matches given data
 * 
 * @param {Object} enrollmentData can consist of ID, ModuleID, and/or StaffID
 */
export const getSupervisions = (supervisionData) => {
  axios.get(`/module`, {
    method: 'GET',
    body: JSON.stringify(supervisionData)
  }).catch((err) => {
    console.log(err);
  });
};

/*
TODO Assignment actions
create/get assignments
create/get questions
create/get rubrics
*/

// TODO remove invalid functions below

/*
export const listModules = (updateModuleList) => {
  axios
    .get("/module/list")
    .then((res) => {
      updateModuleList(res.data);
    })
    .catch((err) => {
      console.error(err);
    });
};

export const getModuleInfo = (moduleCode, updateModuleInfo) => {
  axios
    .get(`/module/${moduleCode}`)
    .then((res) => {
      updateModuleInfo(res.data);
    })
    .catch((err) => {
      if (err.response.status === 503) {
        createModule(moduleCode);
        getModuleInfo(moduleCode, updateModuleInfo);
      }
    });
};

export const addUserToModule = (moduleCode, userHandle) => {
  axios.get(`/module/${moduleCode}/add/${userHandle}`).catch((err) => {
    console.error(err);
  });
};

export const updateGroupings = (moduleCode, groupings) => {
  axios.post(`/module/${moduleCode}/groups/update`, groupings).catch((err) => {
    console.error(err);
  });
};

export const updateAssignments = (moduleCode, assignments) => {
  axios
    .post(`/module/${moduleCode}/assignments/update`, assignments)
    .catch((err) => {
      console.error(err);
    });
};

export const updateQuestion = (moduleCode, newQuestion) => {
  axios.post(`/module/${moduleCode}/version/${newQuestion}`).catch((err) => {
    console.error(err);
  });
};
*/
