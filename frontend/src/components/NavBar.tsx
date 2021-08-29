import { FC, useContext, useState } from "react";

import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import IconButton from "@material-ui/core/IconButton";
import MenuIcon from "@material-ui/icons/Menu";
import Drawer from "@material-ui/core/Drawer";
import { Link, NavLink, useHistory } from "react-router-dom";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemText from "@material-ui/core/ListItemText";

import { AuthContext, TimeoutContext } from "../context/context";
import { logoutUser } from "../actions/userActions";
import { makeStyles } from "@material-ui/core";

export interface Page {
  path: string;
  title: string;
}

export interface MenuProps {
  pageList: Page[];
  currentPage: string;
}

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
  },
  menuButton: {
    marginRight: theme.spacing(5 / 8),
  },

  title: {
    flexGrow: 1,
  },
  drawer: {
    width: theme.spacing(250 / 8),
  },
}));

export const ButtonAppBar: FC<MenuProps> = (props) => {
  const [isDrawn, drawMenu] = useState(false);
  const history = useHistory();
  const { dispatch } = useContext(AuthContext);
  const { cancelTimeout } = useContext(TimeoutContext);
  const classes = useStyles();

  const toggleDrawer = (open: boolean) => () => {
    drawMenu(open);
  };

  return (
    <div className={classes.root}>
      <AppBar>
        <Toolbar>
          {props.pageList.length > 0 && (
            <div>
              <IconButton
                edge="start"
                className={classes.menuButton}
                color="inherit"
                aria-label="menu"
                onClick={toggleDrawer(true)}
              >
                <MenuIcon />
              </IconButton>
              <Drawer
                anchor="left"
                open={isDrawn}
                onClose={toggleDrawer(false)}
              >
                {
                  <div
                    className={classes.drawer}
                    role="presentation"
                    onClick={toggleDrawer(false)}
                  >
                    <List>
                      {props.pageList.map(({ path, title }) => (
                        <ListItem
                          button
                          key={title}
                          onClick={() => { history.push(`${path}`) }}
                        >
                          <ListItemText primary={title} />
                        </ListItem>
                      ))}
                    </List>
                  </div>
                }
              </Drawer>
            </div>
          )}
          <Typography variant="h6" className={classes.title}>
            {props.currentPage}
          </Typography>
          <Button
            color="inherit"
            onClick={() => {
              logoutUser(history, dispatch);
              cancelTimeout();
            }}
            component={Link}
            to="/login"
          >
            Logout
          </Button>
        </Toolbar>
      </AppBar>
    </div>
  );
};
