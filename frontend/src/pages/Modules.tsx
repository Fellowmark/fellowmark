import {
  Card,
  CardActionArea,
  CardContent,
  CardMedia,
  Grid,
  makeStyles,
  Typography,
} from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import { getStaffModules, getStudentModules, getModules } from "../actions/moduleActions";
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

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
  },
  paper: {
    textAlign: "center",
    color: theme.palette.text.secondary,
  },
  add_button: {
    display: "flex",
    height: "100%"
  }
}));

export const ModuleList: FC = () => {
  const [modules, setModules] = useState<ModuleInfo[]>([]);
  const { state } = useContext(AuthContext);
  const classes = useStyles();

  const pageList: Page[] = [];

  useEffect(() => {
    if (state?.role === Role.STUDENT) {
      getStudentModules(setModules);
    } else if (state?.role === Role.STAFF) {
      getStaffModules(setModules);
    } else if (state?.role === Role.ADMIN) {
      getModules({}, setModules);
    }
  }, []);

  const clickAddModule = () => {
    console.log("clickAddModule")
  }

  return (
    <div className={classes.root}>
      <ButtonAppBar pageList={pageList} currentPage="Modules" />
      <Grid container className="page-background" spacing={3}>
        {
          state?.role === Role.ADMIN ? (
          <Grid item className="button-block" xs={12} sm={6} md={4} lg={3} xl={3}>
            <Card className={`${classes.paper} ${classes.add_button}`}>
              <CardActionArea onClick={clickAddModule}>
                <CardContent>
                  <AddIcon />
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>) : null
        }
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
    <Grid item className="button-block" xs={12} sm={6} md={4} lg={3} xl={3}>
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
