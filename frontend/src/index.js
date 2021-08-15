import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import * as serviceWorker from "./serviceWorker";
import { AuthProvider } from "./context/authContext";
import { StateProvider } from "./context/modContext";

ReactDOM.render(
    <BrowserRouter>
        <AuthProvider>
            <StateProvider>
                <App />
            </StateProvider>
        </AuthProvider>
    </BrowserRouter>
    , document.getElementById("root"));
serviceWorker.unregister();
