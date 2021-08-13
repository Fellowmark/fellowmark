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

/**
 * Creates assignment using given data
 * 
 * @param {string} name name of the assignment to be created
 * @param {int} moduleId ModuleID of the module
 * @param {int} groupSize size of each group of student-marker pairings
 * @param {int} [duration=86400] time in seconds assignment should be open for submissions [default: `86400` (1 day)]
 */
export const createAssignment = (name, moduleId, groupSize, duration = 86400) => {
  const deadline = Math.floor(Date.now() / 1000) + duration;
  axios.post(`/assignment`, { Name: name, ModuleID: moduleId, GroupSize: groupSize, Deadline: deadline }).catch((err) => {
    console.error(err);
  });
};

/**
 * Creates question using given data
 * 
 * @param {int} questionNumber question number in the context of the assignment
 * @param {string} questionText question text
 * @param {int} assignmentId AssignmentID of corresponding assignment
 */
export const createQuestion = (questionNumber, questionText, assignmentId) => {
  axios.post(`/assignment/question`, { QuestionNumber: questionNumber, QuestionText: questionText, AssignmentID: assignmentId }).catch((err) => {
    console.error(err);
  });
};

/**
 * Creates rubric for a question using given data
 * 
 * @param {int} questionId QuestionID of question rubric references
 * @param {string} criteria rubric criteria to be marked on
 * @param {string} description description of levels of correctness/answer quality and corresponding marks
 * @param {int} [maxMark=10] maximum amount of marks for the question (default: `10`)  
 * @param {int} [minMark=0] minimum amount of marks for the question (default: `0`) 
 */
export const createRubrics = (questionId, criteria, description, maxMark = 10, minMark = 0) => {
  axios.post(`/assignment/rubric`, {
    QuestionID: questionId,
    Criteria: criteria,
    Description: description,
    MinMark: minMark,
    MaxMark: maxMark
  }).catch((err) => {
    console.error(err);
  });
};

/**
 * Returns Assignments data that matches given data
 * 
 * @param {Object} assignmentData can consist of Name, ModuleID, GroupSize and/or Deadline
 */
export const getAssignments = (assignmentData) => {
  axios.get(`/assignment`, {
    method: 'GET',
    body: JSON.stringify(assignmentData)
  }).catch((err) => {
    console.log(err);
  });
};

/**
 * Returns Questions data that matches given data
 * 
 * @param {Object} questionData can consist of QuestionNumber, QuestionText and/or AssignmentID
 */
export const getQuestions = (questionData) => {
  axios.get(`/assignment/question`, {
    method: 'GET',
    body: JSON.stringify(questionData)
  }).catch((err) => {
    console.log(err);
  });
};

/**
 * Returns Rubrics data that matches given data
 * 
 * @param {Object} rubricData can consist of QuestionID, Criteria, Description, MinMark and/or MaxMark
 */
export const getRubrics = (rubricData) => {
  axios.get(`/assignment/rubric`, {
    method: 'GET',
    body: JSON.stringify(rubricData)
  }).catch((err) => {
    console.log(err);
  });
};


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
