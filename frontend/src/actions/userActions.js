import axios from "axios";
import { authenticate } from "../utils/auth";

// import { addUserToModule } from "./moduleActions";

// export const signupUser = (userData, history) => (dispatch) => {
//   axios
//     .post("/signup", userData)
//     .then((res) => {
//       console.log(res.data.token);
//       setAuthorizationHeader(res.data.token);
//       getUserDetails()(dispatch);
//       addUserToModule(userData.moduleCode, userData.handle);
//       history.push("/");
//     })
//     .catch((err) => {
//       console.error(err.response);
//     });
// };

export const loginUser = (role, userData, history) => (dispatch) => {
  axios
    .get(`${role.toLowerCase()}/auth/login`, {
      params: userData
    })
    .then((res) => {
      console.log(res)
      setAuthorizationHeader(res.data.message);
      console.log(authenticate(dispatch));
      history.push(`/${role.toLowerCase()}`);
    })
    .catch((err) => {
      alert("Email or password incorrect");
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
