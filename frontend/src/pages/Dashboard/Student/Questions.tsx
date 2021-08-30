import {
  Button,
  DialogContentText,
  FormControl,
  Grid,
  IconButton,
  makeStyles,
  TableBody,
  TextField,
  Typography,
} from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  assignPairings,
  createPairings,
  createQuestion,
  getPairings,
  getQuestions,
} from "../../../actions/moduleActions";
import { ButtonAppBar, Page } from "../../../components/NavBar";
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
import { AuthContext, ContextPayload } from "../../../context/context";
import { Pairing, Question } from "../../../models/models";
import { Pagination } from "../../../models/pagination";
import moment from "moment";
import { getPageList } from "./Dashboard";
import { Role } from "../../Login";
import { AuthType } from "../../../reducers/reducer";

export const useFormStyles = makeStyles((theme) => ({
  form: {
    display: "flex",
    flexDirection: "column",
    margin: "auto",
    width: "fit-content",
    maxHeight: "100%",
  },
  formControl: {
    marginTop: theme.spacing(2),
    minWidth: 120,
  },
  formControlLabel: {
    marginTop: theme.spacing(1),
  },
}));

export const getAssignmentPageList = (match): Page[] => {
  const moduleId = (match.params as { moduleId: number }).moduleId;

  return [
    {
      title: "Class",
      path: `/staff/module/${moduleId}/class`,
    },
    {
      title: "Assignments",
      path: `/staff/module/${moduleId}/assignments`,
    },
  ];
};

export const useAssignmentValidCheck = (
  history,
  authContext: ContextPayload,
  match,
  setIsValid?: (boolean) => void
): { moduleId: number; assignmentId: number } => {
  const moduleId: number = Number(
    (match.params as { moduleId: number }).moduleId
  );
  const assignmentId: number = Number(
    (match.params as { assignmentId: number }).assignmentId
  );

  useEffect(() => {
    if (authContext?.role !== Role.STUDENT) {
      history.push("/");
    }
  }, []);

  useEffect(() => {
    if (
      authContext?.module?.ID !== moduleId ||
      authContext?.assignment?.ID !== assignmentId
    ) {
      history.push("/student");
    } else {
      setIsValid(true);
    }
  }, []);

  return { moduleId: moduleId, assignmentId: assignmentId };
};

export const Questions: FC = () => {
  const match = useRouteMatch();
  const { state, dispatch } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [questions, setQuestions] = useState<Pagination<Question>>({});
  const history = useHistory();

  const { assignmentId } = useAssignmentValidCheck(
    history,
    state,
    match,
    setIsValid
  );

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      getQuestions({ assignmentId: assignmentId }, setQuestions);
    }
  }, [isValid]);

  return (
    <div>
      <ButtonAppBar pageList={pageList} currentPage={state?.assignment?.Name} />
      <div><StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Name</StyledTableCell>
          <StyledTableCell align="right">Deadline</StyledTableCell>
        </StyledTableHead>
        <TableBody>
          <StyledTableRow hover={true} key={state?.assignment?.ID}>
            <StyledTableCell component="th" scope="row">
              {state?.assignment?.ID}
            </StyledTableCell>
            <StyledTableCell component="th" scope="row">
              {state?.assignment?.Name}
            </StyledTableCell>
            <StyledTableCell align="right">
              {state?.assignment?.Deadline
                ? moment.unix(state?.assignment?.Deadline).toLocaleString()
                : "No Deadline"}
            </StyledTableCell>
          </StyledTableRow>
        </TableBody>
      </StyledTableContainer>

        <Typography gutterBottom style={{ marginTop: "10px" }} color="primary">
          Assignment Questions
        </Typography>

        <StyledTableContainer>
          <StyledTableHead>
            <StyledTableCell>ID</StyledTableCell>
            <StyledTableCell>Question Number</StyledTableCell>
            <StyledTableCell align="right">Question Text</StyledTableCell>
          </StyledTableHead>
          <TableBody>
            {questions.rows &&
              questions.rows.map((question) => {
                return (
                  <StyledTableRow
                    onClick={() => {
                      dispatch({
                        type: AuthType.QUESTION,
                        payload: {
                          question: question
                        }
                      });
                      history.push(`${match.url}/question/${question.ID}`);
                    }}
                    hover={true}
                    key={question.ID}
                  >
                    <StyledTableCell component="th" scope="row">
                      {question.ID}
                    </StyledTableCell>
                    <StyledTableCell component="th" scope="row">
                      {question.QuestionNumber}
                    </StyledTableCell>
                    <StyledTableCell
                      aria-multiline={true}
                      align="right"
                      component="th"
                      scope="row"
                    >
                      {question.QuestionText}
                    </StyledTableCell>
                  </StyledTableRow>
                );
              })}
          </TableBody>
        </StyledTableContainer></div>
    </div>
  );
};

export const ViewPairings: FC<{ moduleId: number; assignmentId: number }> = (
  props
) => {
  const [view, setView] = useState(false);
  const [pairings, setPairings] = useState<Pagination<Pairing>>({});

  useEffect(() => {
    getPairings(
      props.moduleId,
      { AssignmentID: props.assignmentId },
      setPairings
    );
  }, []);

  return (
    <div>
      <Button color="primary" aria-label="menu" onClick={() => setView(true)}>
        View Pairings
      </Button>
      <MaxWidthDialog
        title="Pairings"
        setOpen={setView}
        open={view}
        width={"xl"}
      >
        <DialogContentText>
          The following marker student pairs were generated
        </DialogContentText>
        <StyledTableContainer>
          <StyledTableHead>
            <StyledTableCell>ID</StyledTableCell>
            <StyledTableCell>Student</StyledTableCell>
            <StyledTableCell align="right">Marker</StyledTableCell>
          </StyledTableHead>
          <TableBody>
            {pairings.rows &&
              pairings.rows.map((pairing) => {
                return (
                  <StyledTableRow hover={true} key={pairing.ID}>
                    <StyledTableCell component="th" scope="row">
                      {pairing.ID}
                    </StyledTableCell>
                    <StyledTableCell component="th" scope="row">
                      {`${pairing.Student.ID}, ${pairing.Student.Name}, ${pairing.Student.Email}`}
                    </StyledTableCell>
                    <StyledTableCell component="th" scope="row">
                      {`${pairing.Marker.ID}, ${pairing.Marker.Name}, ${pairing.Marker.Email}`}
                    </StyledTableCell>
                  </StyledTableRow>
                );
              })}
          </TableBody>
        </StyledTableContainer>
        <MaxWidthDialogActions
          handleClose={() => setView(false)}
        ></MaxWidthDialogActions>
      </MaxWidthDialog>
    </div>
  );
};