import { useContext, useEffect, useState } from "react";
import axios from "axios";

import { authenticate } from "../utils/auth";
import { logoutUser } from "../actions/userActions";
import { AuthContext } from "../context/context";
import { Redirect, useHistory } from "react-router-dom";
import { Role } from "./Login";
import { ProgressBar } from "../components/ProgressBar";
import { AuthType } from "../reducers/reducer";

export const Home: React.FC = () => {
  const { state, dispatch } = useContext(AuthContext);
  const [isLoaded, setIsLoaded] = useState(false);
  const [role, setRole] = useState("");
  const history = useHistory();

  useEffect(() => {
      dispatch({ type: AuthType.AUTHENTICATED, payload: {user: "lol"}});
    if (authenticate(dispatch)) {
      axios.defaults.headers.common["Authorization"] = localStorage.FBIdToken;
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
  return <ProgressBar component={showComponent} isLoaded={isLoaded} />;
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

  const userComponent = <div></div>;

  return <ProgressBar component={userComponent} isLoaded={isLoaded} />;
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

  return <ProgressBar component={userComponent} isLoaded={isLoaded} />;
};
