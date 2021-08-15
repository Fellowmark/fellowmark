import { Auth } from "../../context/authContext";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Paper from "@material-ui/core/Paper";
import { withStyles } from "@material-ui/core/styles";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";

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

function GroupWork(props) {
  const [
    assignments,
    changeAssignedGroup,
    changeAssignedExtract,
    sendAssignments,
  ] = props.assignments;

  const [version, updateVersion] = props.version;
  let currentAssignments;
  if (assignments[version]) {
    currentAssignments = assignments[version];
  } else {
    currentAssignments = {};
  }

  const currentGroups = Object.keys(props.groupings);
  const groupAssignmentRows = currentGroups.map((key) => {
    const otherGroups = currentGroups.filter((e) => e !== key);
    let currentAssignedGroup = currentAssignments[key]
      ? currentAssignments[key]["Assignee"]
      : "";
    currentAssignedGroup = otherGroups.includes(currentAssignedGroup)
      ? currentAssignedGroup
      : "";
    const defaultAssignedExtract = currentAssignments[key]
      ? currentAssignments[key]["Extract"]
      : "";
    return (
      <StyledTableRow key={key}>
        <StyledTableCell component="th" scope="row">
          {key}
        </StyledTableCell>
        <StyledTableCell align="right">
          <Select
            defaultValue={currentAssignedGroup}
            onChange={(e) => changeAssignedGroup(key, e.target.value)}
          >
            {otherGroups.map((group) => {
              return (
                <MenuItem key={group} value={group}>
                  {group}
                </MenuItem>
              );
            })}
          </Select>
        </StyledTableCell>
        <StyledTableCell align="right">
          <TextField
            defaultValue={defaultAssignedExtract}
            onChange={(e) => changeAssignedExtract(key, e.target.value)}
          />
        </StyledTableCell>
      </StyledTableRow>
    );
  });
  return (
    <div>
      <TableContainer component={Paper}>
        <Table className="Groupings table" aria-label="customized table">
          <TableHead>
            <TableRow>
              <StyledTableCell>Group</StyledTableCell>
              <StyledTableCell align="right">Assigned Group</StyledTableCell>
              <StyledTableCell align="right">Extract</StyledTableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {groupAssignmentRows}
            <TableRow>
              <TableCell>
                <TextField
                  label="Version"
                  type="number"
                  onChange={(e) => updateVersion(e.target.value)}
                  defaultValue={version}
                />
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </TableContainer>
      <Button
        align="right"
        variant="contained"
        size="large"
        onClick={sendAssignments}
      >
        Update
      </Button>
    </div>
  );
}
GroupWork.contextType = Auth;

export default GroupWork;
