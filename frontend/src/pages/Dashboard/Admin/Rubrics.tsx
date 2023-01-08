import {
  Button,
  TableBody,
  TextField,
} from "@material-ui/core";
import { FC, useEffect, useState } from "react";
import {
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

  useEffect(() => {
    getRubrics({ QuestionID: props.question.ID }, setRubrics);
  }, []);

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
        </TableBody>
      </StyledTableContainer>
    </div>
  );
};
