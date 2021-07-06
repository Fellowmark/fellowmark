import { Component } from "react";
import Grid from "@material-ui/core/Grid";
import { getOwnDocument, getGroupAssignedExtract } from "../../actions/moduleActions";
import {
  getAllGroupVersionComments,
  stringToColour,
} from "../../actions/commentsActions";
import { getAllMarkers } from "../../actions/commentsActions";
import { Auth } from "../../context/authContext";
import Select from "@material-ui/core/Select";
import MenuItem from "@material-ui/core/MenuItem";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import Avatar from "@material-ui/core/Avatar";
import InputLabel from "@material-ui/core/InputLabel";

class Assessments extends Component {
  constructor(props) {
    super(props);
    this.updateExtract = this.updateExtract.bind(this);
    this.updateVersion = this.updateVersion.bind(this);
    this.updateComments = this.updateComments.bind(this);
    this.parseJSON = this.parseJSON.bind(this);
    this.getTextComponent = this.getTextComponent.bind(this);
    this.updateGroup = this.updateGroup.bind(this);
    this.state = {
      doc: {},
      version: null,
      versionChanged: false,
      comments: {},
      commentsLoaded: {},
      selectedGroup: null,
    };
  }

  updateExtract(doc) {
    console.log(doc);
    this.setState({ doc: doc });
  }

  componentDidMount() {
  }

  updateComments(comments) {
    this.setState({ comments: comments });
  }

  updateVersion(version) {
    this.setState({ version: version, versionChanged: true }, () => {
      if (this.state.selectedGroup) {
        this.updateGroup(this.state.selectedGroup);
      }
    });
  }

  getTextComponent(text, key, isHeading) {
    if (isHeading) {
      return (
        <Typography key={key} variant="h3" component="h3" gutterBottom>
          {text}
        </Typography>
      );
    } else {
      return (
        <Typography key={key} variant="body1" gutterBottom>
          {text}
        </Typography>
      );
    }
  }

  parseJSON(content) {
    return Object.keys(content).map((key) => {
      const object = content[key];
      if (object["sectionBreak"]) {
        return this.getTextComponent("", key, false);
      } else if (object["paragraph"]) {
        if (object["paragraph"]["paragraphStyle"]["headingId"]) {
          return object["paragraph"]["elements"].map((element) => {
            if (element["textRun"]) {
              return this.getTextComponent(
                element["textRun"]["content"],
                key,
                true
              );
            }
          });
        } else {
          return object["paragraph"]["elements"].map((element) => {
            if (element["textRun"]) {
              return this.getTextComponent(
                element["textRun"]["content"],
                key,
                false
              );
            } else if (element["inlineObjectElement"]) {
              return (
                <iframe
                  src={
                    element["inlineObjectElement"]["textStyle"]["link"]["url"]
                  }
                />
              );
            }
          });
        }
      }
    });
  }

  updateGroup(groupId) {
    this.setState({ selectedGroup: groupId, groupSelected: true }, () => {
      getGroupAssignedExtract(this.props.moduleCode, groupId, this.state.version)(this.updateExtract);
      getAllGroupVersionComments(
        this.props.moduleCode,
        groupId,
        this.state.version
      )(this.updateComments);
    });
  }

  render() {
    const versionComponent = (
      <Grid item sm>
        <InputLabel>Version</InputLabel>
        <Select onChange={(e) => this.updateVersion(e.target.value)}>
          {Object.keys(this.props.assignments[0]).map((version) => {
            return <MenuItem value={version}>{version}</MenuItem>;
          })}
        </Select>
      </Grid>
    );


    const groupsComponent = this.state.versionChanged ? (
      <Grid item sm>
        <InputLabel>Group</InputLabel>
        <Select onChange={(e) => this.updateGroup(e.target.value)}>
          {Object.keys(this.props.groupings).map((groupId) => {
            return <MenuItem value={groupId}>{groupId}</MenuItem>;
          })}
        </Select>
      </Grid>
    ) : "";

    const comments = Object.keys(this.state.comments).map((handle) => {
      const color = stringToColour(handle);
      return this.state.comments[handle].reverse().map((comment, index) => {
        return (
          <Grid item key={handle + index}>
            <Grid container alignItems="center" spacing={2} direction="row">
              <Grid item>
                <Avatar
                  style={{
                    backgroundColor: color,
                  }}
                >
                  {("" + handle[0] + handle[1]).toUpperCase()}
                </Avatar>
              </Grid>
              <Grid item>
                <Typography>
                  <i>{'"' + comment.quote + '"'}</i>
                </Typography>
              </Grid>
              <Grid item>
                <Typography>{comment.comment}</Typography>
              </Grid>
            </Grid>
          </Grid>
        );
      });
    });

    return (
      <Grid container spacing={5} justify="center" direction="column">
        <Grid
          item
          container
          alignItems="center"
          spacing={0}
          justify="center"
          direction="row"
        >
          {versionComponent}
          {groupsComponent}
        </Grid>
        <Grid item>
          <Paper variant="elevation" elevation={2} className="login-background">
            {this.parseJSON(this.state.doc)}
          </Paper>
        </Grid>
        {comments}
      </Grid>
    );
  }
}

Assessments.contextType = Auth;

export default Assessments;
