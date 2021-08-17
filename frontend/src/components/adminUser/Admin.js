import { Component } from "react";
import NavBar from "../NavBar";
import { Auth } from "../../context/context";
import Students from "./Students";
import GroupWork from "./GroupWork";
import Assessments from "./Assessments";

class Admin extends Component {
  constructor(props) {
    super(props);
    this.updatePage = this.updatePage.bind(this);
    this.state = {
      page: "Students",
    };
  }

  updatePage(page) {
    this.setState({ page: page });
  }

  pageComponent() {
    switch (this.state.page) {
      case "Students":
        return (
          <Students
            groupings={this.props.groupings}
            updateStudentGroup={this.props.updateStudentGroup}
          />
        );
      case "Split Work":
        return (
          <GroupWork
            groupings={this.props.groupings}
            assignments={this.props.assignments}
            version={this.props.version}
          />
        );
      case "Assessment":
        return (
          <Assessments
            moduleCode={this.props.moduleCode}
            assignments={this.props.assignments}
            groupings={this.props.groupings}
            version={this.props.version}
          />
        );
      default:
        return 0;
    }
  }

  render() {
    const pageList = ["Students", "Split Work", "Assessment"];
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
Admin.contextType = Auth;

export default Admin;
