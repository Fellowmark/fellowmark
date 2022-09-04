import {
  Button,
  Card,
  CardContent,
  DialogContentText,
  Grid,
  makeStyles,
  Typography,
} from "@material-ui/core";
import React, { FC, useContext, useEffect, useRef, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import { getSubmissionMetadata, uploadSubmission } from "../../../actions/moduleActions";
import { ButtonAppBar, Page } from "../../../components/NavBar";
import { MaxWidthDialog, MaxWidthDialogActions } from "../../../components/PopUpDialog";
import { AuthContext } from "../../../context/context";
import { Role } from "../../Login";
import { Gradebook } from "./Gradebook";
import { PeerReview } from "./PeerReview";
import QuillEditor from "../../../components/QuillEditor";

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

export const getPageList = (match): Page[] => {
  const moduleId = (match.params as { moduleId: number }).moduleId;

  return [
    {
      title: "Assignments",
      path: `/student/module/${moduleId}/assignments`,
    },
  ];
};

export const useValidCheck = (history, authContext, match, setIsValid?) => {
  const { moduleId, assignmentId, questionId } = match.params as {
    moduleId: number;
    assignmentId: number;
    questionId: number;
  };
  useEffect(() => {
    if (authContext?.role !== Role.STUDENT) {
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
  const [peerReview, setPeerReview] = useState(false);
  const [gradeBook, setGradeBook] = useState(false);

  const hiddenFileInput = useRef(null);

  const { moduleId, assignmentId, questionId } = useValidCheck(history, state, match, setIsValid);
  const pageList = getPageList(match);

  const handleUpload = async (file: File) => {
    const data = new FormData();
    data.append('uploadFile', file);
    try {
      await uploadSubmission(data, moduleId, questionId, state.user.ID);
    } catch (e) {
      alert("File too big");
    }
    getSubmissionMetadata(state.user.ID, questionId, setSubmitted);
  }

  useEffect(() => {
    if (isValid) {
      getSubmissionMetadata(state.user.ID, questionId, setSubmitted);
    }
  }, [isValid]);

  return (
    <div>
      <ButtonAppBar
        pageList={pageList}
        currentPage={`${state?.assignment?.Name}`}
        username={`${state?.user?.Name}`}
        colour='teal'
      />
      <MaxWidthDialog
        title="Peer Review"
        open={peerReview}
        setOpen={setPeerReview}
        width={"xl"}
      >
        <DialogContentText>
          Review peers assigned to you
        </DialogContentText>

        <PeerReview moduleId={moduleId} assignmentId={assignmentId} questionId={questionId} />

        <MaxWidthDialogActions handleClose={() => setPeerReview(false)} />
      </MaxWidthDialog>
      <MaxWidthDialog
        title="Gradebook"
        open={gradeBook}
        setOpen={setGradeBook}
        width={"xl"}
      >
        <DialogContentText>
          This is the review you have received for your work
        </DialogContentText>

        <Gradebook moduleId={moduleId} assignmentId={assignmentId} questionId={questionId} />

        <MaxWidthDialogActions handleClose={() => setGradeBook(false)} />
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
      <form noValidate>
        <input
          type="file"
          ref={hiddenFileInput}
          style={{ display: "none" }}
          name="file"
          onChange={(e) => {
            handleUpload(e.target.files[0]);
          }}
        />
        <Button
          style={{ marginTop: "10px", marginBottom: "10px" }}
          variant="contained"
          onClick={() => {
            hiddenFileInput.current.click();
          }}
        >
          {submitted ? "Re-submit" : "Upload"}
        </Button>
      </form>

      <QuillEditor studentId={state.user.ID} questionId={questionId}/>

      <Grid
        container
        direction="row"
        justifyContent="center"
        spacing={3}
      >
        <Grid item>
          <Button
            style={{ marginTop: "10px" }}
            variant="contained"
            onClick={() => {
              setPeerReview(true);
            }}
          >
            Peer Review
          </Button>
        </Grid>
        <Grid item>
          <Button
            style={{ marginTop: "10px" }}
            variant="contained"
            onClick={() => {
              setGradeBook(true);
            }}
          >
            Gradebook
          </Button>
        </Grid>
      </Grid>
    </div>
  );
};
