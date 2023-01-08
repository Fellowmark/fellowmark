import { Switch, Route } from "react-router-dom";

//components
import { Login, Role } from "./pages/Login";
import { Home, RoleHome } from "./pages/Home";
import { StaffModuleDashboard } from "./pages/Dashboard/Staff/Dashboard";
import { StudentModuleDashboard } from "./pages/Dashboard/Student/Dashboard";
import { TAModuleDashboard } from "./pages/Dashboard/TA/Dashboard";
import { AdminModuleDashboard } from "./pages/Dashboard/Admin/Dashboard";
import { StaffManagement } from "./pages/Dashboard/Admin/StaffMangement";
import { FC } from "react";
import { Class as StaffClass } from "./pages/Dashboard/Staff/Class";
import { Supervisors as StaffSupervisors } from "./pages/Dashboard/Staff/Supervisors";
import { Supervisors as AdminSupervisors } from "./pages/Dashboard/Admin/Supervisors";
import { TAs as StaffTAs } from "./pages/Dashboard/Staff/TAs";
import { TAs as AdminTAs } from "./pages/Dashboard/Admin/TAs";
import { Class as AdminClass } from "./pages/Dashboard/Admin/Class";
import { Class as TAClass } from "./pages/Dashboard/TA/Class";
import { Assignments as StaffAssignments } from "./pages/Dashboard/Staff/Assignments";
import { Assignments as StudentAssignments} from "./pages/Dashboard/Student/Assignments";
import { Assignments as AdminAssignments} from "./pages/Dashboard/Admin/Assignments";
import { Questions as StaffQuestions } from "./pages/Dashboard/Staff/Questions";
import { Questions as StudentQuestions } from "./pages/Dashboard/Student/Questions";
import { Questions as AdminQuestions } from "./pages/Dashboard/Admin/Questions";
import { SignUp } from "./pages/Signup";
import { QuestionBoard as StudentQuestionBoard } from "./pages/Dashboard/Student/QuestionBoard";
import { QuestionBoard as StaffQuestionBoard } from "./pages/Dashboard/Staff/QuestionBoard";
import { QuestionBoard as AdminQuestionBoard } from "./pages/Dashboard/Admin/QuestionBoard";

const Routes: FC = () => {
  return (
    <Switch>
      <Route exact path="/" component={Home} />
      <Route exact path="/login" component={Login} />
      <Route exact path="/signup" component={SignUp} />

      <Route exact path="/staff" component={() => <RoleHome role={Role.STAFF} />} />
      <Route exact path="/staff/module/:moduleId" component={StaffModuleDashboard} />
      <Route exact path="/staff/module/:moduleId/class" component={StaffClass} />
      <Route exact path="/staff/module/:moduleId/supervisors" component={StaffSupervisors} />
      <Route exact path="/staff/module/:moduleId/tas" component={StaffTAs} />
      <Route exact path="/staff/module/:moduleId/assignments" component={StaffAssignments} />
      <Route exact path="/staff/module/:moduleId/assignments/:assignmentId" component={StaffQuestions} />
      <Route exact path="/staff/module/:moduleId/assignments/:assignmentId/question/:questionId" component={StaffQuestionBoard} />

      <Route exact path="/student" component={() => <RoleHome role={Role.STUDENT} />} />
      <Route exact path="/student/module/:moduleId" component={StudentModuleDashboard} />
      <Route exact path="/student/module/:moduleId/assignments" component={StudentAssignments} />
      <Route exact path="/student/module/:moduleId/assignments/:assignmentId" component={StudentQuestions} />
      <Route exact path="/student/module/:moduleId/assignments/:assignmentId/question/:questionId" component={StudentQuestionBoard} />

      <Route exact path="/student/ta" component={() => <RoleHome role={Role.TA} />} />
      <Route exact path="/student/ta/module/:moduleId" component={TAModuleDashboard} />
      <Route exact path="/student/ta/module/:moduleId/class" component={TAClass} />


      <Route exact path="/admin" component={() => <RoleHome role={Role.ADMIN} />} />
      <Route exact path="/admin/managestaff" component={StaffManagement} />
      <Route exact path="/admin/module/:moduleId" component={AdminModuleDashboard} />
      <Route exact path="/admin/module/:moduleId/class" component={AdminClass} />
      <Route exact path="/admin/module/:moduleId/supervisors" component={AdminSupervisors} />
      <Route exact path="/admin/module/:moduleId/tas" component={AdminTAs} />
      <Route exact path="/admin/module/:moduleId/assignments" component={AdminAssignments} />
      <Route exact path="/admin/module/:moduleId/assignments/:assignmentId" component={AdminQuestions} />
      <Route exact path="/admin/module/:moduleId/assignments/:assignmentId/question/:questionId" component={AdminQuestionBoard} />
    </Switch>
  );
};

export default Routes;
