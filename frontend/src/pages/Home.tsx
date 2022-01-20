import { FC, useContext, useEffect, useState } from "react";

import { AuthContext } from "../context/context";
import { Redirect, useHistory } from "react-router-dom";
import { Role } from "./Login";
import { ProgressBar } from "../components/ProgressBar";
import { ModuleList } from "./Modules";
import { Grid } from "@material-ui/core";

export const Home: React.FC = () => {
  const { state } = useContext(AuthContext);
  const [role, setRole] = useState("");

  useEffect(() => {
    if (state) {
      setRole(state.role);
    }
  }, [state]);

  const showComponent = <Redirect to={`/${role.toLowerCase()}`} />;
  return (
    <Grid container>
      <ProgressBar component={showComponent} isLoaded={true} />
    </Grid>
  );
};

export const RoleHome: FC<{ role: Role }> = (props) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const { state } = useContext(AuthContext);
  const history = useHistory();

  useEffect(() => {
    if (state?.role !== props.role) {
      history.replace("/");
    } else {
      setIsLoaded(true);
    }
  }, []);

  const userComponent = <ModuleList />;

  return (
    <Grid container>
      <ProgressBar component={userComponent} isLoaded={isLoaded} />
    </Grid>
  );
};
