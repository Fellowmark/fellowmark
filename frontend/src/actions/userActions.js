import axios from "axios";
import { addUserToModule } from "./moduleActions";

export const signupUser = (userData, history) => (dispatch) => {
  axios
    .post("/signup", userData)
    .then((res) => {
      console.log(res.data.token);
      setAuthorizationHeader(res.data.token);
      getUserDetails()(dispatch);
      addUserToModule(userData.moduleCode, userData.handle);
      history.push("/");
    })
    .catch((err) => {
      console.error(err.response);
    });
};

export const loginUser = (userData, history) => (dispatch) => {
  axios
    .post("/login", userData)
    .then((res) => {
      console.log(res.data.token);
      setAuthorizationHeader(res.data.token);
      getUserDetails()(dispatch);
      history.push("/");
    })
    .catch((err) => {
      throw new Error(err);
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

export const logoutUser = () => (dispatch) => {
  console.log("Logged out user");
  localStorage.removeItem("FBIdToken");
  delete axios.defaults.headers.common["Authorization"];
  dispatch({
    type: "UNAUTHENTICATED",
    payload: {},
  });
};

const setAuthorizationHeader = (token) => {
  const FBIdToken = `Bearer ${token}`;
  localStorage.setItem("FBIdToken", FBIdToken);
  axios.defaults.headers.common["Authorization"] = FBIdToken;
};
