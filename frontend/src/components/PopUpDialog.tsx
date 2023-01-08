import React, { FC, useEffect } from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';

const useStyles = makeStyles((theme) => ({
  form: {
    display: 'flex',
    flexDirection: 'column',
    margin: 'auto',
    width: 'fit-content',
    maxHeight: '100%',
  },
  formControl: {
    marginTop: theme.spacing(2),
    minWidth: 120,
  },
  formControlLabel: {
    marginTop: theme.spacing(1),
  },
}));

type Size = 'xs' | 'sm' | 'md' | 'lg' | 'xl' | false;

interface MaxWidthDialogProps {
  setOpen?: (boolean) => void,
  handleSubmit?: () => void,
  title: string,
  open: boolean,
  width: Size,
  children: React.ReactNode,
}


export const MaxWidthDialog: FC<MaxWidthDialogProps> = (props) => {
  const classes = useStyles();
  const fullWidth = true;

  const handleClose = () => {
    props.setOpen(false);
  };

  return (
    <React.Fragment>
      <Dialog
        fullWidth={fullWidth}
        maxWidth={props.width}
        open={props.open}
        onClose={handleClose}
        aria-labelledby="max-width-dialog-title"
      >
        <DialogTitle id="max-width-dialog-title">{props.title}</DialogTitle>
        <DialogContent>
          {props.children}
        </DialogContent>
      </Dialog>
    </React.Fragment>
  );
}

interface MaxWidthDialogActionsProps {
  handleClose: () => void,
  children?: React.ReactNode,
}


export const MaxWidthDialogActions: FC<MaxWidthDialogActionsProps> = (props) => {
  return (
    <DialogActions>
      {props.children}
      <Button onClick={props.handleClose} color="primary">
        Close
      </Button>
    </DialogActions>
  );
}
