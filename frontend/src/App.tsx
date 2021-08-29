import "./App.css";
import axios from "axios";
import Routes from "./Routes";
import { useContext, useEffect } from "react";
import { useHistory } from "react-router-dom";
import { AuthContext } from "./context/context";
import { authenticate } from "./utils/auth";
import { logoutUser, setAuthorizationHeader } from "./actions/userActions";

axios.defaults.baseURL = "http://localhost:5000";

function App() {
  const { dispatch } = useContext(AuthContext);
  const history = useHistory();

  useEffect(() => {
    if (authenticate(dispatch)) {
      setAuthorizationHeader(localStorage.FBIdToken);
    } else {
      logoutUser(history, dispatch);
    }
  }, []);

  return (
    <div className="App container">
      <Routes />
    </div>
  );
}

export default App;
