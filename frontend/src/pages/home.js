import { Component } from "react";
import { Auth } from "../context/authContext";
import jwtDecode from "jwt-decode";
import { logoutUser, getUserDetails } from "../actions/userActions";
import CircularProgress from "@material-ui/core/CircularProgress";
import {
  getModuleInfo,
  updateGroupings,
  updateAssignments,
  updateVersion,
} from "../actions/moduleActions";
import axios from "axios";
import Admin from "../components/adminUser/Admin";
import Student from "../components/studentUser/Student";

class home extends Component {
  constructor(props) {
    super(props);
    this.updateVersion = this.updateVersion.bind(this);
    this.updateModuleInfo = this.updateModuleInfo.bind(this);
    this.updateStudentGroup = this.updateStudentGroup.bind(this);
    this.changeAssignedExtract = this.changeAssignedExtract.bind(this);
    this.changeAssignedGroup = this.changeAssignedGroup.bind(this);
    this.sendAssignments = this.sendAssignments.bind(this);
    this.logout = this.logout.bind(this);
    this.state = {
      moduleCode: "",
      groupings: {},
      docVersion: "1",
      assignments: {},
      isAdmin: false,
      isLoaded: false,
    };
  }

  updateVersion(newVersion) {
    if (
      newVersion != this.state.docVersion &&
      !this.state.assignments[newVersion] &&
      this.state.assignments[this.state.docVersion]
    ) {
      const currentAssignments = this.state.assignments;
      currentAssignments[newVersion] =
        currentAssignments[this.state.docVersion];
      this.setState({ assignments: currentAssignments });
    }
    this.setState({
      docVersion: newVersion,
    });
  }

  sendAssignments() {
    console.log(this.state.assignments);
    updateVersion(this.state.moduleCode, this.state.docVersion);
    updateAssignments(this.state.moduleCode, this.state.assignments);
  }

  changeAssignedGroup(groupNumber, assignedGroup) {
    const currentAssignments = this.state.assignments;
    currentAssignments[this.state.docVersion] =
      currentAssignments[this.state.docVersion] || {};
    if (currentAssignments[this.state.docVersion][groupNumber]) {
      currentAssignments[this.state.docVersion][groupNumber][
        "Assignee"
      ] = assignedGroup;
    } else {
      currentAssignments[this.state.docVersion][groupNumber] = {
        Assignee: assignedGroup,
        Extract: "",
      };
    }
    this.setState({
      assignments: currentAssignments,
    });
  }

  changeAssignedExtract(groupNumber, extract) {
    const currentAssignments = this.state.assignments;
    currentAssignments[this.state.docVersion] =
      currentAssignments[this.state.docVersion] || {};
    if (currentAssignments[this.state.docVersion][groupNumber]) {
      currentAssignments[this.state.docVersion][groupNumber][
        "Extract"
      ] = extract;
    } else {
      currentAssignments[this.state.docVersion][groupNumber] = {
        Assignee: "",
        Extract: extract,
      };
    }

    this.setState({
      assignments: currentAssignments,
    });
  }

  updateModuleInfo(moduleInfo) {
    this.setState(moduleInfo);
    this.setState({ isLoaded: true });
  }

  updateStudentGroup(studentHandle, newGroup) {
    let currentGroupings = this.state.groupings;
    Object.keys(currentGroupings).forEach((key) => {
      const group = new Set(currentGroupings[key]);
      group.delete(studentHandle);
      currentGroupings[key] = [...group];
      if (!currentGroupings[key].length) {
        delete currentGroupings[key];
      }
    });
    console.log(currentGroupings);
    try {
      currentGroupings[newGroup].push(studentHandle);
    } catch (err) {
      currentGroupings[newGroup] = [studentHandle];
    }
    //this.setState({ groupings: currentGroupings });
    updateGroupings(this.state.moduleCode, currentGroupings);
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
              newState.payload.credentials.moduleCode
            ] &&
            newState.payload.credentials[
              newState.payload.credentials.moduleCode
            ] != "unassigned";

          if (!isAdmin && !studentAdded) {
            alert("Ask module admin to assign you");
            logoutUser()(dispatch);
            window.location.href = "/login";
          }

          this.setState({
            isAdmin: newState.payload.credentials.status == "Module Admin",
            moduleCode: newState.payload.credentials.moduleCode,
          });
          getModuleInfo(
            newState.payload.credentials.moduleCode,
            this.updateModuleInfo
          );
        });
      }
    } else {
      logoutUser()(dispatch);
      window.location.href = "/login";
    }
  }

  logout() {
    const { dispatch } = this.context;
    logoutUser()(dispatch);
  }

  render() {
    const userComponent = this.state.isAdmin ? (
      <Admin
        moduleCode={this.state.moduleCode}
        logout={this.logout}
        assignments={[
          this.state.assignments,
          this.changeAssignedGroup,
          this.changeAssignedExtract,
          this.sendAssignments,
        ]}
        version={[this.state.docVersion, this.updateVersion]}
        groupings={this.state.groupings}
        updateStudentGroup={this.updateStudentGroup}
      />
    ) : (
      <Student
        moduleCode={this.state.moduleCode}
        logout={this.logout}
        assignments={[
          this.state.assignments,
          this.changeAssignedGroup,
          this.changeAssignedExtract,
          this.sendAssignments,
        ]}
        version={[this.state.docVersion, this.updateVersion]}
        groupings={this.state.groupings}
        updateStudentGroup={this.updateStudentGroup}
      />
    );

    const showComponent = this.state.isLoaded ? (
      userComponent
    ) : (
      <CircularProgress className="pro" />
    );
    return <div>{showComponent}</div>;
  }
}
home.contextType = Auth;

export default home;
