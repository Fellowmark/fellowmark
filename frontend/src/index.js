import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import * as serviceWorker from "./serviceWorker";
import { AuthProvider } from "./context/authContext"

ReactDOM.render(
    <BrowserRouter>
        <AuthProvider>
            <App /> 
        </AuthProvider>
    </BrowserRouter>
    , document.getElementById("root"));
serviceWorker.unregister();
