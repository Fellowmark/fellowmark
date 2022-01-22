import {
  Button,
  Grid,
  makeStyles,
  MenuItem,
  Select,
  TableBody,
  TextField,
} from "@material-ui/core";
import { FC, useEffect, useState } from "react";
import { IHighlight } from "react-pdf-highlighter";
import {
  downloadSubmission,
  getGradesForMarker,
  getPairingAsMarker,
  getRubrics,
  postGrade,
} from "../../../actions/moduleActions";
import { Annotator } from "../../../components/PdfViewer";
import {
  StyledTableCell,
  StyledTableContainer,
  StyledTableHead,
  StyledTableRow,
} from "../../../components/StyledTable";
import { Grade, Pairing, Rubric } from "../../../models/models";
import { Pagination } from "../../../models/pagination";

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

export const PeerReview: FC<{
  moduleId: number;
  assignmentId: number;
  questionId: number;
}> = (props) => {
  const [studentId, setStudent] = useState<number>(null);
  const [pairings, setPairings] = useState<Pagination<Pairing>>({});
  const [rubrics, setRubrics] = useState<Pagination<Rubric>>({});
  const [grades, setGrades] = useState<Map<number, Grade>>(null);
  const [highlights, setHighlights] = useState<Array<IHighlight>>(
    new Array<IHighlight>()
  );
  const [downloadURL, setDownloadURL] = useState<string>(null);

  const { moduleId, questionId } = props;

  useEffect(() => {
    getPairingAsMarker({ assignmentId: props.assignmentId }, setPairings);
    getRubrics({ QuestionID: questionId }, setRubrics);
  }, []);

  useEffect(() => {
    if (studentId) {
      setGrades(null);
      getGradesForMarker(
        moduleId,
        {
          PairingID: pairings.rows.find((pair) => pair.Student.ID === studentId).ID,
        },
        setGrades
      );
    }
  }, [studentId]);

  const handleGrade = () => {
    grades.forEach((value) => {
      return postGrade(moduleId, {
        ...value,
        PairingID: pairings.rows.find((pair) => pair.Student.ID === studentId).ID,
      });
    });
  };

  const handleDownload = async (studentId: number) => {
    console.log(studentId);
    try {
      setDownloadURL(await downloadSubmission(moduleId, questionId, studentId));
    } catch (e) {
      alert("No submission found");
      console.error(e);
    }
  };

  return (
    <div>
      <Grid
        container
        direction={downloadURL ? "row" : "column"}
        alignItems="center"
        justifyContent="center"
        spacing={1}
        style={{
          marginBottom: "10px",
        }}
      >
        {downloadURL && (
          <Grid item>
            <Annotator
              url={downloadURL}
              setHighlights={setHighlights}
              highlights={highlights}
            />
          </Grid>
        )}
        <Grid item>
          <Grid
            container
            direction="column"
            alignItems="center"
            spacing={1}
            style={{
              marginBottom: "10px",
            }}
          >
            <Grid item>
              <Select
                name="status"
                onChange={(e) => {
                  setStudent(e.target.value as number);
                  handleDownload(e.target.value as number);
                }}
              >
                {pairings?.rows?.map((pair, key) => {
                  return (
                    <MenuItem key={key} value={pair.Student.ID}>{`Student ${
                      key + 1
                    }`}</MenuItem>
                  );
                })}
              </Select>
            </Grid>
            <StyledTableContainer>
              <StyledTableHead>
                <StyledTableCell>ID</StyledTableCell>
                <StyledTableCell>Criteria</StyledTableCell>
                <StyledTableCell>Description</StyledTableCell>
                <StyledTableCell>Min</StyledTableCell>
                <StyledTableCell>Max</StyledTableCell>
                {grades && (
                  <>
                    <StyledTableCell>Grade</StyledTableCell>
                    <StyledTableCell>Comment</StyledTableCell>
                  </>
                )}
              </StyledTableHead>
              <TableBody>
                {rubrics.rows?.map((rubric) => {
                  return (
                    <StyledTableRow hover={true} key={rubric.ID}>
                      <StyledTableCell component="th" scope="row">
                        {rubric.ID}
                      </StyledTableCell>
                      <StyledTableCell component="th" scope="row">
                        {rubric.Criteria}
                      </StyledTableCell>
                      <StyledTableCell component="th" scope="row">
                        {rubric.Description}
                      </StyledTableCell>
                      <StyledTableCell component="th" scope="row">
                        {rubric.MinMark}
                      </StyledTableCell>
                      <StyledTableCell component="th" scope="row">
                        {rubric.MaxMark}
                      </StyledTableCell>
                      {grades && (
                        <>
                          <StyledTableCell component="th" scope="row">
                            <TextField
                              type="Grade"
                              placeholder="Grade"
                              name="Grade"
                              variant="outlined"
                              defaultValue={grades?.get(rubric.ID)?.Grade}
                              onChange={(e) => {
                                const copyGrades = new Map(grades);
                                copyGrades.set(rubric.ID, {
                                  ...copyGrades.get(rubric.ID),
                                  RubricID: rubric.ID,
                                  Grade: Number(e.target.value),
                                });
                                console.log(copyGrades);
                                setGrades(copyGrades);
                              }}
                              required
                            />
                          </StyledTableCell>
                          <StyledTableCell component="th" scope="row">
                            <TextField
                              type="Comment"
                              placeholder="Comment"
                              name="Comment"
                              variant="outlined"
                              defaultValue={grades?.get(rubric.ID)?.Comment}
                              onChange={(e) => {
                                const copyGrades = new Map(grades);
                                copyGrades.set(rubric.ID, {
                                  ...copyGrades.get(rubric.ID),
                                  RubricID: rubric.ID,
                                  Comment: e.target.value,
                                });
                                setGrades(copyGrades);
                              }}
                              required
                            />
                          </StyledTableCell>
                        </>
                      )}
                    </StyledTableRow>
                  );
                })}
              </TableBody>
            </StyledTableContainer>
            <Button
              color="primary"
              variant="contained"
              aria-label="menu"
              disabled={!studentId}
              onClick={() => handleGrade()}
              style={{
                marginTop: "10px",
              }}
            >
              Post
            </Button>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
};
