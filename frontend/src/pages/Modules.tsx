import {
  Card,
  CardActionArea,
  CardContent,
  CardMedia,
  Grid,
  makeStyles,
  Typography,
  Dialog,
  DialogContent,
  DialogTitle,
  Button,
  LinearProgress,
  MenuItem,
  Box
} from "@material-ui/core";
import AddIcon from "@material-ui/icons/Add";
import { FC, useContext, useEffect, useState } from "react";
import { useHistory, useRouteMatch } from "react-router-dom";
import { getStaffModules, getStudentModules, getTAModules, getModules } from "../actions/moduleActions";
import { getUserDetails } from "../actions/userActions";
import { ButtonAppBar, Page } from "../components/NavBar";
import { AuthContext } from "../context/context";
import { AuthType } from "../reducers/reducer";
import { Role } from "./Login";
import { createModule } from "../actions/moduleActions"
import { Formik, Form, Field } from "formik"
import { TextField } from 'formik-material-ui';

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
  add_button: {
    display: "flex",
    height: "100%"
  }
}));

export const ModuleList: FC<{ role: Role }> = (props) => {
  const [modules, setModules] = useState<ModuleInfo[]>([]);
  const [open, setOpen] = useState(false);
  const [pageList, setPageList] = useState<Page[]>([]);
  const { state } = useContext(AuthContext);
  const classes = useStyles();

  var colour = '';

  if (state?.role === Role.STUDENT) {
    colour = 'teal';
  } else if (state?.role === Role.STAFF) {
    colour = 'deepPurple';
  } else if (state?.role === Role.ADMIN) {
    colour = 'orange';
  }

  useEffect(() => {
    if (state?.role === Role.STUDENT) {
      if (props.role === Role.STUDENT) {
        getStudentModules(setModules);
      } else if (props.role === Role.TA) {
        getTAModules(setModules);
      }
      setPageList([
        {
          title: "Modules",
          path: "/student",
        },
        {
          title: "TA Modules",
          path: "/student/ta",
        },
      ])
    } else if (state?.role === Role.STAFF) {
      getStaffModules(setModules);
    } else if (state?.role === Role.ADMIN) {
      getModules({}, setModules);
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
    }
  }, []);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const currentYear = new Date().getFullYear()
  const semesterRanges = [
    {value: `AY${currentYear-1}/${currentYear} 1`, label: `AY${currentYear-1}/${currentYear} 1`},
    {value: `AY${currentYear-1}/${currentYear} 2`, label: `AY${currentYear-1}/${currentYear} 2`},
    {value: `AY${currentYear}/${currentYear + 1} 1`, label: `AY${currentYear}/${currentYear + 1} 1`},
    {value: `AY${currentYear}/${currentYear + 1} 2`, label: `AY${currentYear}/${currentYear + 1} 2`}
  ]

  return (
    <div className={classes.root}>
      <ButtonAppBar pageList={pageList} currentPage={props.role === Role.TA? "TA Modules" : "Modules"} username= {`${state?.user?.Name}`} colour={colour}/>
      <Grid container className="page-background" spacing={3}>
        {
          state?.role === Role.STAFF ? (
          <Grid item className="button-block" xs={12} sm={6} md={4} lg={3} xl={3}>
            <Card className={`${classes.paper} ${classes.add_button}`} style={{height: "210px"}}>
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
        state?.role === Role.STAFF ? (
        <Dialog open={open} onClose={handleClose} disableEnforceFocus>
          <DialogTitle>Create a Module</DialogTitle>
          <DialogContent>
          <Formik
            initialValues={{
              Code: "",
              Name: "",
              Semester: ""
            }}
            validate={(values) => {
              const errors: Partial<ModuleInfo> = {};
              values.Code = values.Code.replace(/(^\s*)|(\s*$)/g, "").toUpperCase()
              values.Name = values.Name.replace(/(^\s*)|(\s*$)/g, "")
              if (!values.Code) {
                errors.Code = 'Required';
              }
              if (!values.Name.replace(/(^\s*)|(\s*$)/g, "")) {
                errors.Name= 'Required';
              }
              if (!values.Semester) {
                errors.Semester = 'Required';
              }
              return errors;
            }}
            onSubmit={(values, {setSubmitting, resetForm}) => {
              createModule(values).then(_ => {
                setSubmitting(false)
                getStaffModules(setModules);
                alert("Successfully created!")
                resetForm()
                setOpen(false)
              }).catch(err => {
                setSubmitting(false)
                if (err && err.response  && err.response.data && err.response.data.message) {
                  alert(err.response.data.message)
                } else {
                  alert("Creation failed.")
                }
              })
            }}
            render={({ submitForm, resetForm, isSubmitting }) => (
              <Form>
                <Box margin={1}>
                  <Field
                    component={TextField}
                    name="Code"
                    type="text"
                    label="Code"
                    helperText="Please Enter Module Code"
                  />
                </Box>
                <Box margin={1}>
                  <Field
                    component={TextField}
                    name="Name"
                    type="text"
                    label="Name"
                    helperText="Please Enter Module Name"
                  />
                </Box>
                <Box margin={2}>
                  <Field
                    component={TextField}
                    name="Semester"
                    type="text"
                    label="Semester"
                    select
                    variant="standard"
                    helperText="Please Select Semester"
                    margin="normal"
                    InputLabelProps={{
                      shrink: true,
                    }}
                  >
                    {semesterRanges.map((option) => (
                      <MenuItem key={option.value} value={option.value}>
                        {option.label}
                      </MenuItem>
                    ))}
                  </Field>
                </Box>
                {isSubmitting && <LinearProgress />}
                <Box 
                  sx={{
                    display: "flex",
                    justifyContent: "space-between"
                  }}
                  margin={1}>
                  <Button
                    variant="contained"
                    color="primary"
                    disabled={isSubmitting}
                    onClick={submitForm}
                  >
                    Submit
                  </Button>
                  <Button
                    variant="contained"
                    color="secondary"
                    disabled={isSubmitting}
                    onClick={() => { resetForm() }}
                  >
                    Reset
                  </Button>
                </Box>
              </Form>
            )}
          >
          </Formik>
          </DialogContent>
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
      <Card className={classes.paper} style={{height: "210px"}}>
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
