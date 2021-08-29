import React, { FormEvent, useContext, useState } from "react";
import { Link, RouteComponentProps } from "react-router-dom";

import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";

import { AuthContext } from "../context/context";
import { loginUser } from "../actions/userActions";

import "./login.css";
import { Select } from "@material-ui/core";

export enum Role {
  STUDENT = "Student",
  STAFF = "Staff",
  ADMIN = "Admin",
}

export const Login: React.FC<RouteComponentProps> = (props) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState(Role.STUDENT);
  const { dispatch } = useContext(AuthContext);

  const submitLogin = (e: FormEvent) => {
    e.preventDefault();
    const userDetails = { email: email, password: password };
    try {
      loginUser(role, userDetails, props.history);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div>
      <Grid container spacing={0} justifyContent="center" direction="row">
        <Grid item>
          <Grid
            container
            direction="column"
            justifyContent="center"
            spacing={2}
            className="login-form"
          >
            <Paper
              variant="elevation"
              elevation={2}
              className="login-background"
            >
              <Grid item>
                <Typography component="h1" variant="h5">
                  Welcome to NUS Peermark!
                </Typography>
              </Grid>
              <Grid item>
                <form onSubmit={(e) => submitLogin(e)}>
                  <Grid container direction="column" spacing={2}>
                    <Grid item>
                      <TextField
                        type="email"
                        placeholder="Email"
                        fullWidth
                        name="email"
                        variant="outlined"
                        onChange={(e) => setEmail(e.target.value)}
                        required
                        autoFocus
                      />
                    </Grid>
                    <Grid item>
                      <TextField
                        type="password"
                        placeholder="Password"
                        fullWidth
                        name="password"
                        variant="outlined"
                        onChange={(e) => setPassword(e.target.value)}
                        required
                      />
                    </Grid>
                    <Grid item>
                      <Select
                        native
                        fullWidth
                        defaultValue="student"
                        name="role"
                        onChange={(e) => setRole(e.target.value as Role)}
                        required
                      >
                        <option value={"student"}>Student</option>
                        <option value={"staff"}>Staff</option>
                        <option value={"admin"}>Admin</option>
                      </Select>
                    </Grid>
                    <Grid item>
                      <Button
                        variant="contained"
                        color="primary"
                        type="submit"
                        className="button-block"
                      >
                        Login
                      </Button>
                    </Grid>
                    <Grid item>
                      <Link to="/signup">
                        <Button
                          variant="contained"
                          color="primary"
                          className="button-block"
                        >
                          Sign Up
                        </Button>
                      </Link>
                    </Grid>
                  </Grid>
                </form>
              </Grid>
              <Grid item></Grid>
            </Paper>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
};
