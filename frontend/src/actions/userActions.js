import axios from "axios";

export const signupUser = (role, userData, history) => (dispatch) => {
  axios
    .post(`${role.toLowerCase()}/auth/signup`, userData)
    .then(() => {
      history.push("/login");
    })
    .catch((err) => {
      console.error(err.response);
    });
};

export const loginUser = (role, userData, history) => {
  axios
    .get(`${role.toLowerCase()}/auth/login`, {
      params: userData,
    })
    .then((res) => {
      setAuthorizationHeader(res.data.message)(role);
      window.location.href = `/${role.toLowerCase()}`;
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
  localStorage.removeItem("jwt");
  localStorage.removeItem("role");
  delete axios.defaults.headers.common["Authorization"];
  dispatch({
    type: "UNAUTHENTICATED",
    payload: {},
  });
  history.push("/login");
};

export const setAuthorizationHeader = (token) => (role) => {
  const jwt = `Bearer ${token}`;
  localStorage.setItem("jwt", token);
  localStorage.setItem("role", role);
  axios.defaults.headers.common["Authorization"] = jwt;
};
