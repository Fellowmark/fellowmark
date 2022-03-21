import {
  Button,
  FormControl,
  Grid,
  IconButton,
  TableBody,
  TextField,
  DialogContent,
  DialogActions,
  makeStyles
} from "@material-ui/core";
import DateFnsUtils from "@date-io/moment"; // choose your lib
import AddIcon from "@material-ui/icons/Add";
import { DateTimePicker, MuiPickersUtilsProvider } from "@material-ui/pickers";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  createAssignment,
  getAssignments,
  getSupervisions,
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
import { Assignment, Supervision } from "../../../models/models";
import { Pagination } from "../../../models/pagination";
import moment from "moment";
import { getPageList, useFormStyles, useValidCheck } from "./Dashboard";
import { createSupervision, deleteSupervision } from "../../../actions/moduleActions";

const useStyles = makeStyles((theme) => ({
  error: {
    color: "#f44336;",
    fontSize: "0.75rem"
  },
}));

export const Supervisors: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [supervisions, setSupervisions] = useState<Pagination<Supervision>>({});
  const [open, setOpen] = useState(false);
  const [superviseEmails, setSuperviseEmails] = useState("");
  const [superviseErrorMessages, setSuperviseErrorMessages] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false)
  const classes = useStyles()
  const history = useHistory();

  const moduleId: number = useValidCheck(history, state, match, setIsValid);

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      getSupervisions({ moduleId: moduleId }, setSupervisions);
    }
  }, [isValid]);

  const handleSuperviseEmailsChange = (event) => {
    setSuperviseEmails(event.target.value)
    if (superviseErrorMessages.length > 0){
      setSuperviseErrorMessages([])
    }
  };

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const addStaffs = () => {
    setIsSubmitting(true)
    let emailArr = superviseEmails.split(/,|\n|\s/)
    emailArr = emailArr.filter(email => email)
    const emailCount = emailArr.length
    createSupervision(moduleId, emailArr).then(res => {
      const successCount = res.data.success
      const errMessages =  res.data.superviseErrors
      if (successCount == emailCount) {
        setSuperviseEmails("")
        setSuperviseErrorMessages([])
        getSupervisions({ moduleId: moduleId }, setSupervisions);
        alert("All supervisors are added successfully!")
        setIsSubmitting(false)
      } else {
        const failedEmails = []
        const emailAndErrors = []
        for (let i = 0; i < emailCount; i++) {
          if (errMessages[i].length > 0) {
            failedEmails.push(emailArr[i])
            emailAndErrors.push(emailArr[i] + " : " + errMessages[i])
          }
        }
        setSuperviseErrorMessages(emailAndErrors)
        setSuperviseEmails(failedEmails.join('\n'))
        if (failedEmails.length < emailCount) {
          getSupervisions({ moduleId: moduleId }, setSupervisions);
        }
        alert("Not all supervisors are added successfully. Please check the error messages.")
        setIsSubmitting(false)
      }
    }).catch(err => {
      let message = ""
      if (err && err.response && err.response.data && err.response.data.message) {
        message = "Add supervisors failed:" + err.response.data.message
      } else {
        message = "Add supervisors failed"
      }
      setSuperviseErrorMessages([message])
      alert(message)
      setIsSubmitting(false)
    })
  }

  const deleteStaff = (supervision) => {
    const confirmed = window.confirm(`Are you sure you want to delete supervisor ${supervision.Staff.Name}?`)
    if (!confirmed) {
      return
    }
    deleteSupervision(moduleId, supervision.Staff.ID).then(res => {
      getSupervisions({ moduleId: moduleId }, setSupervisions);
      alert("Successfully deleted!")
    }).catch(err => {
      if (err && err.response && err.response.data && err.response.data.message) {
        alert("Deletion failed:" + err.response.data.message)
      } else {
        alert("Deletion failed")
      }
    })
  }

  return (
    <div>
      <ButtonAppBar pageList={pageList} currentPage="Supervisors" />
      <MaxWidthDialog
        title="Add Supervisors"
        setOpen={setOpen}
        open={open}
        width={"xl"}>
          <DialogContent>
            <TextField
              label="Please enter staff emails seperated by comma or space or return"
              multiline
              minRows={5}
              maxRows={10}
              fullWidth
              value={superviseEmails}
              onChange={handleSuperviseEmailsChange}
              variant="outlined"
              error = {superviseErrorMessages.length > 0}
            />
            {
              superviseErrorMessages.length > 0 ? (
                <div className={classes.error}>
                  {superviseErrorMessages.map((message,i) => (
                    <div key={i}>{message}</div>
                  ))}
                </div>
              ) : null
            }
          </DialogContent>
          <DialogActions>
            <Button onClick={addStaffs} color="primary" disabled={isSubmitting}>
              Add
            </Button>
            <Button onClick={handleClose} color="primary" disabled={isSubmitting}>
              Close
            </Button>
          </DialogActions>
      </MaxWidthDialog>
      <StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Name</StyledTableCell>
          <StyledTableCell>Email</StyledTableCell>
          <StyledTableCell align="right">Delete</StyledTableCell>
        </StyledTableHead>
        <TableBody>
          <IconButton
            edge="end"
            color="primary"
            aria-label="add"
            onClick={handleOpen}
          >
            <AddIcon />
          </IconButton>
          {supervisions.rows?.map((supervision) => {
            return (
              <StyledTableRow key={supervision.Staff.ID}>
                <StyledTableCell component="th" scope="row">
                  {supervision.Staff.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {supervision.Staff.Name}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {supervision.Staff.Email}
                </StyledTableCell>
                <StyledTableCell align="right">
                  <Button onClick={() => deleteStaff(supervision)} color="primary">Delete</Button>
                </StyledTableCell>
              </StyledTableRow>
            );
          })}
        </TableBody>
      </StyledTableContainer>
    </div>
  );
};
