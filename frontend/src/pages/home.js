import { Component } from "react";
import jwtDecode from "jwt-decode";

import CircularProgress from "@material-ui/core/CircularProgress";

import { logoutUser, getUserDetails } from "../actions/userActions";
import { Auth } from "../context/authContext";

class Home extends Component {
    constructor(props) {
        super(props);
        this.state = {
            moduleId: localStorage.getItem('module'),
            isAdmin: false,
            isLoaded: false
        };
    }

    updateModuleInfo(moduleInfo) {
        this.setState(moduleInfo);
        this.setState({ isLoaded: true });
    }

    logout() {
        const { dispatch } = this.context;
        logoutUser()(dispatch);
    }

    componentDidMount() {
        console.log("Welcome");
        const token = localStorage.FBIdToken;
        const { dispatch } = this.context;
        if (token) {
            const decodedToken = jwtDecode(token);
            if (decodedToken.exp * 1000 < Date.now()) {
                dispatch(logoutUser());
                window.location.href = "/login";
            } else {
                axios.defaults.headers.common["Authorization"] = token;
                getUserDetails()((newState) => {
                    dispatch(newState);
                    const isAdmin = newState.payload.credentials.status == "Module Admin" || newState.payload.credentials.status == "Tutor";
                    const studentAdded =
                        newState.payload.credentials[
                        newState.payload.credentials.moduleId
                        ] &&
                        newState.payload.credentials[
                        newState.payload.credentials.moduleId
                        ] != "unassigned";

                    if (!isAdmin && !studentAdded) {
                        alert("Ask module admin to assign you");
                        logoutUser()(dispatch);
                        window.location.href = "/login";
                    }

                    this.setState({
                        isAdmin: newState.payload.credentials.status == "Module Admin", // TODO check on this part
                        moduleId: newState.payload.credentials.moduleId,
                    });
                    getModuleInfo(
                        newState.payload.credentials.moduleId,
                        this.updateModuleInfo
                    );
                });
            }
        } else {
            logoutUser()(dispatch);
            window.location.href = "/login";
        }
    }

    render() {
        const userComponent = (
            <div>
                // TODO fill with Admin and Student components
            </div>
        );
        const showComponent = this.state.isLoaded ? (
            userComponent
        ) : (
            <CircularProgress className="pro" />
        );
        return <div>{showComponent}</div>;
    }
}

Home.contextType = Auth;
