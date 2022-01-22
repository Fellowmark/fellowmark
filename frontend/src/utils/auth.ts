import axios from "axios";
import jwtDecode from "jwt-decode";
import { Dispatch } from "react";
import { ContextState } from "../context/context";
import { Role } from "../pages/Login";
import { AuthType } from "../reducers/reducer";

interface Claims {
  data: any;
  exp: any;
}

export function authenticate(dispatch: Dispatch<ContextState>) {
  const token = localStorage.jwt;
  const role = localStorage.role;
  if (token) {
    const claims: Claims = jwtDecode(token);
    if (claims.exp * 1000 > Date.now()) {
      const secondsLeft = claims.exp * 1000 - Date.now();
      setUserContext(claims, role, dispatch);
      axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;
      return secondsLeft;
    }
  }
  return null;
}

export function setUserContext(
  claims: Claims,
  role: Role,
  dispatch: Dispatch<ContextState>
) {
  const context: ContextState = {
    type: AuthType.AUTHENTICATED,
    payload: {
      user: claims.data,
      role: role
    },
  };
  dispatch(context);
}
