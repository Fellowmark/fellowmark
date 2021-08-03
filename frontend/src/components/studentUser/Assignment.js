import { Component } from 'react';
import Button from '@material-ui/core/Button';
import { Auth } from '../../context/authContext';

/*
Assignment consists of:
  - Question Dropdown
  - Question Text
  - Submission Card
    - File Upload
    - Text Input
    - Submit Button       

Data required:
  - QuestionIDs

Data retrieved:
  - Question text

Data submitted:
  - Submission file / submission content
*/
class Assignment extends Component {
  constructor(props) {
    super(props);
    this.submitDocument = this.submitDocument.bind(this);
    this.state = {
      questionId: null,
      submitted: false
    };
  }

  submitDocument() {

  }

  render() {
    // const submitButton = <Button variant="contained" size="large" onClick={this.submitDocument}>Submit</Button>;

  }
}

Assignment.contextType = Auth;

export default Assignment;