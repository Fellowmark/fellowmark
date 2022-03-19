import {
  Card,
  CardActionArea,
  CardContent,
  CardMedia,
  Grid,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import { getStaffModules, getStudentModules } from "../actions/moduleActions";
import { getUserDetails } from "../actions/userActions";
import { ButtonAppBar, Page } from "../components/NavBar";
import { AuthContext } from "../context/context";
import { AuthType } from "../reducers/reducer";
import { Role } from "./Login";

export interface ModuleInfo {
  ID?: number;
  Code?: string;
  Semester?: string;
  Name?: string;
}

export interface UserInfo{
  Email?: string;
  Name?: string;
  Password?: string
}

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
  },
  paper: {
    textAlign: "center",
    color: theme.palette.text.secondary,
  },
}));

export const ModuleList: FC = () => {
  const [modules, setModules] = useState<ModuleInfo[]>([]);
  const { state } = useContext(AuthContext);
  const classes = useStyles();
  var colour = '';
  const pageList: Page[] = [];

  if (state?.role === Role.STUDENT) {
    colour = 'pink';
  } else if (state?.role === Role.STAFF) {
    colour = 'deepPurple';
  }

  useEffect(() => {
    if (state?.role === Role.STUDENT) {
      getStudentModules(setModules);
    } else if (state?.role === Role.STAFF) {
      getStaffModules(setModules);
    } else {
    }
  }, []);

  return (
    <div className={classes.root}>
      <ButtonAppBar pageList={pageList} currentPage="Modules" username= {`${state?.user?.Name}`} colour={colour} />
      <Grid container className="page-background" spacing={3}>
        {modules?.map((module) => {
          return <Module key={module.ID} {...module} />;
        })}
      </Grid>
    </div>
  );
};

export const Module: FC<ModuleInfo> = (props) => {
  const classes = useStyles();
  const match = useRouteMatch();
  const history = useHistory();
  const { dispatch } = useContext(AuthContext);

  const clickModule = () => {
    dispatch({
      type: AuthType.MODULE,
      payload: {
        module: props,
      },
    });
    history.push(`${match.url}/module/${props.ID}`);
  };

  return (
    <Grid item className="button-block">
      <Card className={classes.paper}>
        <CardActionArea onClick={clickModule}>
          <CardMedia
            component="img"
            height="80"
            image="https://picsum.photos/400/600"
            title="Contemplative Reptile"
          />
          <CardContent>
            <Typography gutterBottom variant="h4" component="h2">
              {props.Code}
            </Typography>
            <Typography variant="body1" color="textSecondary" component="p">
              {props.Name}
            </Typography>
            <Typography variant="body2" color="textSecondary" component="p">
              {props.Semester}
            </Typography>
          </CardContent>
        </CardActionArea>
      </Card>
    </Grid>
  );
};
