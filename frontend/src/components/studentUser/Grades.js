import { Component } from 'react';
import { Auth } from "../../context/authContext";

/**
 * Grades consists of:
 *  - Table (Row: Assignment-QuestionNo, Col: Grades received) with search feature
 * 
 * Data required:
 *  - Student ID of User
 *  - Module ID
 * 
 * Data retrieved:
 *  - Grade data of those given by group members (check on backend for whether student
 *    has submitted the necessary peer reviews)
 * 
 * Data submitted: None
 */
class Grades extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }

    render() { }
};

Grades.contextType = Auth;

export default Grades;