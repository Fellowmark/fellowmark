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
import moment from "moment";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  createPairings,
  createQuestion,
  getAllPairings,
  getAllPairingsId,
  getGradesForStudent,
  getQuestions,
  getTotalGradeForStudent,
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
import { AuthType } from "../../../reducers/reducer";
import { Role } from "../../Login";
import { getPageList } from "./Dashboard";
import { Rubrics } from "./Rubrics";
import axios from "axios";

export var array_name:number[];        //declaration 
array_name = [12, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0];

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
    if (authContext?.role !== Role.STAFF) {
      history.push("/");
    }
  }, []);

  useEffect(() => {
    if (
      authContext?.module?.ID !== moduleId ||
      authContext?.assignment?.ID !== assignmentId
    ) {
      history.push("/staff");
    } else {
      setIsValid(true);
    }
  }, []);

  return { moduleId: moduleId, assignmentId: assignmentId };
};

export const Questions: FC = () => {
  const classes = useFormStyles();

  const match = useRouteMatch();
  const { state, dispatch } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [createNew, setCreateNew] = useState(false);
  const [questions, setQuestions] = useState<Pagination<Question>>({});
  const [selectedQuestion, selectQuestion] = useState<Question>(null);
  const history = useHistory();

  const { assignmentId, moduleId } = useAssignmentValidCheck(
    history,
    state,
    match,
    setIsValid
  );
  const [newQuestion, setNewQuestion] = useState<Question>({
    AssignmentID: assignmentId,
  });

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      getQuestions({ assignmentId: assignmentId }, setQuestions);
    }
  }, [isValid]);

  const addQuestion = async () => {
    await createQuestion(
      newQuestion.QuestionNumber,
      newQuestion.QuestionText,
      assignmentId
    );
    setCreateNew(false);
    setNewQuestion({ AssignmentID: assignmentId });
    getQuestions({ moduleId: moduleId }, setQuestions);
  };

  return (
    <div>
      <ButtonAppBar pageList={pageList} currentPage={state?.assignment?.Name} username={`${state?.user?.Name}`} />
      {isValid && (
        <ViewPairings moduleId={moduleId} assignmentId={assignmentId} />
      )}
      <MaxWidthDialog
        title="Rubric"
        setOpen={(open) => {
          !open && selectQuestion(null);
        }}
        open={Boolean(selectedQuestion)}
        width={"xl"}
      >
        <DialogContentText>
          Rubric provides marking criteria to the markers
        </DialogContentText>
        <Rubrics question={{ ...selectedQuestion }} />
        <MaxWidthDialogActions handleClose={() => selectQuestion(null)} />
      </MaxWidthDialog>

      <MaxWidthDialog
        title="Create Question"
        setOpen={setCreateNew}
        open={createNew}
        width={"xl"}
      >
        <DialogContentText>Please fill in the details</DialogContentText>
        <form className={classes.form} noValidate>
          <FormControl className={classes.formControl}>
            <Grid container direction="column" spacing={2}>
              <Grid item>
                <TextField
                  type="QuestionNumber"
                  placeholder="Question Number"
                  name="QuestionNumber"
                  variant="outlined"
                  onChange={(e) =>
                    setNewQuestion({
                      ...newQuestion,
                      QuestionNumber: Number(e.target.value),
                    })
                  }
                  required
                  autoFocus
                />
              </Grid>
              <Grid item>
                <TextField
                  type="QuestionText"
                  placeholder="Question Text"
                  style={{
                    width: "70vw",
                  }}
                  fullWidth
                  multiline={true}
                  name="QuestionText"
                  variant="outlined"
                  onChange={(e) =>
                    setNewQuestion({
                      ...newQuestion,
                      QuestionText: e.target.value,
                    })
                  }
                  required
                />
              </Grid>
            </Grid>
          </FormControl>
        </form>
        <MaxWidthDialogActions handleClose={() => setCreateNew(false)}>
          <Button onClick={addQuestion} color="primary">
            Add
          </Button>
        </MaxWidthDialogActions>
      </MaxWidthDialog>
      <StyledTableContainer>
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
                        question: question,
                      },
                    });
                    history.push(`${match.url}/question/${question.ID}`);
                    selectQuestion(question);
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

export const ViewPairings: FC<{
  moduleId: number;
  assignmentId: number;
  setPairing?: (pairing: Pairing) => void;
}> = (props) => {
  const [view, setView] = useState(false);
  const [pairings, setPairings] = useState<Pagination<Pairing>>({});

  useEffect(() => {
    getAllPairings(
      { assignmentId: props.assignmentId },
      setPairings
    );
  }, []);

  const generateNewPairings = async () => {
    await createPairings({ id: props.assignmentId });
    getAllPairings(
      { assignmentId: props.assignmentId },
      setPairings
    );
  };

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
        <Button
          color="primary"
          aria-label="menu"
          onClick={() => generateNewPairings()}
        >
          Generate
        </Button>
        <PairingsList assignmentId={props.assignmentId} pairings={pairings} setPairing={props.setPairing} />
        <MaxWidthDialogActions
          handleClose={() => setView(false)}
        ></MaxWidthDialogActions>
      </MaxWidthDialog>
    </div>
  );
};

export const PairingsList: FC<{
  assignmentId : number;
  pairings: Pagination<Pairing>;
  setPairing?: (pairing: Pairing) => void;
}> = (props) => {

  const [pairingsId, setPairingsId] = useState([]);
  //const[grades, setTotalGrade] = useState([]);
  
  useEffect(() => {
    getAllPairingsId(
      { assignmentId: props.assignmentId },
      setPairingsId
    );
  }, []);

  // useEffect(() => {
  //   if (pairingsId) {
  //     //console.log(pairingsId);
  //     getTotalGradeForStudent({ pairingsId: pairingsId }, setTotalGrade);   
  //   }
  // }, [pairingsId]);

  return (
    <>
      <StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Student</StyledTableCell>
          <StyledTableCell>Marker</StyledTableCell>
          <StyledTableCell>Grade</StyledTableCell>
        </StyledTableHead>
        <TableBody>
          {props.pairings?.rows &&
            props.pairings?.rows.map((pairing, index) => {
              return (
                <StyledTableRow
                  hover={true}
                  key={pairing?.ID}
                  onClick={() => {
                    props.setPairing && props.setPairing(pairing);
                    //console.log(pairingsId);
                    //console.log(grades);
                    //getTotalGradeForStudent({ pairingsId: pairingsId }, setTotalGrade);
                  }}
                >
                  <StyledTableCell component="th" scope="row">
                    {pairing?.ID}
                  </StyledTableCell>
                  <StyledTableCell component="th" scope="row">
                    {`${pairing?.Student?.ID}, ${pairing?.Student?.Name}, ${pairing?.Student?.Email}`}
                  </StyledTableCell>
                  <StyledTableCell component="th" scope="row">
                    {`${pairing?.Marker?.ID}, ${pairing?.Marker?.Name}, ${pairing?.Marker?.Email}`}
                  </StyledTableCell>
                  <StyledTableCell component="th" scope="row">
                    {array_name[index]}
                  </StyledTableCell>
                </StyledTableRow>
              );
           })}
        </TableBody>
      </StyledTableContainer>
    </>
  );
};
