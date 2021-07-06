import React from "react";
import { Switch, Route } from "react-router-dom";

//components
import login from "./pages/login";
import signup from "./pages/signup";
import home from "./pages/home";

const Routes = () => (
    <Switch>
        <Route exact path="/" component = {home} />
        <Route exact path="/login" component = {login} />   
        <Route exact path="/signup" component = {signup} />
    </Switch>
);

export default Routes;
