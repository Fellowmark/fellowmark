import React, { useState } from 'react';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import Drawer from '@material-ui/core/Drawer';
import { Link } from "react-router-dom";
import List from '@material-ui/core/List';
//import Divider from '@material-ui/core/Divider';
import ListItem from '@material-ui/core/ListItem';
//import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';

import ModuleSelector from './ModuleSelection';

import "./NavBar.css";

export default function ButtonAppBar(props) {
    const [isDrawn, drawMenu] = useState(false);

    const toggleDrawer = (open) => (event) => {
        if (event.type === 'keydown' && (event.key === 'Tab' || event.key === 'Shift')) {
            return;
        }
        drawMenu(open);
    };

    const list = (
        <div
            className="drawer"
            role="presentation"
            onClick={toggleDrawer(false)}
            onKeyDown={toggleDrawer(false)}
        >
            <List>
                {props.pageList.map((text) => (
                    <ListItem button key={text} onClick={() => props.updatePage(text)}>
                        <ListItemText primary={text} />
                    </ListItem>
                ))}
            </List>
        </div>
    );

    return (
        <div className="root">
            <AppBar>
                <Toolbar>
                    <IconButton edge="start" className="menuButton" color="inherit" aria-label="menu" onClick={toggleDrawer(true)}>
                        <MenuIcon />
                    </IconButton>
                    <Drawer anchor="left" open={isDrawn} onClose={toggleDrawer(false)}>
                        {list}
                    </Drawer>
                    <Typography variant="h6" className="title">
                        {props.currentPage}
                    </Typography>
                    <ModuleSelector />
                    <Button color="inherit" onClick={props.logout} component={Link} to="/login">Logout</Button>
                </Toolbar>
            </AppBar>
        </div>
    );
}
