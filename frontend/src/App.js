import "./App.css";
import { useContext } from "react";
import axios from "axios";
import { Auth } from "./context/authContext";
import Routes from "./routes";

axios.defaults.baseURL =
    "http://localhost:5000";


function App() {
    const { state, dispatch } = useContext(Auth);

    return (
        <div className="App container">
            <Routes />
        </div>
    );
}

export default App;
