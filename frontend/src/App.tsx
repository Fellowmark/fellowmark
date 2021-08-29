import "./App.css";
import axios from "axios";
import Routes from "./Routes";
import { useContext, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { AuthContext, TimeoutContext } from "./context/context";
import { authenticate } from "./utils/auth";
import { logoutUser, setAuthorizationHeader } from "./actions/userActions";

axios.defaults.baseURL = "http://localhost:5000";

const useAuthHook = () => {
  const { dispatch } = useContext(AuthContext);
  const { createTimeout, cancelTimeout } = useContext(TimeoutContext);
  const history = useHistory();

  useEffect(() => {
    const secondsLeft = authenticate(dispatch);
    if (secondsLeft) {
      cancelTimeout();
      createTimeout(setTimeout(() => {
        logoutUser(history, dispatch);;
        alert("Session expired");
      }, secondsLeft - 30));
      setAuthorizationHeader(localStorage.FBIdToken);
    } else {
      logoutUser(history, dispatch);
      cancelTimeout();
    }
  }, []);
}

function App() {
  useAuthHook();

  return (
    <div className="App container">
      <Routes />
    </div>
  );
}

export default App;
