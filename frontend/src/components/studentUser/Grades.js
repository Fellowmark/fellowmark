import { Component } from 'react';
import { Auth } from "../../context/authContext";

/**
 * Grades consists of:
 *  - Table (Row: Assignment-QuestionNo, Col: Grades received) with search feature
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