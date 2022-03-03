import { FC, useContext, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { AuthContext } from "../../../context/context";
import { ButtonAppBar, Page } from "../../../components/NavBar";
import { Role } from "../../Login";
import { makeStyles, TableBody } from "@material-ui/core";
import Button from '@material-ui/core/Button';
import Typography from "@material-ui/core/Typography";
import {
  StyledTableCell,
  StyledTableContainer,
  StyledTableHead,
  StyledTableRow,
} from "../../../components/StyledTable";
import { approveStaff, getPendingStaffs, getStaffs } from "../../../actions/userActions";
import { User } from "../../../models/models";
import { Pagination } from "../../../models/pagination";


export const StaffManagement: FC = () => {
  const [pageList, setPageList] = useState<Page[]>([]);
  const [pendingStaffs, setPendingStaffs] = useState<Pagination<User>>({});
  const [staffs, setStaffs] = useState<Pagination<User>>({});
  const { state } = useContext(AuthContext);
  const history = useHistory();
  const useStyles = makeStyles(() => ({
    root: {
      flexGrow: 1,
    }
  }));
  const classes = useStyles();

  useEffect(() => {
    if (state?.role !== Role.ADMIN) {
      history.push("/");
    } else {
      setPageList([
        {
          title: "Modules",
          path: "/admin",
        },
        {
          title: "Staff Signup Management",
          path: "/admin/managestaff",
        },
      ])
      getPendingStaffs({}, setPendingStaffs)
      getStaffs({}, setStaffs)
    }
  }, [])

  const handleApprove = (stf) => {
    approveStaff(stf).then((res) => {
      if (res.success) {
        getPendingStaffs({}, setPendingStaffs)
        getStaffs({}, setStaffs)
      }
    })
  }

  return (
    <div className={classes.root}>
      <ButtonAppBar pageList={pageList} currentPage="Staff Signup Management" />
      <Typography variant="h5" component="div" color="primary" gutterBottom style={{ paddingTop: "10px" }} >
        Pending
      </Typography>
      <StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Name</StyledTableCell>
          <StyledTableCell>Email</StyledTableCell>
          <StyledTableCell align="right">Action</StyledTableCell>
        </StyledTableHead>
        <TableBody>
          {pendingStaffs.rows?.map((stf) => {
            return (
              <StyledTableRow key={stf.ID}>
                <StyledTableCell component="th" scope="row">
                  {stf.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {stf.Name}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {stf.Email}
                </StyledTableCell>
                <StyledTableCell align="right">
                  <Button onClick={() => { handleApprove(stf) }} color="primary">Approve</Button>
                </StyledTableCell>
              </StyledTableRow>
            );
          })}
        </TableBody>
      </StyledTableContainer>

      <Typography variant="h5" component="div" color="primary"  style={{ paddingTop: "10px" }} gutterBottom>
        Approved
      </Typography>
      <StyledTableContainer>
        <StyledTableHead>
          <StyledTableCell>ID</StyledTableCell>
          <StyledTableCell>Name</StyledTableCell>
          <StyledTableCell align="right">Email</StyledTableCell>
        </StyledTableHead>
        <TableBody>
          {staffs.rows?.map((stf) => {
            return (
              <StyledTableRow key={stf.ID}>
                <StyledTableCell component="th" scope="row">
                  {stf.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {stf.Name}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row"  align="right">
                  {stf.Email}
                </StyledTableCell>
              </StyledTableRow>
            );
          })}
        </TableBody>
      </StyledTableContainer>
    </div>
  )
}