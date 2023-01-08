import {
  Button,
  TableBody,
  makeStyles
} from "@material-ui/core";
import PaginationMui from '@material-ui/lab/Pagination';
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  getAssistances,
} from "../../../actions/moduleActions";
import { ButtonAppBar } from "../../../components/NavBar";
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
  const [page, setPage] = useState(1)
  const PAGE_SIZE = 15 //to test
  const [noPagination, setNoPagination] = useState(false)
  const classes = useStyles()
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

  return (
    <div>
      <ButtonAppBar pageList={pageList} currentPage="TAs" username= {`${state?.user?.Name}`} colour='orange'/>
      <StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Name</StyledTableCell>
          <StyledTableCell align="right">Email</StyledTableCell>
        </StyledTableHead>
        <TableBody>
          {assistances.rows?.map((assistance) => {
            return (
              <StyledTableRow key={assistance.Student.ID}>
                <StyledTableCell component="th" scope="row">
                  {assistance.Student.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {assistance.Student.Name}
                </StyledTableCell>
                <StyledTableCell align="right">
                  {assistance.Student.Email}
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
