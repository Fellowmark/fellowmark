import { createContext, Dispatch, FC, useReducer } from "react";
import { Role } from "../pages/Login";
import { AuthType, updateContext } from "../reducers/reducer";

export interface ContextPayload {
  role?: Role;
  user?: any;
  module?: any;
}

export interface ContextState {
  type: AuthType;
  payload: ContextPayload;
}

export interface Context {
  state: ContextPayload;
  dispatch?: Dispatch<ContextState>;
}

const initialState: Context = {
  state: {
  },
};

export const AuthContext = createContext(initialState);

export const AuthProvider: FC = (props) => {
  const [state, dispatch] = useReducer(updateContext, null);
  const value = { state, dispatch };

  return (
    <AuthContext.Provider value={value}>{props.children}</AuthContext.Provider>
  );
};
