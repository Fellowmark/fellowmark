import {
  Button,
  Card,
  CardContent,
  Grid,
  Input,
  makeStyles,
  MenuItem,
  Select,
  TableBody,
  TextField,
  Typography,
} from "@material-ui/core";
import { FC, useEffect, useRef, useState } from "react";
import { IHighlight } from "react-pdf-highlighter";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  downloadSubmission,
  getGradesForMarker,
  getRubrics,
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
import { PairingsList } from "./Questions";

export const Review: FC<{
  moduleId: number;
  assignmentId: number;
  questionId: number;
  pair: Pairing;
}> = (props) => {
  const [rubrics, setRubrics] = useState<Pagination<Rubric>>({});
  const [grades, setGrades] = useState<Map<number, Grade>>(null);
  const [highlights, setHighlights] = useState<Array<IHighlight>>(
    new Array<IHighlight>()
  );
  const [downloadURL, setDownloadURL] = useState<string>(null);
  const ref = useRef(null);

  const { moduleId, questionId } = props;

  useEffect(() => {
    getRubrics({ QuestionID: questionId }, setRubrics);
  }, []);

  useEffect(() => {
    getGradesForMarker(moduleId, { PairingID: props.pair.ID }, setGrades);
  }, []);

  const handleDownload = async () => {
    try {
      setDownloadURL(
        await downloadSubmission(moduleId, questionId, props.pair.Student.ID)
      );
    } catch (e) {
      alert("No submission found");
    }
  };

  return (
    <div>
      <PairingsList pairings={{ rows: [props.pair] }} />
      <a style={{ display: "none" }} href="empty" ref={ref}>
        ref
      </a>
      <Grid
        container
        direction={downloadURL ? "row" : "column"}
        alignItems="center"
        justifyContent="center"
        spacing={1}
        style={{
          marginBottom: "20px",
        }}
      >
        <Grid item>
          {!downloadURL ? (
            <Button
              color="primary"
              variant="contained"
              aria-label="menu"
              onClick={() => handleDownload()}
            >
              Open
            </Button>
          ) : (
            <Annotator
              url={downloadURL}
              setHighlights={setHighlights}
              highlights={highlights}
            />
          )}
        </Grid>
        <Grid item>
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
    </div>
  );
};
