import { Switch, Route } from "react-router-dom";

//components
import { Login } from "./pages/Login";
// import SignUp from "./pages/signup";
import { Home, StaffHome, StudentHome } from "./pages/Home";
import { StaffModuleDashboard } from "./pages/Dashboard/Staff/Dashboard";
import { StudentModuleDashboard } from "./pages/Dashboard/Student/Dashboard";
import { FC } from "react";
import { Class } from "./pages/Dashboard/Staff/Class";
import { Assignments } from "./pages/Dashboard/Staff/Assignments";

const Routes: FC = () => (
  <Switch>
    <Route exact path="/" component={Home} />
    <Route exact path="/login" component={Login} />

    <Route exact path="/staff" component={StaffHome} />
    <Route exact path="/staff/module/:moduleId" component={StaffModuleDashboard} />
    <Route exact path="/staff/module/:moduleId/class" component={Class} />
    <Route exact path="/staff/module/:moduleId/assignments" component={Assignments} />

    <Route exact path="/student" component={StudentHome} />
    <Route exact path="/student/module/:moduleId" component={StudentModuleDashboard} />
  </Switch>
);

export default Routes;
