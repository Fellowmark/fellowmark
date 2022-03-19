import {
  Button,
  DialogContentText,
  FormControl,
  Grid,
  IconButton,
  TableBody,
  TextField,
} from "@material-ui/core";
import DateFnsUtils from "@date-io/moment"; // choose your lib
import AddIcon from "@material-ui/icons/Add";
import { DateTimePicker, MuiPickersUtilsProvider } from "@material-ui/pickers";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  createAssignment,
  getAssignments,
} from "../../../actions/moduleActions";
import { ButtonAppBar } from "../../../components/NavBar";
import {
  MaxWidthDialog,
  MaxWidthDialogActions,
} from "../../../components/PopUpDialog";
import {
  StyledTableCell,
  StyledTableContainer,
  StyledTableHead,
  StyledTableRow,
} from "../../../components/StyledTable";
import { AuthContext } from "../../../context/context";
import { Assignment } from "../../../models/models";
import { Pagination } from "../../../models/pagination";
import moment from "moment";
import { getPageList, useFormStyles, useValidCheck } from "./Dashboard";
import { AuthType } from "../../../reducers/reducer";

export const Assignments: FC = () => {
  const classes = useFormStyles();

  const match = useRouteMatch();
  const { state, dispatch } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [createNew, setCreateNew] = useState(false);
  const [assignments, setAssignments] = useState<Pagination<Assignment>>({});
  const history = useHistory();

  const moduleId: number = useValidCheck(history, state, match, setIsValid);
  const [newAssignment, setNewAssignment] = useState<Assignment>({
    ModuleID: moduleId,
  });

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      getAssignments({ moduleId: moduleId }, setAssignments);
    }
  }, [isValid]);

  const addAssignment = async () => {
    await createAssignment(newAssignment);
    setCreateNew(false);
    setNewAssignment({ ModuleID: moduleId });
    getAssignments({ moduleId: moduleId }, setAssignments);
  };

  return (
    <div>
      <ButtonAppBar pageList={pageList} currentPage="Assignments" username={`${state?.user?.Name}`} colour='deepPurple'/>
      <MaxWidthDialog
        title="Create Assignment"
        setOpen={setCreateNew}
        open={createNew}
        width={"sm"}
      >
        <DialogContentText>
          Please fill in the details
        </DialogContentText>
        <form className={classes.form} noValidate>
          <FormControl className={classes.formControl}>
            <Grid container direction="column" spacing={2}>
              <Grid item>
                <TextField
                  type="Name"
                  placeholder="Name"
                  fullWidth
                  name="Name"
                  variant="outlined"
                  onChange={(e) =>
                    setNewAssignment({ ...newAssignment, Name: e.target.value })
                  }
                  required
                  autoFocus
                />
              </Grid>
              <Grid item>
                <TextField
                  type="GroupSize"
                  placeholder="GroupSize"
                  fullWidth
                  name="GroupSize"
                  variant="outlined"
                  onChange={(e) =>
                    setNewAssignment({
                      ...newAssignment,
                      GroupSize: Number(e.target.value),
                    })
                  }
                  required
                />
              </Grid>
              <Grid item>
                <MuiPickersUtilsProvider utils={DateFnsUtils}>
                  <DateTimePicker
                    label="Deadline"
                    inputVariant="outlined"
                    value={
                      newAssignment.Deadline
                        ? moment.unix(newAssignment.Deadline).local().toDate()
                        : moment().toDate()
                    }
                    onChange={(e) =>
                      setNewAssignment({
                        ...newAssignment,
                        Deadline: e.unix(),
                      })
                    }
                  />
                </MuiPickersUtilsProvider>
              </Grid>
            </Grid>
          </FormControl>
        </form>
        <MaxWidthDialogActions handleClose={() => setCreateNew(false)}>
          <Button onClick={addAssignment} color="primary">
            Add
          </Button>
        </MaxWidthDialogActions>
      </MaxWidthDialog>
      <StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Name</StyledTableCell>
          <StyledTableCell>Group Size</StyledTableCell>
          <StyledTableCell align="right">Deadline</StyledTableCell>
        </StyledTableHead>
        <TableBody>
          {assignments.rows?.map((assignment) => {
            return (
              <StyledTableRow onClick={() => {
                dispatch({ type: AuthType.ASSIGNMENT, payload: { assignment: assignment } });
                history.push(`${match.url}/${assignment.ID}`);
              }} hover={true} key={assignment.ID}>
                <StyledTableCell component="th" scope="row">
                  {assignment.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {assignment.Name}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {assignment.GroupSize}
                </StyledTableCell>
                <StyledTableCell align="right">
                  {assignment.Deadline
                    ? moment.unix(assignment.Deadline).toLocaleString()
                    : "No Deadline"}
                </StyledTableCell>
              </StyledTableRow>
            );
          })}
        </TableBody>
      </StyledTableContainer>

      <IconButton
        edge="start"
        color="primary"
        aria-label="menu"
        onClick={() => setCreateNew(true)}
      >
        <AddIcon />
      </IconButton>
    </div>
  );
};
