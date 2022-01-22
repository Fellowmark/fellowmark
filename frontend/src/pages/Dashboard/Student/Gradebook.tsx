import {
  Grid,
  makeStyles,
  MenuItem,
  Select,
  TableBody,
} from "@material-ui/core";
import { FC, useContext, useEffect, useRef, useState } from "react";
import { IHighlight } from "react-pdf-highlighter";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  downloadSubmission,
  getGradesForStudent,
  getPairingAsReviewee,
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
import { AuthContext } from "../../../context/context";
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

export const Gradebook: FC<{
  moduleId: number;
  assignmentId: number;
  questionId: number;
}> = (props) => {
  const { state } = useContext(AuthContext);
  const [student, setStudent] = useState<number>(null);
  const [pairings, setPairings] = useState<Pagination<Pairing>>({});
  const [rubrics, setRubrics] = useState<Pagination<Rubric>>({});
  const [grades, setGrades] = useState<Map<number, Grade>>(null);
  const [highlights, setHighlights] = useState<Array<IHighlight>>(
    new Array<IHighlight>()
  );
  const [downloadURL, setDownloadURL] = useState<string>(null);

  const { moduleId, questionId } = props;

  useEffect(() => {
    getPairingAsReviewee({ assignmentId: props.assignmentId }, setPairings);
    getRubrics({ QuestionID: questionId }, setRubrics);
  }, []);

  useEffect(() => {
    if (student) {
      getGradesForStudent(moduleId, { PairingID: student }, setGrades);
    }
  }, [student]);

  const handleDownload = async () => {
    try {
      setDownloadURL(
        await downloadSubmission(moduleId, questionId, state.user.ID)
      );
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
            <Select
              name="status"
              onChange={(e) => {
                setStudent(e.target.value as number);
                setGrades(null);
                handleDownload();
              }}
            >
              {pairings?.rows?.map((pair, key) => {
                return (
                  <MenuItem key={key} value={pair.ID}>{`Student ${
                    key + 1
                  }`}</MenuItem>
                );
              })}
            </Select>
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
                            {grades?.get(rubric.ID)?.Grade}
                          </StyledTableCell>
                          <StyledTableCell component="th" scope="row">
                            {grades?.get(rubric.ID)?.Comment}
                          </StyledTableCell>
                        </>
                      )}
                    </StyledTableRow>
                  );
                })}
              </TableBody>
            </StyledTableContainer>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
};
