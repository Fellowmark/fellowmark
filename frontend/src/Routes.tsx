import { Switch, Route } from "react-router-dom";

//components
import { Login, Role } from "./pages/Login";
// import SignUp from "./pages/signup";
import { Home, RoleHome } from "./pages/Home";
import { StaffModuleDashboard } from "./pages/Dashboard/Staff/Dashboard";
import { StudentModuleDashboard } from "./pages/Dashboard/Student/Dashboard";
import { FC } from "react";
import { Class } from "./pages/Dashboard/Staff/Class";
import { Assignments } from "./pages/Dashboard/Staff/Assignments";
import { Questions } from "./pages/Dashboard/Staff/Assignment";
import { SignUp } from "./pages/Signup";

const Routes: FC = () => {
  return (
    <Switch>
      <Route exact path="/" component={Home} />
      <Route exact path="/login" component={Login} />
      <Route exact path="/signup" component={SignUp} />

      <Route exact path="/staff" component={() => <RoleHome role={Role.STAFF} />} />
      <Route exact path="/staff/module/:moduleId" component={StaffModuleDashboard} />
      <Route exact path="/staff/module/:moduleId/class" component={Class} />
      <Route exact path="/staff/module/:moduleId/assignments" component={Assignments} />
      <Route exact path="/staff/module/:moduleId/assignments/:assignmentId" component={Questions} />

      <Route exact path="/student" component={() => <RoleHome role={Role.STUDENT} />} />
      <Route exact path="/student/module/:moduleId" component={StudentModuleDashboard} />
    </Switch>
  );
};

export default Routes;
