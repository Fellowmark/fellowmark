import {
  makeStyles,
} from "@material-ui/core";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import { Page } from "../../../components/NavBar";
import { AuthContext } from "../../../context/context";
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

export const getPageList = (match): Page[] => {
  const moduleId = (match.params as { moduleId: number }).moduleId;

  return [
    {
      title: "Class",
      path: `/staff/module/${moduleId}/class`,
    },
    {
      title: "Supervisors",
      path: `/staff/module/${moduleId}/supervisors`,
    },
    {
      title: "TAs",
      path: `/staff/module/${moduleId}/tas`,
    },
    {
      title: "Assignments",
      path: `/staff/module/${moduleId}/assignments`,
    },
    {
      title: "Go Back to All Modules",
      path: `/staff`,
    }
  ];
};

export const useValidCheck = (history, authContext, match, setIsValid?) => {
  const moduleId: number = Number(
    (match.params as { moduleId: number }).moduleId
  );
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
};

export const StaffModuleDashboard: FC = () => {
  const match = useRouteMatch();
  const { state } = useContext(AuthContext);
  const history = useHistory();
  const [isValid, setIsValid] = useState(false);

  useValidCheck(history, state, match, setIsValid);

  useEffect(() => {
    if (isValid) {
      history.push(`${match.url}/class`);
    }
  }, [isValid]);

  return <div></div>;
};
