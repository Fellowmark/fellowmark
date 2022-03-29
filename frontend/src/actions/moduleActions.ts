import axios from 'axios';
import { Dispatch, SetStateAction } from 'react';
import { Assignment, Enrollment, Grade, Pairing, Question, Rubric } from '../models/models';
import { Pagination } from '../models/pagination';

/**
 * Creates module using given data
 *
 * @param {string} moduleCode e.g. "CS2113T"
 * @param {string} semester e.g. "2122-1" for AY2021/2022 Semester 1
 * @param {string} name  e.g. "Software Engineering & Object-Oriented Programming"
 */
export const createModule = (moduleCode: string, semester: string, name: string) => {
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
export const createEnrollment = (moduleId: number, studentId: number) => {
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
export const createSupervision = (moduleId: number, staffId: number) => {
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
export const assignPairings = async (assignmentData: unknown) => {
  await axios.post(`/assignment/pairs/assign`, assignmentData);
};

export const getPairingAsReviewee = ({ assignmentId }, setPairings?: Dispatch<SetStateAction<Pagination<Pairing>>>) => {
  axios
    .get(`/assignment/pairs/mymarkers`, {
      params: {
        assignmentId,
      },
    })
    .then((res) => {
      setPairings(res.data);
    })
    .catch((err) => {
      console.error(err);
    });
};

export const getPairingAsMarker = ({ assignmentId }, setPairings: Dispatch<SetStateAction<Pagination<Pairing>>>) => {
  axios
    .get(`/assignment/pairs/myreviewees`, {
      params: {
        assignmentId: assignmentId,
      },
    })
    .then((res) => {
      setPairings(res.data);
    })
    .catch((err) => {
      console.error(err);
    });
};

export const getAllPairings = ({ assignmentId }, setPairings: Dispatch<SetStateAction<Pagination<Pairing>>>) => {
  axios
    .get(`/assignment/pairs`, {
      params: {
        assignmentId: assignmentId,
      },
    })
    .then((res) => {
      setPairings(res.data);
    })
    .catch((err) => {
      console.error(err);
    });
};

/**
 * Deletes old pairings of assignment (if any) and reinserts all possible pairings, marking them as inactive
 *
 * @param {int} moduleId ID of module the assignment belongs to
 * @param {Object} assignmentData must be (Name + ModuleID) or (ID)
 */
export const initializePairings = async (assignmentData: unknown) => {
  await axios.post(`/assignment/pairs/initialize`, assignmentData);
};

export const createPairings = async (assignmentData: any) => {
  await initializePairings({
    id: assignmentData.AssignmentID,
  });
  await assignPairings({
    id: assignmentData.AssignmentID,
  });
};

/**
 * Returns Modules data that matches given data
 *
 * @param {Object} moduleData can consist of attributes ID, Code, Name, and/or Semester
 */
export const getModules = (moduleData: unknown, setModules: Dispatch<SetStateAction<any>>) => {
  axios
    .get(`/module`, {
      method: 'GET',
      params: moduleData,
    })
    .then((res) => {
      return setModules(res);
    })
    .catch((err) => {
      console.log(err);
    });
};

/**
 * Returns Enrollments data that matches given data
 *
 * @param {Object} enrollmentData can consist of ID, ModuleID, and/or StudentID
 */
export const getEnrollments = (enrollmentData: any, setEnrollments: Dispatch<SetStateAction<Pagination<Enrollment>>>) => {
  axios
    .get(`/module/enrolls`, {
      method: 'GET',
      params: enrollmentData,
    })
    .then((res) => {
      setEnrollments(res.data);
    })
    .catch((err) => {
      console.log(err);
    });
};

// /**
//  * Returns Enrollments data that matches given data
//  *
//  * @param {Object} enrollmentData can consist of ID, ModuleID, and/or StudentID
//  */
// export const getModu = (enrollmentData) => {
//   axios.get(`/module/enroll`, {
//     method: 'GET',
//     params: enrollmentData
//   }).then((res) => {
//     return res.data;
//   }).catch((err) => {
//     console.log(err);
//   });
// };

/**
 * Returns Enrollments data that matches given data
 *
 * @param {Object} enrollmentData can consist of ID, ModuleID, and/or StudentID
 */
export const getStudentModules = (setModules: Dispatch<SetStateAction<any>>) => {
  axios
    .get(`/module/enroll`)
    .then((res) => {
      return res.data;
    })
    .then((res) => {
      setModules(res);
    })
    .catch((err) => {
      console.log(err);
    });
};

/**
 * Returns Supervisions data that matches given data
 *
 * @param {Object} supervisionData can consist of ID, ModuleID, and/or StaffID
 */
export const getSupervisions = (setModules: Dispatch<SetStateAction<any>>) => {
  axios
    .get(`/module/supervises`)
    .then((res) => {
      return res.data;
    })
    .then((res) => {
      setModules(res.data);
    })
    .catch((err) => {
      console.log(err);
    });
};

/**
 * Returns Supervisions data that matches given data
 *
 * @param {Object} supervisionData can consist of ID, ModuleID, and/or StaffID
 */
export const getStaffModules = (setModules: Dispatch<SetStateAction<any>>) => {
  axios
    .get(`/module/supervise`)
    .then((res) => {
      return res.data;
    })
    .then((res) => {
      setModules(res);
    })
    .catch((err) => {
      console.log(err);
    });
};

/**
 * Creates assignment using given data
 *
 * @param {string} name name of the assignment to be created
 * @param {int} moduleId ModuleID of the module
 * @param {int} groupSize size of each group of student-marker pairings
 * @param {int} [duration=86400] time in seconds assignment should be open for submissions [default: `86400` (1 day)]
 */
export const createAssignment = async (assignment: Assignment) => {
  return await axios.post(`/assignment`, {
    name: assignment.Name,
    moduleId: assignment.ModuleID,
    groupSize: assignment.GroupSize,
    deadline: assignment.Deadline,
  });
};

/**
 * Creates question using given data
 *
 * @param {int} questionNumber question number in the context of the assignment
 * @param {string} questionText question text
 * @param {int} assignmentId AssignmentID of corresponding assignment
 */
export const createQuestion = async (questionNumber: number, questionText: string, assignmentId: number) => {
  await axios.post(`/assignment/question`, {
    questionNumber: questionNumber,
    questionText: questionText,
    assignmentId: assignmentId,
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
export const createRubrics = async ({ QuestionID, Criteria, Description, MaxMark = 10, MinMark = 0 }: any) => {
  await axios.post(`/assignment/rubric`, {
    questionId: QuestionID,
    criteria: Criteria,
    description: Description,
    minMark: MinMark,
    maxMark: MaxMark,
  });
};

/**
 * Returns Assignments data that matches given data
 *
 * @param {Object} assignmentData can consist of AssignmentID, Name, ModuleID, GroupSize and/or Deadline
 */
export const getAssignments = async (
  assignmentData: any,
  setAssignments?: Dispatch<SetStateAction<Pagination<Assignment>>>
): Promise<Assignment[]> => {
  const { data } = await axios.get(`/assignment`, {
    method: 'GET',
    params: {
      ...assignmentData,
      sort: 'deadline asc',
    },
  });
  if (setAssignments) setAssignments(data);
  return data;
};

/**
 * Returns Questions data that matches given data
 *
 * @param {Object} questionData can consist of QuestionID, QuestionNumber, QuestionText and/or AssignmentID
 */
export const getQuestions = async (
  questionData: unknown,
  setQuestions?: Dispatch<SetStateAction<Pagination<Question>>>
): Promise<Question[]> => {
  const { data } = await axios.get(`/assignment/question`, {
    params: questionData,
  });
  if (setQuestions) setQuestions(data);
  return data;
};

/**
 * Returns Rubrics data that matches given data
 *
 * @param {Object} rubricData can consist of RubricID, QuestionID, Criteria, Description, MinMark and/or MaxMark
 */
export const getRubrics = async (rubricData: Rubric, setRubrics?: Dispatch<SetStateAction<Pagination<Rubric>>>) => {
  const { data } = await axios.get(`/assignment/rubric`, {
    params: {
      questionId: rubricData.QuestionID,
    },
  });
  if (setRubrics) setRubrics(data);
  return data;
};

export const uploadSubmission = async (fileFormData: unknown, moduleId: number, questionId: number, studentId: number) => {
  await axios.post(`/submission?questionId=${questionId}&studentId=${studentId}`, fileFormData);
};

export const downloadSubmission = async (moduleId: number, questionId: number, studentId: number) => {
  try {
    const res = await axios.get(`/submission`, {
      params: {
        questionId: questionId,
        studentId: studentId,
      },
      responseType: 'blob', // important
    });
    let blob = new Blob([res.data], { type: 'application/octet-stream' });
    return URL.createObjectURL(blob);
  } catch (e) {
    alert('No submission found');
  }
};

export const getSubmissionMetadata = (studentId: number, questionId: number, setSubmission: Dispatch<SetStateAction<boolean>>) => {
  axios
    .get(`/assignment/submission`, {
      params: {
        studentId: studentId,
        questionId: questionId,
      },
    })
    .then((res) => {
      if (res.data.totalRows) {
        setSubmission(true);
      } else {
        setSubmission(false);
      }
    });
};

export const createGrade = (pairingId: number, rubricId: number, grade: Grade) => {
  // const rubric = getRubrics({ RubricID: rubricId })[0];
  // if (grade < rubric.MinMark || grade > rubric.MaxMark) {
  //   console.error('Please provide a valid grade');
  //   return;
  // }
  axios
    .post(`/assignment/rubric`, {
      PairingID: pairingId,
      RubricID: rubricId,
      Grade: grade,
    })
    .catch((err) => {
      console.error(err);
    });
};

export const getGradesForStudent = async (moduleId: number, gradeData: Grade, setGrades?: Dispatch<SetStateAction<Map<number, Grade>>>) => {
  const { data } = await axios.get(`/grade/my/reviewee`, {
    params: {
      pairingId: gradeData.PairingID,
    },
  });
  let grades = new Map();
  data.rows.forEach((grade: Grade) => {
    grades.set(grade.RubricID, grade);
  });
  if (setGrades) setGrades(grades);
  return grades;
};

export const getGradesForMarker = (moduleId: number, gradeData: Grade, setGrades: Dispatch<SetStateAction<Map<number, Grade>>>) => {
  axios
    .get(`grade/my/marker`, {
      method: 'GET',
      params: {
        pairingId: gradeData.PairingID,
      },
    })
    .then((res) => {
      let grades = new Map();
      res.data.rows.forEach((grade) => {
        grades.set(grade.RubricID, grade);
      });
      setGrades(grades);
    })
    .catch((err) => {
      console.log(err);
    });
};

export const postGrade = async (moduleId: number, gradeData: Grade) => {
  await axios.post(`/grade`, {
    pairingId: gradeData.PairingID,
    rubricId: gradeData.RubricID,
    comment: gradeData.Comment,
    grade: gradeData.Grade,
  });
};

export const getGradesForStaff = async (gradeData: Grade) => {
  const { data } = await axios.get('/grading');
  return data;
};

export const getAssignmentGradesForStudent = async (assignmentData: Assignment) => {
  let assignment = await getAssignments(assignmentData)[0];
  let questions = await getQuestions({ AssignmentID: assignment.ID });
  let questionsGrades = {};
  questions.forEach((question) => {
    let rubric = getRubrics({ QuestionID: question.ID })[0];
    questionsGrades[question.ID] = getGradesForStudent(null, { RubricID: rubric.ID });
  });
  return questionsGrades;
};

export const getQuestionGradesForStudent = (questionData: Question) => {
  let question = getQuestions(questionData)[0];
  let rubric = getRubrics({ QuestionID: question.ID })[0];
  return getGradesForStudent(null, { RubricID: rubric.ID });
};
