import { Component, FC, FormEvent, useContext, useState } from "react";
import { signupUser } from "../actions/userActions";
import { AuthContext } from "../context/context";
import { Link, useHistory } from "react-router-dom";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";
import InputLabel from "@material-ui/core/InputLabel";
import "./Login.css";
import { User } from "../models/models";
import { Role } from "./Login";

export const SignUp: FC = (props) => {
  const [user, setUser] = useState<User>({});
  const [role, setRole] = useState<Role>(Role.STUDENT);
  const { dispatch } = useContext(AuthContext);
  const history = useHistory();

  const submitSignup = (e: FormEvent) => {
    e.preventDefault();
    try {
      signupUser(role, user, history)(dispatch);
    } catch (err) {
      console.error(err);
    }
  }

  return (
    <div>
      <Grid container spacing={0} justify="center" direction="row">
        <Grid item>
          <Grid
            container
            direction="column"
            justify="center"
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
                <form onSubmit={(e) => submitSignup(e)}>
                  <Grid container direction="column" spacing={2}>
                    <Grid item>
                      <TextField
                        type="Name"
                        placeholder="Name"
                        fullWidth
                        name="handle"
                        variant="outlined"
                        onChange={(e) =>
                          setUser({ ...user, Name: e.target.value })
                        }
                        required
                        autoFocus
                      />
                    </Grid>
                    <Grid item>
                      <TextField
                        type="email"
                        placeholder="Email"
                        fullWidth
                        name="email"
                        variant="outlined"
                        onChange={(e) =>
                          setUser({ ...user, Email: e.target.value })
                        }
                        required
                      />
                    </Grid>
                    <Grid item>
                      <TextField
                        type="password"
                        placeholder="Password"
                        fullWidth
                        name="password"
                        variant="outlined"
                        onChange={(e) =>
                          setUser({ ...user, Password: e.target.value })
                        }
                        required
                      />
                    </Grid>
                    <Grid item>
                      <InputLabel>Role</InputLabel>
                      <Select
                        name="status"
                        defaultValue={Role.STUDENT}
                        fullWidth
                        onChange={(e) =>
                          setRole(e.target.value as Role)
                        }
                      >
                        <MenuItem value={Role.STUDENT}>Student</MenuItem>
                        <MenuItem value={Role.STAFF}>Staff</MenuItem>
                      </Select>
                    </Grid>
                    <Grid item>
                      <Button
                        variant="contained"
                        color="primary"
                        type="submit"
                        className="button-block"
                      >
                        Sign Up
                      </Button>
                    </Grid>
                    <Grid item>
                      <Link to="/login">
                        <Button
                          variant="contained"
                          color="primary"
                          className="button-block"
                        >
                          Login
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
}
