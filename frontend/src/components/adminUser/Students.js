import { Component } from "react";

import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Paper from "@material-ui/core/Paper";
import { withStyles } from "@material-ui/core/styles";
import TextField from "@material-ui/core/TextField";

import { Auth } from "../../context/context";

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
    this.setState({ groupings: this.props.groupings });
  }

  render() {
    const currentGroupings = this.state.groupings;
    const studentRows = Object.keys(currentGroupings).map((key) => {
      return currentGroupings[key].map((studentHandle) => {
        return (
          <StyledTableRow key={studentHandle}>
            <StyledTableCell component="th" scope="row">
              {studentHandle}
            </StyledTableCell>
            <StyledTableCell align="right">
              <TextField
                type="number"
                key={studentHandle}
                defaultValue={key}
                onChange={(e) => {
                  console.log(studentHandle);
                  console.log(e.target.value);
                  this.props.updateStudentGroup(studentHandle, e.target.value);
                }}
              />
            </StyledTableCell>
          </StyledTableRow>
        );
      });
    });
    return (
      <div>
        <TableContainer component={Paper}>
          <Table className="Groupings table" aria-label="customized table">
            <TableHead>
              <TableRow>
                <StyledTableCell>Student</StyledTableCell>
                <StyledTableCell align="right">Group</StyledTableCell>
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
