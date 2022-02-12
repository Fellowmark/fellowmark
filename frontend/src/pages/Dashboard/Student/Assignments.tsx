import {
  TableBody,
} from "@material-ui/core";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  getAssignments,
} from "../../../actions/moduleActions";
import { ButtonAppBar } from "../../../components/NavBar";
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
import { getPageList, useValidCheck } from "./Dashboard";
import { AuthType } from "../../../reducers/reducer";

export const Assignments: FC = () => {
  const match = useRouteMatch();
  const { state, dispatch } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [assignments, setAssignments] = useState<Pagination<Assignment>>({});
  const history = useHistory();

  const moduleId: number = useValidCheck(history, state, match, setIsValid);

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      getAssignments({ moduleId: moduleId }, setAssignments);
    }
  }, [isValid]);

  return (
    <div>
      <ButtonAppBar pageList={pageList} currentPage="Assignments" username={`${state?.user?.Name}`}/>
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
    </div>
  );
};
