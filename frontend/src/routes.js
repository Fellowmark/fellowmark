import React from "react";
import { Switch, Route } from "react-router-dom";

//components
import Login from "./pages/login";
import SignUp from "./pages/signup";
import Home from "./pages/home";

const Routes = () => (
    <Switch>
        <Route exact path="/" component={Home} />
        <Route exact path="/login" component={Login} />
        <Route exact path="/signup" component={SignUp} />
    </Switch>
);

export default Routes;