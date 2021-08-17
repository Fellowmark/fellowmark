/**
 * Staff/Student uses this component to select module after login
 * A Popup with Module Selection
 */

import React, { useState } from 'react';
import { Link } from 'react-router-dom';

import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import Select from '@material-ui/core/Select';

import { Auth } from '../context/context';

export default function ModuleSelector(props) {
    const { state, dispatch } = useContext(Auth);
    const [isOpen, setDialog] = useState(false);
    const [moduleId, setModuleId] = useState('');

    const handleDialogClickOpen = () => {
        setDialog(true);
    };

    const handleDialogClose = () => {
        setDialog(false);
    };

    const setModule = () => {
        localStorage.remove('module');
        localStorage.setItem('module', moduleId);
        dispatch({
            type: "MODULE",
            payload: moduleId
        });
        window.location.reload();
    };

    const handleIdChange = (event) => {
        setModuleId(event.target.value);
    };

    return (
        <div>
            <Button color="inherit" onClick={handleDialogClickOpen} component={Link} to="/login">Change Module</Button>
            <Dialog open={isOpen} onClose={handleDialogClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Select Module</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Please a select a module from below
                    </DialogContentText>
                    // TODO update with list of enrolled/supervised modules based on whether user is a staff or student
                    <FormControl className={classes.formControl}>
                        <Select
                            value={moduleId}
                            onChange={handleIdChange}
                        >
                            <MenuItem value={10}>Ten</MenuItem>
                            <MenuItem value={20}>Twenty</MenuItem>
                            <MenuItem value={30}>Thirty</MenuItem>
                        </Select>
                    </FormControl>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose} color="primary">
                        Cancel
                    </Button>
                    <Button onClick={setModule} color="primary">
                        Confirm
                    </Button>
                </DialogActions>
            </Dialog>
        </div>
    );
}