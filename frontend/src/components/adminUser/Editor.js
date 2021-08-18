import { Component } from "react";
import { Auth } from "../../context/context";

/**
 * Editor is the assignment editor, where the module supervisor can edit an assignment,
 * consisting of
 *  1. Adding Question (number, text) (ID of created question used in rubrics creation)
 *  1. Adding Rubrics (criteria, description, min_mark, max_mark)
 * with both of the above occurring in the same Action
 * 
 * Data required:
 *  - Assignment ID (user creates new Assignment in Assignments Component)
 * 
 * Data retrieved:
 *  - Questions from Assignment
 *  - Rubrics for retrieved Questions
 * 
 * Data submitted:
 *  - New Questions+Rubrics
 *  - Edited Questions+Rubrics
 */

class Editor extends Component {
    constructor(props) {
        super(props);
    }

    render() { }
}

Editor.contextType = Auth;

export default Editor;