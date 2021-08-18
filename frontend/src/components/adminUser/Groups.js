import { Component } from "react";
import { Auth } from "../../context/context";

/**
 * Groups consists of the displaying the different groups for a chosen assignment
 * includes options to reset or redo the groupings based on a certain group size
 *  - Table of Pairings (optional: Group Details)
 *  - Reset and Redo Groupings Buttons
 *  - Group Size text (numerical) input to use for pairing
 *
 * Data required:
 *  - Assignment ID
 *
 * Data retrieved:
 *  - Active pairings (optional: translating active pairings into different groups)
 *
 * Data submitted:
 *  - New group size for assignment
 *
 */
class Groups extends Component {
    constructor(props) {
        super(props);
        this.state = {
            groupSize: 0
        };
    }

    render() { }
}

Groups.contextType = Auth;

export default Groups;