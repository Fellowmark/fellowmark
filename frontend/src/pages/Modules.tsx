import {
  Card,
  CardActionArea,
  CardContent,
  CardMedia,
  Grid,
  makeStyles,
  Typography,
  TextField,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Button
} from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import { getStaffModules, getStudentModules, getModules } from "../actions/moduleActions";
import { ButtonAppBar, Page } from "../components/NavBar";
import { AuthContext } from "../context/context";
import { AuthType } from "../reducers/reducer";
import { Role } from "./Login";
import { createModule } from "../actions/moduleActions"

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
  const [open, setOpen] = useState(false);
  const [code, setCode] = useState("");
  const [semester, setSemester] = useState("");
  const [name, setName] = useState("");
  const [createSuccess, setCreateSuccess] = useState(false);
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
      setCreateSuccess(false);
    }
  }, [createSuccess]);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleCreate = () => {
    if (code.length == 0) {
      alert("Please enter a module code!")
    } 
    if (semester.length == 0) {
      alert("Please enter a semester!")
    } 
    if (name.length == 0) {
      alert("Please enter a name!")
    } 
    createModule(code, semester, name).then(
      (res) => {
        console.log(res.data)
        alert("Module successfully created!")
        setOpen(false)
        setCreateSuccess(true)
      }
    ).catch(
      (err) => {
        console.log(err)
        alert("Module creation failed!")
      }
    )
  }

  const codeOnBlur = (e) => {
    console.log(e.target.value)
    setCode(e.target.value)
  }

  const semesterOnBlur = (e) => {
    console.log(e.target.value)
    setSemester(e.target.value)
  }

  const nameOnBlur = (e) => {
    console.log(e.target.value)
    setName(e.target.value)
  }

  return (
    <div className={classes.root}>
      <ButtonAppBar pageList={pageList} currentPage="Modules" />
      <Grid container className="page-background" spacing={3}>
        {
          state?.role === Role.ADMIN ? (
          <Grid item className="button-block" xs={12} sm={6} md={4} lg={3} xl={3}>
            <Card className={`${classes.paper} ${classes.add_button}`}>
              <CardActionArea onClick={handleOpen}>
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
      {
        state?.role === Role.ADMIN ? (
        <Dialog open={open} onClose={handleClose}>
          <DialogTitle>Create a Module</DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              margin="dense"
              id="module_code"
              label="Code"
              type="text"
              fullWidth
              variant="standard"
              onBlur={codeOnBlur}
            />
            <TextField
              margin="dense"
              id="module_semester"
              label="Semester"
              type="text"
              fullWidth
              variant="standard"
              onBlur={semesterOnBlur}
            />
            <TextField
              margin="dense"
              id="module_name"
              label="Name"
              type="text"
              fullWidth
              variant="standard"
              onBlur={nameOnBlur}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose}>Cancel</Button>
            <Button onClick={handleCreate}>Create</Button>
          </DialogActions>
        </Dialog>): null
      }
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
