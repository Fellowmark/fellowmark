import axios from "axios";

export const signupUser = (role, userData, history) => (dispatch) => {
  axios
    .post(`${role.toLowerCase()}/auth/signup`, userData)
    .then(() => {
      if (role.toLowerCase() == "staff") {
        alert("To be approved as a staff, please send a email to fellowmarksystem@gmail.com")
      } else {
        history.push("/login");
      }
    })
    .catch((err) => {
      console.error(err.response);
      if (err.response && err.response.data && err.response.data.message) {
        alert(err.response.data.message)
      }
    });
};

export const approveStaff = (stf) => {
  return axios
    .post(`/staff/approve`, stf)
    .then((res) => {
      console.log(res.data)
      alert("Successfully approved!")
      return {success: true, data: res.data}
    })
    .catch((err) => {
      console.error(err);
      alert("Approval failed!")
      return {success: false, err}
    });
};

export const updateStaff = (stf) => {
    return axios
        .post(`/staff`, stf)
        .then((res) => {
            console.log(res.data)
            alert("Successfully updated!")
            return {success: true, data: res.data}
        })
        .catch((err) => {
            console.error(err);
            alert("Update failed!")
            return {success: false, err}
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
      console.error(err);
      if (err && err.response && err.response.data && err.response.data.message) {
        alert(err.response.data.message);
      } else {
        alert("Email or password incorrect");
      }
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

/**
 * Returns PendingStaffs data that matches given data
 *
 * @param {Object} pendingStaffData can consist of ID, Email, Name
 */
 export const getPendingStaffs = (pendingStaffData, setPendingStaffs) => {
  axios
    .get(`/staff/pending`, {
      method: "GET",
      params: pendingStaffData,
    })
    .then((res) => {
      setPendingStaffs(res.data);
    })
    .catch((err) => {
      console.log(err);
    });
};

/**
 * Returns Staffs data that matches given data
 *
 * @param {Object} staffData can consist of ID, Email, Name
 */
 export const getStaffs = (staffData, setStaffs) => {
  axios
    .get(`/staff/approve`, {
      method: "GET",
      params: staffData,
    })
    .then((res) => {
      setStaffs(res.data);
    })
    .catch((err) => {
      console.log(err);
    });
};