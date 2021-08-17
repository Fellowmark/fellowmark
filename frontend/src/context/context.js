import React from "react";
import { updateContext } from "../reducers/reducer";

export const Auth = React.createContext();
const initialState = {
    user: {},
    module: localStorage.getItem("module")
};

export const AuthProvider = (props) => {

    const [state, dispatch] = React.useReducer(updateContext, initialState);
    const value = { state, dispatch };

    return <Auth.Provider value={value}>
        {props.children}
    </Auth.Provider>;
};
