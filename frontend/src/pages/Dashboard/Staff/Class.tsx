import {
  Button,
  IconButton,
  TableBody,
  TextField,
  DialogContent,
  DialogActions,
  makeStyles
} from "@material-ui/core";
import PaginationMui from '@material-ui/lab/Pagination';
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
import { createEnrollment, deleteEnrollment } from "../../../actions/moduleActions";

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
  const [enrollments, setEnrollments] = useState<Pagination<Enrollment>>({});
  const [open, setOpen] = useState(false);
  const [enrollEmails, setEnrollEmails] = useState("");
  const [enrollErrorMessages, setEnrollErrorMessages] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [page, setPage] = useState(1)
  const classes = useStyles()
  const PAGE_SIZE = 3 //to test
  const [noPagination, setNoPagination] = useState(false)

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
      if (noPagination) {
        getEnrollments({ moduleId: moduleId }, setEnrollments);
      } else {
        getEnrollments({ moduleId: moduleId, page: page, limit: PAGE_SIZE }, setEnrollments);
      }
    }
  }, [isValid, page, noPagination]);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handlePageChange = (event, page) => {
    setPage(page)
  }

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
        const totalRowsAfter = enrollments.totalRows + successCount
        const lastPageAfter = Math.ceil(totalRowsAfter / PAGE_SIZE)
        getEnrollments({ moduleId: moduleId, page: lastPageAfter, limit: PAGE_SIZE }, setEnrollments);
        setPage(lastPageAfter)
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
          const totalRowsAfter = enrollments.totalRows + successCount
          const lastPageAfter = Math.ceil(totalRowsAfter / PAGE_SIZE)
          getEnrollments({ moduleId: moduleId, page: lastPageAfter, limit: PAGE_SIZE }, setEnrollments);
          setPage(lastPageAfter)
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

  const deleteStudent = (enrollment) => {
    const confirmed = window.confirm(`Are you sure you want to delete student ${enrollment.Student.Name}?`)
    if (!confirmed) {
      return
    }
    deleteEnrollment(moduleId, enrollment.Student.ID).then(res => {
      let showPage = page
      if (page == enrollments.totalPages && enrollments.totalRows % PAGE_SIZE == 1) {//last page && only 1 row in last page
        showPage--
      }
      getEnrollments({ moduleId: moduleId, page: showPage, limit: PAGE_SIZE }, setEnrollments);
      setPage(showPage)
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
      <ButtonAppBar pageList={pageList} currentPage="Class" username={`${state?.user?.Name}`} colour='deepPurple'/>
      <MaxWidthDialog
        title="Enroll Students"
        setOpen={setOpen}
        open={open}
        width={"xl"}>
          <DialogContent>
            <TextField
              label="Please enter student emails seperated by comma or space or return"
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
          {enrollments.rows?.map((enrollment) => {
            return (
              <StyledTableRow key={enrollment.Student.ID}>
                <StyledTableCell component="th" scope="row">
                  {enrollment.Student.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {enrollment.Student.Name}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {enrollment.Student.Email}
                </StyledTableCell>
                <StyledTableCell align="right">
                  <Button onClick={() => deleteStudent(enrollment)} color="primary">Delete</Button>
                </StyledTableCell>
              </StyledTableRow>
            );
          })}
        </TableBody>
      </StyledTableContainer>
      {
        !noPagination && enrollments.totalPages > 1 ?
        <div style={{marginTop: 20, display: 'flex', justifyContent: 'center'}}>
          <PaginationMui count={enrollments.totalPages} page={page} onChange={handlePageChange} variant="outlined" color="primary" />
          <Button color="primary" onClick={()=>{setNoPagination(true)}}>Show full list</Button>
        </div> : null 
      }
    </div>
  );
};
