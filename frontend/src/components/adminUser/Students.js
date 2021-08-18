import { Component } from "react";

import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Paper from "@material-ui/core/Paper";
import { withStyles } from "@material-ui/core/styles";

import { Auth } from "../../context/context";
import { getEnrollments } from "../../actions/moduleActions";

const StyledTableCell = withStyles((theme) => ({
  head: {
    backgroundColor: "purple",
    color: theme.palette.common.white,
  },
  body: {
    fontSize: 14,
  },
}))(TableCell);

const StyledTableRow = withStyles((theme) => ({
  root: {
    "&:nth-of-type(odd)": {
      backgroundColor: theme.palette.action.hover,
    },
  },
}))(TableRow);

class Students extends Component {
  constructor(props) {
    super(props);
    this.state = {
      groupings: {},
    };
  }

  componentDidMount() {
    const { state, dispatch } = this.context;
    let enrollments = getEnrollments({ ModuleID: state.Module }).rows;
    this.setState({ enrollments: enrollments });
  }

  render() {
    const enrollments = this.state.enrollments;
    const students = enrollments.map((enrollment) => enrollment.Student);
    const studentRows = students.map((student) => {
      <StyledTableRow key={student}>
        <StyledTableCell component="th" scope="row">
          {student.Name}
        </StyledTableCell>
        <StyledTableCell align="right">
          {student.Email}
        </StyledTableCell>
      </StyledTableRow>;
    });

    return (
      <div>
        <TableContainer component={Paper}>
          <Table className="Groupings table" aria-label="customized table">
            <TableHead>
              <TableRow>
                <StyledTableCell>Student</StyledTableCell>
                <StyledTableCell align="right">Email</StyledTableCell>
              </TableRow>
            </TableHead>
            <TableBody>{studentRows}</TableBody>
          </Table>
        </TableContainer>
      </div>
    );
  }
}

Students.contextType = Auth;

export default Students;
