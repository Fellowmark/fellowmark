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
import { FC, useContext, useEffect, useRef, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import {
  getGradesForMarker,
  getGradesForStudent,
  getPairings,
  getRubrics,
  getSubmissionMetadata,
  postGrade,
  uploadSubmission,
} from "../../../actions/moduleActions";
import { ButtonAppBar, Page } from "../../../components/NavBar";
import {
  StyledTableCell,
  StyledTableContainer,
  StyledTableHead,
  StyledTableRow,
} from "../../../components/StyledTable";
import { AuthContext } from "../../../context/context";
import { Grade, Pairing, Rubric } from "../../../models/models";
import { Pagination } from "../../../models/pagination";
import { Role } from "../../Login";

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
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const history = useHistory();
  const [submitted, setSubmitted] = useState(false);
  const [student, setStudent] = useState<number>(null);
  const [pairings, setPairings] = useState<Pagination<Pairing>>({});
  const [rubrics, setRubrics] = useState<Pagination<Rubric>>({});
  const [grades, setGrades] = useState<Map<number, Grade>>(null);

  const hiddenFileInput = useRef(null);

  const { moduleId, questionId } = props;

  useEffect(() => {
    getPairings(moduleId, { StudentID: state.user.ID }, setPairings);
    getRubrics({ QuestionID: questionId }, setRubrics);
  }, []);

  useEffect(() => {
    if (student) {
      getGradesForStudent(moduleId, { PairingID: student }, setGrades);
    }
  }, [student]);

  const handleGrade = () => {
    grades.forEach((value) => {
      return postGrade(moduleId, { ...value, PairingID: student });
    });
  };

  return (
    <div>
      <Select
        name="status"
        onChange={(e) => {
          setStudent(e.target.value as number);
          setGrades(null);
        }}
      >
        {pairings?.rows?.map((pair, key) => {
          return (
            <MenuItem key={key} value={pair.ID}>{`Student ${key + 1
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
    </div>
  );
};
