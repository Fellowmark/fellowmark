import React, { Component } from "react";
import { signupUser } from "../actions/userActions";
import { listModules } from "../actions/moduleActions";
import { Auth } from "../context/authContext";
import { Link } from "react-router-dom";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import Select from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';
import InputLabel from '@material-ui/core/InputLabel';
import "./login.css";
import {Menu} from "@material-ui/core";


class signup extends Component {
    constructor(props) {
        super(props);
        this.updateModuleList = this.updateModuleList.bind(this);
        this.state = {
            moduleList: [],
            status: "Module Admin",
            handle: "",
            moduleCode: "CS2113T",
            email: "",
            password: ""
        }
    }

    componentWillMount() {
        listModules(this.updateModuleList);
    }

    updateModuleList(modules) {
        this.setState({ moduleList: modules });
    }

    submitSignup (e) {
        e.preventDefault();
        const {_, dispatch} = this.context;
        const userDetails = {email: this.state.email, password: this.state.password, status: this.state.status, handle: this.state.handle, moduleCode: this.state.moduleCode};
        console.log(userDetails);
        try {
            signupUser(userDetails, this.props.history)(dispatch);
        } catch (err) {
            console.error(err);
        }
    }

    render () {
        const moduleOptions = ["CS2113T"].map(code => {
            return <MenuItem value={code}>{code}</MenuItem>
        });

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
                                        Welcome to NUSColab!
                                    </Typography>
                                </Grid>
                                <Grid item>
                                    <form onSubmit={(e) => this.submitSignup(e)}>
                                        <Grid container direction="column" spacing={2}>
                                            <Grid item>
                                                <TextField
                                                    type="handle"
                                                    placeholder="Username"
                                                    fullWidth
                                                    name="handle"
                                                    variant="outlined"
                                                    onChange={(e) => this.setState({handle: e.target.value})}
                                                    required
                                                    autoFocus
                                                />
                                            </Grid>
                                            <Grid item>
                                                <TextField
                                                    type="email"
                                                    placeholder="Email"
                                                    fullWidth
                                                    name="username"
                                                    variant="outlined"
                                                    onChange={(e) => this.setState({email: e.target.value})}
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
                                                    onChange={(e) => this.setState({password: e.target.value})}
                                                    required
                                                />
                                            </Grid>
                                            <Grid item>
                                                <InputLabel>Module</InputLabel>
                                                <Select name="moduleCode" fullWidth onChange={(e) => this.setState({moduleCode: e.target.value})}>{moduleOptions}</Select>
                                            </Grid>
                                            <Grid item>
                                                <InputLabel>Status</InputLabel>
                                                <Select name="status" fullWidth onChange={(e) => this.setState({status: e.target.value})}>
                                                    <MenuItem value="Module Admin">Module Admin</MenuItem>
                                                    <MenuItem value="Student">Student</MenuItem>
                                                    <MenuItem value="Tutor">Tutor</MenuItem>
                                                </Select>
                                            </Grid>
                                            <Grid item>
                                                <Button
                                                    variant="contained"
                                                    color="primary"
                                                    type="submit"
                                                    className="button-block"
                                                >
                                                    Sign up
                                                </Button>
                                            </Grid>
                                            <Grid item>
                                            <Link to="/login">
                                                <Button variant="contained" color="primary" className="button-block" > 
                                                    Login 
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
        )
    }
}

signup.contextType = Auth;

export default signup;
