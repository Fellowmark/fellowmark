import { User } from '../models/models';
import { Role } from '../pages/Login';
import { History } from 'history';
import { Dispatch } from 'react';
import axios from 'axios';

export const signupUser = (role: Role, userData: User, history: History) => (dispatch: Dispatch<any>) => {
  axios
    .post(`${role.toLowerCase()}/auth/signup`, userData)
    .then(() => {
      history.push('/login');
    })
    .catch((err) => {
      console.error(err.response);
    });
};

export const loginUser = (role: Role, userData: User, history: History) => {
  axios
    .get(`${role.toLowerCase()}/auth/login`, {
      params: userData,
    })
    .then((res) => {
      setAuthorizationHeader(res.data.message);
      window.location.href = `/${role}`;
    })
    .catch((err) => {
      alert('Email or password incorrect');
      console.error(err);
    });
};

export const getUserDetails = () => (dispatch: Dispatch<any>) => {
  axios
    .get('/user')
    .then((res) => {
      const context = {
        type: 'AUTHENTICATED',
        payload: res.data,
      };
      dispatch(context);
    })
    .catch((err) => {
      throw new Error(err);
    });
};

export const logoutUser = (history: History, dispatch: Dispatch<any>) => {
  localStorage.removeItem('FBIdToken');
  delete axios.defaults.headers.common['Authorization'];
  dispatch({
    type: 'UNAUTHENTICATED',
    payload: {},
  });
  history.push('/login');
};

export const setAuthorizationHeader = (token: string) => {
  const FBIdToken = `Bearer ${token}`;
  localStorage.setItem('FBIdToken', token);
  axios.defaults.headers.common['Authorization'] = FBIdToken;
};
