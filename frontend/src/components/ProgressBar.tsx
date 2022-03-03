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
  return <Grid item xs={12} sm={12} md={12} lg={12} xl={12}>{showComponent}</Grid>;
}
