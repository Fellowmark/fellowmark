import { Component } from "react";
import NavBar from "../NavBar";
import { Auth } from "../../context/authContext";
import GDocs from "./GDocs";
import Comments from "./Comments";
import PeerReview from "./PeerReview";

class Student extends Component {
  constructor(props) {
    super(props);
    this.updatePage = this.updatePage.bind(this);
    this.state = {
      page: "Google Docs",
    };
  }

  updatePage(page) {
    this.setState({ page: page });
  }

  pageComponent() {
    switch (this.state.page) {
      case "Google Docs":
        return (
          <GDocs
            moduleCode={this.props.moduleCode}
            updateStudentGroup={this.props.updateStudentGroup}
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
      case "Comments":
        return (
          <Comments
            moduleCode={this.props.moduleCode}
            assignments={this.props.assignments}
            version={this.props.version}
          />
        );
      default:
        return 0;
    }
  }

  render() {
    const pageList = ["Google Docs", "Peer Review", "Comments"];
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
