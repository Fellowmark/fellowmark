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
import EditIcon from "@material-ui/icons/Edit";
import { DateTimePicker, MuiPickersUtilsProvider } from "@material-ui/pickers";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  createAssignment,
  editAssignmentCall,
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
import { Select } from "formik-material-ui";

export const Assignments: FC = () => {
  const classes = useFormStyles();

  const match = useRouteMatch();
  const { state, dispatch } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [createNew, setCreateNew] = useState(false);
  const [editNew, setEditNew] = useState(false);
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

  const editAssignment = async () => {
    await editAssignmentCall(newAssignment);
    setEditNew(false);
    setNewAssignment({ ModuleID: moduleId });
    getAssignments({ moduleId: moduleId }, setAssignments);
  }

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
      <MaxWidthDialog
        title="Edit Assignment"
        setOpen={setEditNew}
        open={editNew}
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
                  label="Name"
                  defaultValue={newAssignment.Name}
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
                  label="Group Size"
                  defaultValue={newAssignment.GroupSize}
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
                    defaultValue={newAssignment.Deadline}
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
        <MaxWidthDialogActions handleClose={() => setEditNew(false)}>
          <Button onClick={editAssignment} color="primary">
            Edit
          </Button>
        </MaxWidthDialogActions>
      </MaxWidthDialog>
      <StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Name</StyledTableCell>
          <StyledTableCell>Group Size</StyledTableCell>
          <StyledTableCell>Deadline</StyledTableCell>
          <StyledTableCell align="right">Edit</StyledTableCell>
        </StyledTableHead>
        <TableBody>
          {assignments.rows?.map((assignment) => {
            const { ID, Name, GroupSize, Deadline } = assignment
            return (
              <StyledTableRow onClick={() => {
                console.log("ONCLICK");
                dispatch({ type: AuthType.ASSIGNMENT, payload: { assignment: assignment } });
                history.push(`${match.url}/${assignment.ID}`);
              }} hover={true} key={assignment.ID}>
                <StyledTableCell component="th" scope="row">
                  {ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {Name}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {GroupSize}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {Deadline
                    ? moment.unix(Deadline).toLocaleString()
                    : "No Deadline"}
                </StyledTableCell>
                <StyledTableCell align="right">
                <IconButton
                  edge="start"
                  color="primary"
                  aria-label="menu"
                  onClick={(e: React.MouseEvent<HTMLElement>) => {
                    e.stopPropagation();
                    setEditNew(true);
                    setNewAssignment(assignment => ({
                      ...assignment, 
                      ID,
                      Name,
                      GroupSize, 
                      Deadline 
                    }))
                    //console.log("edit button clicked");
                  }}
                >
                  <EditIcon />
                </IconButton>
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
