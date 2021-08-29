import {
  Button,
  TableBody,
  TextField,
} from "@material-ui/core";
import { FC, useEffect, useState } from "react";
import {
  createRubrics,
  getRubrics,
} from "../../../actions/moduleActions";
import {
  StyledTableCell,
  StyledTableContainer,
  StyledTableHead,
  StyledTableRow,
} from "../../../components/StyledTable";
import { Question, Rubric } from "../../../models/models";
import { Pagination } from "../../../models/pagination";

export const Rubrics: FC<{ question: Question }> = (props) => {
  const [rubrics, setRubrics] = useState<Pagination<Rubric>>({});
  const [newRubric, setNewRubric] = useState<Rubric>({});

  useEffect(() => {
    getRubrics({ QuestionID: props.question.ID }, setRubrics);
  }, []);

  const createNewRubric = async () => {
    await createRubrics({ ...newRubric, QuestionID: props.question.ID });
    getRubrics({ QuestionID: props.question.ID }, setRubrics);
  }

  return (
    <div>
      <StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Criteria</StyledTableCell>
          <StyledTableCell>Description</StyledTableCell>
          <StyledTableCell>Min</StyledTableCell>
          <StyledTableCell>Max</StyledTableCell>
        </StyledTableHead>
        <TableBody>
          {
            rubrics.rows && rubrics.rows.map((rubric) => {
              return <StyledTableRow
                hover={true}
                key={rubric.ID}
              >
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
              </StyledTableRow>
            })
          }
          <StyledTableRow
            hover={true}
          >
            <StyledTableCell component="th" scope="row">
            </StyledTableCell>
            <StyledTableCell component="th" scope="row">
              <TextField
                type="Criteria"
                placeholder="Criteria"
                name="Criteria"
                multiline
                variant="outlined"
                onChange={(e) =>
                  setNewRubric({ ...newRubric, Criteria: e.target.value })
                }
                required
                autoFocus
              />
            </StyledTableCell>
            <StyledTableCell component="th" scope="row">
              <TextField
                type="Description"
                placeholder="Description"
                name="Description"
                multiline
                variant="outlined"
                onChange={(e) =>
                  setNewRubric({ ...newRubric, Description: e.target.value })
                }
                required
              />
            </StyledTableCell>
            <StyledTableCell component="th" scope="row">
              <TextField
                type="Min"
                placeholder="Min"
                name="Min"
                variant="outlined"
                onChange={(e) =>
                  setNewRubric({ ...newRubric, MinMark: Number(e.target.value) })
                }
                required
              />
            </StyledTableCell>
            <StyledTableCell component="th" scope="row">
              <TextField
                type="Max"
                placeholder="Max"
                name="Max"
                variant="outlined"
                onChange={(e) =>
                  setNewRubric({ ...newRubric, MaxMark: Number(e.target.value) })
                }
                required
              />
            </StyledTableCell>
          </StyledTableRow>
        </TableBody>
      </StyledTableContainer>
      <Button
        color="primary"
        aria-label="menu"
        onClick={() => createNewRubric()}
      >Add Rubric</Button>
    </div>
  );
};
