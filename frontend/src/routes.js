import { Switch, Route } from "react-router-dom";

//components
import { Login } from "./pages/Login";
// import SignUp from "./pages/signup";
import { Home, StaffHome, StudentHome } from "./pages/Home";

const Routes = () => (
  <Switch>
    <Route exact path="/" component={Home} />
    <Route exact path="/student" component={StudentHome} />
    <Route exact path="/staff" component={StaffHome} />
    <Route exact path="/login" component={Login} />
  </Switch>
);

export default Routes;
