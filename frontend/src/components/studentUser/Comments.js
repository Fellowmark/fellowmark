import { Component } from "react";
import Grid from "@material-ui/core/Grid";
import { getOwnDocument } from "../../actions/moduleActions";
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

class Comments extends Component {
  constructor(props) {
    super(props);
    this.updateMarkersList = this.updateMarkersList.bind(this);
    this.updateOwnDoc = this.updateOwnDoc.bind(this);
    this.updateHeading = this.updateHeading.bind(this);
    this.updateVersion = this.updateVersion.bind(this);
    this.updateComments = this.updateComments.bind(this);
    this.parseJSON = this.parseJSON.bind(this);
    this.getTextComponent = this.getTextComponent.bind(this);
    this.state = {
      markers: {},
      doc: {},
      heading: null,
      team: null,
      headingChanged: false,
      version: null,
      versionChanged: false,
      comments: {},
    };
  }

  updateOwnDoc(doc) {
    this.setState({ doc: doc });
  }

  componentDidMount() {
    const { state } = this.context;
    getAllMarkers(
      this.props.moduleCode,
      state.user.credentials[this.props.moduleCode],
      this.props.version[0]
    )(this.updateMarkersList);
  }

  updateMarkersList(markers) {
    this.setState({ markers: markers });
  }

  updateHeading(heading) {
    this.setState({ heading: heading, headingChanged: true }, () => {
      if (this.state.team) {
        this.updateMarker(this.state.team);
      }
    });
  }

  updateMarker(team) {
    this.setState({ team: team });
    console.log(this.state.version);
    getAllGroupVersionComments(
      this.props.moduleCode,
      team,
      this.state.version
    )(this.updateComments);
  }

  updateComments(comments) {
    this.setState({ comments: comments });
  }

  updateVersion(version) {
    this.setState({ 
      version: version,
      versionChanged: true,
      //headingChanged: false
    }, () => {
      console.log("New version: " + this.state.version);
      if (this.state.heading && this.state.team) {
        this.updateMarker(this.state.team);
      }
      getOwnDocument(this.props.moduleCode, version)(this.updateOwnDoc);
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

  render() {
    const versionComponent = (
      <Grid item sm>
        <InputLabel>Version</InputLabel>
        <Select onChange={(e) => { e.preventDefault(); this.updateVersion(e.target.value) }}>
          {Object.keys(this.props.assignments[0]).map((version) => {
            return <MenuItem value={version}>{version}</MenuItem>;
          })}
        </Select>
      </Grid>
    );

    //const splitDoc = splitByHeading(this.state.doc);
    const headingsOptions = this.state.versionChanged ? (
      <Grid item sm>
        <InputLabel>Heading</InputLabel>
        <Select
          onChange={(e) => this.updateHeading(e.target.value)}
        >
          {Object.keys(this.state.markers).map((heading) => {
            return <MenuItem value={heading}>{heading}</MenuItem>;
          })}
        </Select>
      </Grid>
    ) : (
      ""
    );

    const markerSelect = this.state.headingChanged ? (
      <Grid item sm>
        <InputLabel>Group</InputLabel>
        <Select onChange={(e) => this.updateMarker(e.target.value)}>
          {this.state.markers[this.state.heading].map((team) => {
            return <MenuItem value={team}>{team}</MenuItem>;
          })}
        </Select>
      </Grid>
    ) : (
      ""
    );

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
          {headingsOptions}
          {markerSelect}
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

Comments.contextType = Auth;

export default Comments;
