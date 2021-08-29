import axios from "axios";
import { authenticate } from "../utils/auth";

export const signupUser = (role, userData, history) => (dispatch) => {
  axios
    .post(`${role.toLowerCase()}/auth/login`, userData)
    .then(() => {
      history.push("/login");
    })
    .catch((err) => {
      console.error(err.response);
    });
};

export const loginUser = (role, userData, history) => (dispatch) => {
  axios
    .get(`${role.toLowerCase()}/auth/login`, {
      params: userData
    })
    .then((res) => {
      setAuthorizationHeader(res.data.message);
      authenticate(dispatch);
      history.push(`/${role.toLowerCase()}`);
    })
    .catch((err) => {
      alert("Email or password incorrect");
      console.error(err);
    });
};

export const getUserDetails = () => (dispatch) => {
  axios
    .get("/user")
    .then((res) => {
      const context = {
        type: "AUTHENTICATED",
        payload: res.data,
      };
      dispatch(context);
    })
    .catch((err) => {
      throw new Error(err);
    });
};

export const logoutUser = (history, dispatch) => {
  localStorage.removeItem("FBIdToken");
  delete axios.defaults.headers.common["Authorization"];
  dispatch({
    type: "UNAUTHENTICATED",
    payload: {},
  });
  history.push("/login");
};

export const setAuthorizationHeader = (token) => {
  const FBIdToken = `Bearer ${token}`;
  localStorage.setItem("FBIdToken", token);
  axios.defaults.headers.common["Authorization"] = FBIdToken;
};
