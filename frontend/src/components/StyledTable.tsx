import { FC } from "react";

import Table from "@material-ui/core/Table";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Paper from "@material-ui/core/Paper";
import { withStyles } from "@material-ui/core/styles";

export const StyledTableCell = withStyles((theme) => ({
  head: {
    backgroundColor: "purple",
    color: theme.palette.common.white,
  },
  body: {
    fontSize: 14,
  },
}))(TableCell);

export const StyledTableRow = withStyles((theme) => ({
  root: {
    "&:nth-of-type(odd)": {
      backgroundColor: theme.palette.background.default,
    },
  },
}))(TableRow);

export const StyledTableContainer: FC = (props) => {
  return (
    <TableContainer component={Paper}>
      <Table className="Groupings table" aria-label="customized table">
        {props.children}
      </Table>
    </TableContainer>
  );
}

export const StyledTableHead: FC = (props) => {
  return (
    <TableHead>
      <TableRow>
        {props.children}
      </TableRow>
    </TableHead>
  );
}
