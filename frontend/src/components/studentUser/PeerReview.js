import { Component } from "react";

import { Auth } from "../../context/context";

import "./PeerReview.css";

/**
 * PeerReview consists of:
 *  - Question Dropdown
 *  - Question Text
 *  - Student Submission
 *  - Rubric (with corresponding dropdowns)
 *  - Submit Button
 * 
 * Data required:
 *  - Student ID of User
 *  - Module ID
 * 
 * Data retrieved:
 *  - Question IDs of questions open for marking
 *    - Question text
 *    - Rubrics
 *  - Student IDs of User's Group
 *    - Submissions by those in Group
 * 
 * Data submitted:
 *  - Marks for each rubric criteria
 *  - Comments
 */
class PeerReview extends Component {
  constructor(props) {
    super(props);
    this.state = {
      submitted: false
    };
  }

  componentDidMount() { }

  componentWillUnmount() { }

  render() { }
}

PeerReview.contextType = Auth;

export default PeerReview;
