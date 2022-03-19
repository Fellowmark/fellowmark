import {
  Button,
  Card,
  CardContent,
  DialogContentText,
  Grid,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  getAllPairings,
  getAllPairingsId,
  getGradesForStudent,
  getSubmissionMetadata,
} from "../../../actions/moduleActions";
import { ButtonAppBar } from "../../../components/NavBar";
import {
  MaxWidthDialog,
  MaxWidthDialogActions,
} from "../../../components/PopUpDialog";
import { AuthContext } from "../../../context/context";
import { Grade, Pairing } from "../../../models/models";
import { Pagination } from "../../../models/pagination";
import { Role } from "../../Login";
import { getPageList } from "./Dashboard";
import { PairingsList } from "./Questions";
import { Review } from "./Review";
import { Rubrics } from "./Rubrics";

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

export const useValidCheck = (history, authContext, match, setIsValid?) => {
  const { moduleId, assignmentId, questionId } = match.params as {
    moduleId: number;
    assignmentId: number;
    questionId: number;
  };
  useEffect(() => {
    if (authContext?.role !== Role.STAFF) {
      history.push("/");
    }
  }, []);

  useEffect(() => {
    if (
      authContext?.module?.ID != moduleId ||
      authContext?.assignment?.ID != assignmentId ||
      authContext?.question?.ID != questionId
    ) {
      history.push("/staff");
    } else {
      setIsValid(true);
    }
  }, []);

  return { moduleId, assignmentId, questionId };
};

export const QuestionBoard: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const history = useHistory();
  const [isValid, setIsValid] = useState(false);
  const [submitted, setSubmitted] = useState(false);
  const [review, setReview] = useState(false);
  const [student, setStudent] = useState<number>(null);
  const [viewRubric, setViewRubric] = useState(false);
  const [grades, setGrades] = useState<Map<number, Grade>>(null);
  const [pairings, setPairings] = useState<Pagination<Pairing>>(null);
  const [selectedPair, selectPair] = useState<Pairing>(null);

  const { moduleId, assignmentId, questionId } = useValidCheck(
    history,
    state,
    match,
    setIsValid
  );
  const pageList = getPageList(match);

  useEffect(() => {
    if (student) {
      getGradesForStudent(moduleId, { PairingID: student }, setGrades);
    }
  }, [student]);
  
  useEffect(() => {
    if (isValid) {
      getSubmissionMetadata(state.user.ID, questionId, setSubmitted);
      getAllPairings({ assignmentId }, setPairings);
    }
  }, [isValid]);

  return (
    <div>
      <ButtonAppBar
        pageList={pageList}
        currentPage={`${state?.assignment?.Name}`}
        username={`${state?.user?.Name}`}
        colour='deepPurple'
      />
      <MaxWidthDialog
        title="Rubric"
        setOpen={setViewRubric}
        open={viewRubric}
        width={"xl"}
      >
        <DialogContentText>
          Rubric provides marking criteria to the markers
        </DialogContentText>
        <Rubrics question={state?.question} />
        <MaxWidthDialogActions handleClose={() => setViewRubric(false)} />
      </MaxWidthDialog>

      <MaxWidthDialog
        title="Review"
        open={Boolean(selectedPair)}
        setOpen={(open) => {
          !open && selectPair(null);
        }}
        width={"xl"}
      >
        <DialogContentText>
          Review submission and feedback of the following pair
        </DialogContentText>

        <Review
          moduleId={moduleId}
          assignmentId={assignmentId}
          questionId={questionId}
          pair={selectedPair}
        />

        <MaxWidthDialogActions handleClose={() => selectPair(null)} />
      </MaxWidthDialog>
      <Card>
        <CardContent>
          <Typography gutterBottom variant="h3">
            {`Question ${state?.question?.QuestionNumber}`}
          </Typography>

          <Typography gutterBottom variant="body1">
            {state?.question?.QuestionText}
          </Typography>
        </CardContent>
      </Card>

      <PairingsList assignmentId={assignmentId} pairings={pairings} setPairing={selectPair} />

      <Grid container direction="row" justifyContent="center" spacing={3}>
        <Grid item>
          <Button
            style={{ marginTop: "10px" }}
            variant="contained"
            onClick={() => {
              setViewRubric(true);
            }}
          >
            View Rubric
          </Button>
        </Grid>
      </Grid>
    </div>
  );
};
