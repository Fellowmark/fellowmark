import {
  Button,
  FormControl,
  Grid,
  IconButton,
  TableBody,
  TextField,
} from "@material-ui/core";
import PaginationMui from '@material-ui/lab/Pagination';
import DateFnsUtils from "@date-io/moment"; // choose your lib
import AddIcon from "@material-ui/icons/Add";
import { DateTimePicker, MuiPickersUtilsProvider } from "@material-ui/pickers";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  createAssignment,
  getAssignments,
  getEnrollments,
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
import { Assignment, Enrollment } from "../../../models/models";
import { Pagination } from "../../../models/pagination";
import moment from "moment";
import { getPageList, useFormStyles, useValidCheck } from "./Dashboard";

export const Class: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [students, setStudents] = useState<Pagination<Enrollment>>({});
  const history = useHistory();
  const [page, setPage] = useState(1)
  const PAGE_SIZE = 5 //to test
  const [noPagination, setNoPagination] = useState(false)

  const moduleId: number = useValidCheck(history, state, match, setIsValid);

  const pageList = getPageList(match);

  const handlePageChange = (event, page) => {
    setPage(page)
  }

  useEffect(() => {
    if (isValid) {
      if (noPagination) {
        getEnrollments({ moduleId: moduleId }, setStudents);
      } else {
        getEnrollments({ moduleId: moduleId, page: page, limit: PAGE_SIZE }, setStudents);
      }
      }
  }, [isValid, page, noPagination]);

  return (
    <div>
      <ButtonAppBar pageList={pageList} currentPage="Class" username= {`${state?.user?.Name}`} colour='deepPurple'/>
      <StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Name</StyledTableCell>
          <StyledTableCell align="right">Email</StyledTableCell>
        </StyledTableHead>
        <TableBody>
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
      {
        !noPagination && students.totalPages > 1 ?
        <div style={{marginTop: 20, display: 'flex', justifyContent: 'center'}}>
          <PaginationMui count={students.totalPages} page={page} onChange={handlePageChange} variant="outlined" color="primary" />
          <Button color="primary" onClick={()=>{setNoPagination(true)}}>Show full list</Button>
        </div> : null 
      }
    </div>
  );
};
