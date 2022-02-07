import {
  Button,
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

export const Supervisors: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [supervisions, setSupervisions] = useState<Pagination<Supervision>>({});
  const history = useHistory();

  const moduleId: number = useValidCheck(history, state, match, setIsValid);

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      getSupervisions({ moduleId: moduleId }, setSupervisions);
    }
  }, [isValid]);

  return (
    <div>
      <ButtonAppBar pageList={pageList} currentPage="Supervisors" />
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
            onClick={() => console.log("add supervisor")}
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
                <StyledTableCell align="right">
                  {supervision.Staff.Email}
                </StyledTableCell>
              </StyledTableRow>
            );
          })}
        </TableBody>
      </StyledTableContainer>
    </div>
  );
};
