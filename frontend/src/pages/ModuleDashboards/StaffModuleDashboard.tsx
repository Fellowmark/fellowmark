import { TableBody } from "@material-ui/core";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import { getEnrollments } from "../../actions/moduleActions";
import { ButtonAppBar, Page } from "../../components/NavBar";
import { StyledTableCell, StyledTableContainer, StyledTableHead, StyledTableRow } from "../../components/StyledTable";
import { AuthContext } from "../../context/context";
import { Enrollment } from "../../models/models";
import { Pagination } from "../../models/pagination";
import { Role } from "../Login";

const getPageList = (match): Page[] => {
  const moduleId = (match.params as { moduleId: number }).moduleId;

  return [{
    title: "Class",
    path: `/staff/module/${moduleId}/class`
  }, {
    title: "Assignments",
    path: `/staff/module/${moduleId}/assignments`
  }]
};

const useValidCheck = (history, authContext, match, setIsValid?) => {
  const moduleId: number = Number((match.params as { moduleId: number }).moduleId);
  useEffect(() => {
    if (authContext?.role !== Role.STAFF) {
      history.push("/");
    }
  }, []);

  useEffect(() => {
    if (authContext?.module?.ID !== moduleId) {
      history.push("/staff");
    } else {
      setIsValid(true);
    }
  }, []);

  return moduleId;
}

export const StaffModuleDashboard: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const history = useHistory();
  const [isValid, setIsValid] = useState(false);

  const pageList = getPageList(match);

  useValidCheck(history, state, match, setIsValid);

  useEffect(() => {
    if (isValid) {
      history.push(`${match.url}/class`);
    }
  }, [isValid]);


  return <div>
  </div>;
}

export const Class: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [students, setStudents] = useState<Pagination<Enrollment>>({});
  const history = useHistory();

  const moduleId: number = useValidCheck(history, state, match, setIsValid);

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      getEnrollments({ moduleId: moduleId }, setStudents);
    }
  }, [isValid]);

  return <div>
    <ButtonAppBar pageList={pageList} currentPage="Class" />
    <StyledTableContainer>
      <StyledTableHead>
        <StyledTableCell>ID</StyledTableCell>
        <StyledTableCell>Name</StyledTableCell>
        <StyledTableCell align="right">Email</StyledTableCell>
      </StyledTableHead>
      <TableBody>
        {
          students.rows?.map((student) => {
            return (
              <StyledTableRow key={student.Student.ID}>
                <StyledTableCell component="th" scope="row">
                  {student.Student.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {student.Student.Name}
                </StyledTableCell>
                <StyledTableCell align="right">
                  {student.Student.Email}
                </StyledTableCell>
              </StyledTableRow>
            );
          })
        }
      </TableBody>
    </StyledTableContainer>
  </div>;
}

export const Assignments: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const [isValid, setIsValid] = useState(false);
  const [students, setStudents] = useState<Pagination<Enrollment>>({});
  const history = useHistory();

  const moduleId: number = useValidCheck(history, state, match, setIsValid);

  const pageList = getPageList(match);

  useEffect(() => {
    if (isValid) {
      getEnrollments({ moduleId: moduleId }, setStudents);
    }
  }, [isValid]);

  return <div>
    <ButtonAppBar pageList={pageList} currentPage="Assignments" />
    <StyledTableContainer>
      <StyledTableHead>
        <StyledTableCell>ID</StyledTableCell>
        <StyledTableCell>Name</StyledTableCell>
        <StyledTableCell align="right">Email</StyledTableCell>
      </StyledTableHead>
      <TableBody>
        {
          students.rows?.map((student) => {
            return (
              <StyledTableRow key={student.Student.ID}>
                <StyledTableCell component="th" scope="row">
                  {student.Student.ID}
                </StyledTableCell>
                <StyledTableCell component="th" scope="row">
                  {student.Student.Name}
                </StyledTableCell>
                <StyledTableCell align="right">
                  {student.Student.Email}
                </StyledTableCell>
              </StyledTableRow>
            );
          })
        }
      </TableBody>
    </StyledTableContainer>
  </div>;
}
