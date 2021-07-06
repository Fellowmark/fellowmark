import { Component } from "react";
import { getAssignedExtract } from "../../actions/moduleActions";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import {
  commentOnExtract,
  getAllGroupVersionComments,
  removeCommentByIndex,
  stringToColour
} from "../../actions/commentsActions";
import { Auth } from "../../context/authContext";
import Avatar from "@material-ui/core/Avatar";
import "./PeerReview.css";
import IconButton from "@material-ui/core/IconButton";
import DeleteIcon from "@material-ui/icons/Delete";

class PeerReview extends Component {
  constructor(props) {
    super(props);
    this.updateExtractJson = this.updateExtractJson.bind(this);
    this.parseJSON = this.parseJSON.bind(this);
    this.getTextComponent = this.getTextComponent.bind(this);
    this.commentHandler = this.commentHandler.bind(this);
    this.updateCommentHistory = this.updateCommentHistory.bind(this);
    this.pollAllComments = this.pollAllComments.bind(this);
    this.poll = 0;
    this.state = {
      extract: {},
      commentMetaData: {},
      showCommentBox: false,
      commentHistory: [],
    };
  }

  updateCommentHistory(comments) {
    this.setState({ commentHistory: comments });
  }

  pollAllComments() {
    const { state } = this.context;
    getAllGroupVersionComments(
      this.props.moduleCode,
      state.user.credentials[this.props.moduleCode],
      this.props.version[0]
    )(this.updateCommentHistory);
  }

  updateExtractJson(extract) {
    this.setState({ extract: extract });
  }

  componentDidMount() {
    const { state } = this.context;
    this.userHandle = state.user.credentials.handle;
    getAssignedExtract(
      this.props.moduleCode,
      "current"
    )(this.updateExtractJson);
    this.poll = setInterval(this.pollAllComments, 1000);
  }

  componentWillUnmount() {
    clearInterval(this.poll);
  }

  getSelectionForComment(lineNum) {
    let selectedText;
    if (window.getSelection) {
      selectedText = window.getSelection();
    } else if (document.getSelection) {
      selectedText = document.getSelection();
    } else if (document.selection) {
      selectedText = document.selection.createRange().text;
    }
    this.setState({
      commentMetaData: {
        lineNum: lineNum,
        quote: selectedText.toString(),
      },
      showCommentBox: true,
    });
  }

  commentHandler(e) {
    if (e.keyCode != 13) {
      return;
    }

    const comment = this.state.commentMetaData;
    comment["comment"] = e.target.value;
    console.log(comment);
    commentOnExtract(this.props.moduleCode, comment);
    this.setState({
      showCommentBox: false,
    });
  }

  getTextComponent(text, key, isHeading) {
    if (isHeading) {
      return (
        <Typography
          key={key}
          variant="h3"
          component="h3"
          onMouseUp={() => this.getSelectionForComment(key)}
          gutterBottom
        >
          {text}
        </Typography>
      );
    } else {
      return (
        <Typography
          key={key}
          variant="body1"
          onMouseUp={() => this.getSelectionForComment(key)}
          gutterBottom
        >
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

  deleteComment(commentIndex) {
    removeCommentByIndex(this.props.moduleCode, commentIndex);
  }

  render() {
    const commentBox = this.state.showCommentBox ? (
      <Grid item>
        <TextField
          label="Comment"
          fullWidth
          autoFocus
          onKeyDown={this.commentHandler}
        />
      </Grid>
    ) : (
      ""
    );
    const commentHistoryList = Object.keys(this.state.commentHistory).map(
      (handle) => {
        const color = stringToColour(handle);
        return this.state.commentHistory[handle]
          .reverse()
          .map((comment, index) => {
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
                  {this.userHandle == handle ? (
                    <Grid item>
                      <IconButton
                        aria-label="delete"
                        onClick={() => this.deleteComment(index)}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Grid>
                  ) : (
                    ""
                  )}
                </Grid>
              </Grid>
            );
          });
      }
    );
    return (
      <div>
        <Grid
          container
          className="extract-grid"
          spacing={2}
          justify="center"
          direction="column"
        >
          <Grid item>
            <Paper
              variant="elevation"
              elevation={2}
              className="login-background"
            >
              {this.parseJSON(this.state.extract)}
            </Paper>
          </Grid>
          {commentBox}
          {commentHistoryList}
        </Grid>
      </div>
    );
  }
}

PeerReview.contextType = Auth;

export default PeerReview;
