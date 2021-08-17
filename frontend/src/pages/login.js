import React, { Component } from "react";
import { Link } from "react-router-dom";

import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";

import { Auth } from "../context/context";
import { loginUser } from "../actions/userActions";

import "./login.css";

class Login extends Component {
    constructor(props) {
        super(props);
        this.state = {
            email: "",
            password: ""
        };
    }

    submitLogin(e) {
        e.preventDefault();
        const { _, dispatch } = this.context;
        const userDetails = { email: this.state.email, password: this.state.password };
        console.log(userDetails);
        try {
            loginUser(userDetails, this.props.history)(dispatch);
        } catch (err) {
            console.error(err);
        }
    }

    render() {
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
                                    <form onSubmit={(e) => this.submitLogin(e)}>
                                        <Grid container direction="column" spacing={2}>
                                            <Grid item>
                                                <TextField
                                                    type="email"
                                                    placeholder="Email"
                                                    fullWidth
                                                    name="email"
                                                    variant="outlined"
                                                    onChange={(e) => this.setState({ email: e.target.value })}
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
                                                    value={this.state.password}
                                                    onChange={(e) => this.setState({ password: e.target.value })}
                                                    required
                                                />
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
                                                    <Button variant="contained" color="primary" className="button-block" >
                                                        Sign Up
                                                    </Button>
                                                </Link>
                                            </Grid>
                                        </Grid>
                                    </form>
                                </Grid>
                                <Grid item>
                                </Grid>
                            </Paper>
                        </Grid>
                    </Grid>
                </Grid>
            </div>
        );
    }
}

Login.contextType = Auth;

export default login;
