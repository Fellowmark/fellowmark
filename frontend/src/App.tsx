/* eslint-disable react/jsx-no-undef */
import "./App.css";
import axios from "axios";
import Routes from "./Routes";
import { useContext, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { AuthContext, TimeoutContext } from "./context/context";
import { authenticate } from "./utils/auth";
import { logoutUser, setAuthorizationHeader } from "./actions/userActions";
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from "@material-ui/core";

axios.defaults.baseURL =
  process.env.REACT_APP_API_URL ?? "http://localhost:5050/api";

const useAuthHook = () => {
  const { dispatch } = useContext(AuthContext);
  const { createTimeout, cancelTimeout } = useContext(TimeoutContext);
  const history = useHistory();

  useEffect(() => {
    const secondsLeft = authenticate(dispatch);
    if (secondsLeft) {
      cancelTimeout();
      createTimeout(
        setTimeout(() => {
          logoutUser(history, dispatch);
          alert("Session expired");
        }, secondsLeft - 30)
      );
      setAuthorizationHeader(localStorage.jwt);
    } else {
      logoutUser(history, dispatch);
      cancelTimeout();
    }
  }, []);
};

function App() {
  const [open, setOpen] = useState(true);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  useAuthHook();

  return (
    <div className="App container">
      <>
        <Dialog
          open={open}
          onClose={handleClose}
          aria-labelledby="alert-dialog-title"
          aria-describedby="alert-dialog-description"
        >
          <DialogTitle id="alert-dialog-title">{"Privacy Notice"}</DialogTitle>
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              This site uses cookies. By clicking accept or continuing to
              use this site, you agree to our use of cookies. For more details,
              please see our{" "}
              <a href="https://www.nus.edu.sg/privacy-notice/">
                Privacy Policy
              </a>
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>
              Agree
            </Button>
          </DialogActions>
        </Dialog>
      </>
      <Routes />
    </div>
  );
}

export default App;
