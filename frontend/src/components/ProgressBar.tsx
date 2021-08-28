import { Grid, LinearProgress } from "@material-ui/core";
import { FC, ReactElement } from "react";

interface ProgressProps {
  component: ReactElement,
  isLoaded: Boolean
}

export const ProgressBar: FC<ProgressProps> = (props) => {
  const showComponent = props.isLoaded ? (
    props.component
  ) : (
    <LinearProgress className="progressBar" />
  );
  return <Grid item>{showComponent}</Grid>;
}
