import { FC, useContext, useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import { AuthContext } from "../../../context/context";
import { ButtonAppBar, Page } from "../../../components/NavBar";
import { Role } from "../../Login";
import { makeStyles } from "@material-ui/core";

export const StaffManagement: FC = () => {
  const [pageList, setPageList] = useState<Page[]>([]);
  const { state } = useContext(AuthContext);
  const history = useHistory();
  const useStyles = makeStyles(() => ({
    root: {
      flexGrow: 1,
    }
  }));
  const classes = useStyles();

  useEffect(() => {
    if (state?.role !== Role.ADMIN) {
      history.push("/");
    } else {
      setPageList([
        {
          title: "Modules",
          path: "/admin",
        },
        {
          title: "Staff Signup Management",
          path: "/admin/managestaff",
        },
      ])
    }
  }, [])

  return (
    <div className={classes.root}>
      <ButtonAppBar pageList={pageList} currentPage="Staff Signup Management" />
    </div>
  )
}