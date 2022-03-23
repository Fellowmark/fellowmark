import { Switch, Route } from "react-router-dom";

//components
import { Login, Role } from "./pages/Login";
import { Home, RoleHome } from "./pages/Home";
import { StaffModuleDashboard } from "./pages/Dashboard/Staff/Dashboard";
import { StudentModuleDashboard } from "./pages/Dashboard/Student/Dashboard";
import { AdminModuleDashboard } from "./pages/Dashboard/Admin/Dashboard";
import { StaffManagement } from "./pages/Dashboard/Admin/StaffMangement";
import { FC } from "react";
import { Class as StaffClass } from "./pages/Dashboard/Staff/Class";
import { Supervisors as StaffSupervisors } from "./pages/Dashboard/Staff/Supervisors";
import { Class as AdminClass } from "./pages/Dashboard/Admin/Class";
import { Assignments as StaffAssignments } from "./pages/Dashboard/Staff/Assignments";
import { Assignments as StudentAssignments} from "./pages/Dashboard/Student/Assignments";
import { Questions as StaffQuestions } from "./pages/Dashboard/Staff/Questions";
import { Questions as StudentQuestions } from "./pages/Dashboard/Student/Questions";
import { SignUp } from "./pages/Signup";
import { QuestionBoard as StudentQuestionBoard } from "./pages/Dashboard/Student/QuestionBoard";
import { QuestionBoard as StaffQuestionBoard } from "./pages/Dashboard/Staff/QuestionBoard";
import { Supervisors } from "./pages/Dashboard/Admin/Supervisors";

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
      <Route exact path="/staff/module/:moduleId/assignments" component={StaffAssignments} />
      <Route exact path="/staff/module/:moduleId/assignments/:assignmentId" component={StaffQuestions} />
      <Route exact path="/staff/module/:moduleId/assignments/:assignmentId/question/:questionId" component={StaffQuestionBoard} />

      <Route exact path="/student" component={() => <RoleHome role={Role.STUDENT} />} />
      <Route exact path="/student/module/:moduleId" component={StudentModuleDashboard} />
      <Route exact path="/student/module/:moduleId/assignments" component={StudentAssignments} />
      <Route exact path="/student/module/:moduleId/assignments/:assignmentId" component={StudentQuestions} />
      <Route exact path="/student/module/:moduleId/assignments/:assignmentId/question/:questionId" component={StudentQuestionBoard} />

      <Route exact path="/admin" component={() => <RoleHome role={Role.ADMIN} />} />
      <Route exact path="/admin/managestaff" component={StaffManagement} />
      <Route exact path="/admin/module/:moduleId" component={AdminModuleDashboard} />
      <Route exact path="/admin/module/:moduleId/class" component={AdminClass} />
      <Route exact path="/admin/module/:moduleId/supervisors" component={Supervisors} />
    </Switch>
  );
};

export default Routes;
