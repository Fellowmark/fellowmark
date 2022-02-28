import {
  Button,
  IconButton,
  TableBody,
  TextField,
  DialogContent,
  DialogActions,
  makeStyles
} from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  getEnrollments,
} from "../../../actions/moduleActions";
import { ButtonAppBar } from "../../../components/NavBar";
import {
  MaxWidthDialog,
} from "../../../components/PopUpDialog";
import {
  StyledTableCell,
  StyledTableContainer,
  StyledTableHead,
  StyledTableRow,
} from "../../../components/StyledTable";
import { AuthContext } from "../../../context/context";
import { Enrollment } from "../../../models/models";
import { Pagination } from "../../../models/pagination";
import { getPageList, useValidCheck } from "./Dashboard";
import { createEnrollment } from "../../../actions/moduleActions";

const useStyles = makeStyles((theme) => ({
  error: {
    color: "#f44336;",
    fontSize: "0.75rem"
  },
}));

export const Class: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [students, setStudents] = useState<Pagination<Enrollment>>({});
  const [open, setOpen] = useState(false);
  const [enrollEmails, setEnrollEmails] = useState("");
  const [enrollErrorMessages, setEnrollErrorMessages] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false)
  const classes = useStyles()

  const handleEnrollEmailsChange = (event) => {
    setEnrollEmails(event.target.value)
    if (enrollErrorMessages.length > 0){
      setEnrollErrorMessages([])
    }
  };

  const history = useHistory();

  const moduleId: number = useValidCheck(history, state, match, setIsValid);

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      getEnrollments({ moduleId: moduleId }, setStudents);
    }
  }, [isValid]);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const enrollStudents = () => {
    setIsSubmitting(true)
    let emailArr = enrollEmails.split(/,|\n|\s/)
    emailArr = emailArr.filter(email => email)
    const emailCount = emailArr.length
    createEnrollment(moduleId, emailArr).then(res => {
      const successCount = res.data.success
      const errMessages =  res.data.enrollErrors
      if (successCount == emailCount) {
        setEnrollEmails("")
        setEnrollErrorMessages([])
        getEnrollments({ moduleId: moduleId }, setStudents);
        alert("All students are enrolled successfully!")
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
        setEnrollErrorMessages(emailAndErrors)
        setEnrollEmails(failedEmails.join('\n'))
        if (failedEmails.length < emailCount) {
          getEnrollments({ moduleId: moduleId }, setStudents);
        }
        alert("Not all students are enrolled successfully. Please check the error messages.")
        setIsSubmitting(false)
      }
    }).catch(err => {
      let message = ""
      if (err && err.response && err.response.data && err.response.data.message) {
        message = "Enrollment failed:" + err.response.data.message
      } else {
        message = "Enrollment failed"
      }
      setEnrollErrorMessages([message])
      alert(message)
      setIsSubmitting(false)
    })
  }

  return (
    <div>
      <ButtonAppBar pageList={pageList} currentPage="Class" />
      <MaxWidthDialog
        title="Enroll Stundents"
        setOpen={setOpen}
        open={open}
        width={"xl"}>
          <DialogContent>
            <TextField
              label="Please enter students emails seperated by comma or space or return"
              multiline
              minRows={5}
              maxRows={10}
              fullWidth
              value={enrollEmails}
              onChange={handleEnrollEmailsChange}
              variant="outlined"
              error = {enrollErrorMessages.length > 0}
            />
            {
              enrollErrorMessages.length > 0 ? (
                <div className={classes.error}>
                  {enrollErrorMessages.map((message,i) => (
                    <div key={i}>{message}</div>
                  ))}
                </div>
              ) : null
            }
          </DialogContent>
          <DialogActions>
            <Button onClick={enrollStudents} color="primary" disabled={isSubmitting}>
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
          <StyledTableCell align="right">Email</StyledTableCell>
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
          {students.rows?.map((student) => {
            return (
              <StyledTableRow key={student.Student.ID}>
                <StyledTableCell component="th" scope="row">
                  {student.Student.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {student.Student.Name}
                </StyledTableCell>
                <StyledTableCell align="right">
                  {student.Student.Email}
                </StyledTableCell>
              </StyledTableRow>
            );
          })}
        </TableBody>
      </StyledTableContainer>
    </div>
  );
};
