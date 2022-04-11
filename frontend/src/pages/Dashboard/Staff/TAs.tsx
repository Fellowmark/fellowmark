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
  getAssistances,
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
import { Assistance } from "../../../models/models";
import { Pagination } from "../../../models/pagination";
import { getPageList, useValidCheck } from "./Dashboard";
import { createAssistance, deleteAssistance } from "../../../actions/moduleActions";

const useStyles = makeStyles((theme) => ({
  error: {
    color: "#f44336;",
    fontSize: "0.75rem"
  },
}));

export const TAs: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [assistances, setAssistances] = useState<Pagination<Assistance>>({});
  const [open, setOpen] = useState(false);
  const [taEmails, setTAEmails] = useState("");
  const [assistanceErrorMessages, setAssistanceErrorMessages] = useState([]);
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [page, setPage] = useState(1)
  const PAGE_SIZE = 2 //to test
  const [noPagination, setNoPagination] = useState(false)
  const classes = useStyles()

  const handleTAEmailsChange = (event) => {
    setTAEmails(event.target.value)
    if (assistanceErrorMessages.length > 0){
      setAssistanceErrorMessages([])
    }
  };

  const history = useHistory();

  const moduleId: number = useValidCheck(history, state, match, setIsValid);

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      if (noPagination) {
        getAssistances({ moduleId: moduleId }, setAssistances);
      } else {
        getAssistances({ moduleId: moduleId, page: page, limit: PAGE_SIZE }, setAssistances);
      }
    }
  }, [isValid, page, noPagination]);

  const handlePageChange = (event, page) => {
    setPage(page)
  }

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const addTAs = () => {
    setIsSubmitting(true)
    let emailArr = taEmails.split(/,|\n|\s/)
    emailArr = emailArr.filter(email => email)
    const emailCount = emailArr.length
    createAssistance(moduleId, emailArr).then(res => {
      const successCount = res.data.success
      const errMessages =  res.data.assistanceErrors
      if (successCount == emailCount) {
        setTAEmails("")
        setAssistanceErrorMessages([])
        const totalRowsAfter = assistances.totalRows + successCount
        const lastPageAfter = Math.ceil(totalRowsAfter / PAGE_SIZE)
        getAssistances({ moduleId: moduleId, page: lastPageAfter, limit: PAGE_SIZE }, setAssistances);
        setPage(lastPageAfter)
        alert("All TAs are added successfully!")
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
        setAssistanceErrorMessages(emailAndErrors)
        setTAEmails(failedEmails.join('\n'))
        if (failedEmails.length < emailCount) {
          const totalRowsAfter = assistances.totalRows + successCount
          let lastPageAfter = Math.ceil(totalRowsAfter / PAGE_SIZE)
          getAssistances({ moduleId: moduleId, page: lastPageAfter, limit: PAGE_SIZE }, setAssistances);
          setPage(lastPageAfter)
        }
        alert("Not all TAs are added successfully. Please check the error messages.")
        setIsSubmitting(false)
      }
    }).catch(err => {
      let message = ""
      if (err && err.response && err.response.data && err.response.data.message) {
        message = "Add TAs failed:" + err.response.data.message
      } else {
        message = "Add TAs failed"
      }
      setAssistanceErrorMessages([message])
      alert(message)
      setIsSubmitting(false)
    })
  }

  const deleteStudent = (assistance) => {
    const confirmed = window.confirm(`Are you sure you want to delete student ${assistance.Student.Name}?`)
    if (!confirmed) {
      return
    }
    deleteAssistance(moduleId, assistance.Student.ID).then(res => {
      let showPage = page
      if (page == assistances.totalPages && assistances.totalRows % PAGE_SIZE == 1) {//last page && only 1 row in last page
        showPage--
      }
      getAssistances({ moduleId: moduleId, page: showPage, limit: PAGE_SIZE }, setAssistances);
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
      <ButtonAppBar pageList={pageList} currentPage="TAs" username= {`${state?.user?.Name}`} colour='deepPurple'/>
      <MaxWidthDialog
        title="Add TAs"
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
              value={taEmails}
              onChange={handleTAEmailsChange}
              variant="outlined"
              error = {assistanceErrorMessages.length > 0}
            />
            {
              assistanceErrorMessages.length > 0 ? (
                <div className={classes.error}>
                  {assistanceErrorMessages.map((message,i) => (
                    <div key={i}>{message}</div>
                  ))}
                </div>
              ) : null
            }
          </DialogContent>
          <DialogActions>
            <Button onClick={addTAs} color="primary" disabled={isSubmitting}>
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
          {assistances.rows?.map((assistance) => {
            return (
              <StyledTableRow key={assistance.Student.ID}>
                <StyledTableCell component="th" scope="row">
                  {assistance.Student.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {assistance.Student.Name}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {assistance.Student.Email}
                </StyledTableCell>
                <StyledTableCell align="right">
                  <Button onClick={() => deleteStudent(assistance)} color="primary">Delete</Button>
                </StyledTableCell>
              </StyledTableRow>
            );
          })}
        </TableBody>
      </StyledTableContainer>
      {
        !noPagination && assistances.totalPages > 1 ?
        <div style={{marginTop: 20, display: 'flex', justifyContent: 'center'}}>
          <PaginationMui count={assistances.totalPages} page={page} onChange={handlePageChange} variant="outlined" color="primary" />
          <Button color="primary" onClick={()=>{setNoPagination(true)}}>Show full list</Button>
        </div> : null 
      }
    </div>
  );
};
