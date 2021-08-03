import { Component } from "react";
import NavBar from "../NavBar";
import { Auth } from "../../context/authContext";
import Assignment from "./Assignment";
import Grades from "./Grades";
import PeerReview from "./PeerReview";

class Student extends Component {
  constructor(props) {
    super(props);
    this.updatePage = this.updatePage.bind(this);
    this.state = {
      page: "Assignment",
    };
  }

  updatePage(page) {
    this.setState({ page: page });
  }

  pageComponent() {
    switch (this.state.page) {
      case "Assignment":
        return (
          <Assignment
            module={this.props.moduleId}
          />
        );
      case "Peer Review":
        return (
          <PeerReview
            moduleCode={this.props.moduleCode}
            groupings={this.props.groupings}
            assignments={this.props.assignments}
            version={this.props.version}
          />
        );
      case "Grades":
        return (
          <Grades
          />
        );
      default:
        return 0;
    }
  }

  render() {
    const pageList = ["Assignment", "Peer Review", "Grades"];
    return (
      <div>
        <NavBar
          logout={this.props.logout}
          pageList={pageList}
          updatePage={this.updatePage}
          currentPage={this.state.page}
        />
        {this.pageComponent()}
      </div>
    );
  }
}
Student.contextType = Auth;

export default Student;
