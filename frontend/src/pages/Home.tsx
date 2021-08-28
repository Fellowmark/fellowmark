import { useContext, useEffect, useState } from "react";

import { authenticate } from "../utils/auth";
import { logoutUser, setAuthorizationHeader } from "../actions/userActions";
import { AuthContext } from "../context/context";
import { Redirect, useHistory } from "react-router-dom";
import { Role } from "./Login";
import { ProgressBar } from "../components/ProgressBar";
import { AuthType } from "../reducers/reducer";
import { ModuleList } from "./Modules";
import { Grid } from "@material-ui/core";

export const Home: React.FC = () => {
  const { state, dispatch } = useContext(AuthContext);
  const [isLoaded, setIsLoaded] = useState(false);
  const [role, setRole] = useState("");
  const history = useHistory();

  useEffect(() => {
    dispatch({ type: AuthType.AUTHENTICATED, payload: { user: "lol" } });
    if (authenticate(dispatch)) {
      setAuthorizationHeader(localStorage.FBIdToken);
      setIsLoaded(true);
    } else {
      logoutUser(history, dispatch);
    }
  }, []);

  useEffect(() => {
    if (state) {
      setRole(state.role);
    }
  }, [state]);

  const showComponent = <Redirect to={`/${role.toLowerCase()}`} />;
  return (
    <Grid container>
      <ProgressBar component={showComponent} isLoaded={isLoaded} />
    </Grid>
  );
};

export const StudentHome: React.FC = () => {
  const [isLoaded, setIsLoaded] = useState(false);
  const { state } = useContext(AuthContext);
  const history = useHistory();

  useEffect(() => {
    if (state?.role !== Role.STUDENT) {
      history.replace("/");
    } else {
      setIsLoaded(true);
    }
  }, []);

  const userComponent = <ModuleList />;

  return (
    <Grid container>
      <ProgressBar component={userComponent} isLoaded={isLoaded} />
    </Grid>
  );
};

export const StaffHome: React.FC = () => {
  const [isLoaded, setIsLoaded] = useState(false);
  const { state } = useContext(AuthContext);
  const history = useHistory();

  useEffect(() => {
    if (state?.role !== Role.STAFF) {
      history.replace("/");
    } else {
      setIsLoaded(true);
    }
  }, []);

  const userComponent = <div></div>;

  return (
    <Grid container>
      <ProgressBar component={userComponent} isLoaded={isLoaded} />
    </Grid>
  );
};
