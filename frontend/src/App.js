import "./App.css";
import axios from "axios";
import Routes from "./routes";

axios.defaults.baseURL = "http://localhost:5000";

function App() {
    return (
        <div className="App container">
            <Routes />
        </div>
    );
}

export default App;
