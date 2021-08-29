import ReactDOM from "react-dom";
import { BrowserRouter } from "react-router-dom";
import App from "./App";
import * as serviceWorker from "./serviceWorker";
import { AuthProvider, TimeoutProvider } from "./context/context";

ReactDOM.render(
  <BrowserRouter>
    <AuthProvider>
      <TimeoutProvider>
        <App />
      </TimeoutProvider>
    </AuthProvider>
  </BrowserRouter>
  , document.getElementById("root"));
serviceWorker.unregister();
