import { createContext, Dispatch, FC, useReducer, useState } from "react";
import { Assignment, Question } from "../models/models";
import { Role } from "../pages/Login";
import { ModuleInfo } from "../pages/Modules";
import { AuthType, updateContext } from "../reducers/reducer";

export interface ContextPayload {
  role?: Role;
  user?: any;
  module?: ModuleInfo;
  assignment?: Assignment;
  question?: Question;
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

export interface TimeoutState {
  timeout?: NodeJS.Timeout,
  createTimeout?: (timeout: NodeJS.Timeout) => void,
  cancelTimeout?: () => void,
}

export const TimeoutContext = createContext<TimeoutState>({});

export const TimeoutProvider: FC = (props) => {
  const [timeout, createTimeout] = useState<NodeJS.Timeout>(null);

  const cancelTimeout = () => {
    if (timeout) {
      clearTimeout(timeout);
      createTimeout(null);
    }
  }

  const value = { timeout, createTimeout, cancelTimeout }
  return (
    <TimeoutContext.Provider value={value}>{props.children}</TimeoutContext.Provider>
  );
};
