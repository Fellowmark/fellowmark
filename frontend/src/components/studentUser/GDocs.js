import { Component } from 'react';
import Button from '@material-ui/core/Button';
import { getGoogleDocId, updateGoogleDocId, submitGoogleDoc } from '../../actions/moduleActions';
import { gapi } from 'gapi-script';
import { Auth } from '../../context/authContext';
import "./GDocs.css";

var CLIENT_ID = '744114636968-c4kcndfhfejccc3ijsqghv6n7fe5f5id.apps.googleusercontent.com'
var API_KEY = 'AIzaSyC-VDkWcebILsSu89pwoKk69TQm7W6jB_I';

let SCOPES = "https://www.googleapis.com/auth/documents https://www.googleapis.com/auth/drive";

let DISCOVERY_DOCS = [
  'https://docs.googleapis.com/$discovery/rest?version=v1'
];

class GDocs extends Component {
  constructor(props) {
    super(props);
    this.handleGoogleSigninClick = this.handleGoogleSigninClick.bind(this);
    this.updateGoogleDocId = this.updateGoogleDocId.bind(this);
    this.createDoc = this.createDoc.bind(this);
    this.getJSONfromDoc = this.getJSONfromDoc.bind(this);
    this.submitDocument = this.submitDocument.bind(this);
    this.state = {
      signedIn: false,
      documentId: null 
    };
  }

  handleGoogleSigninClick(signedIn) {
    this.setState({signedIn: signedIn});

    if (!signedIn) {
      gapi.auth2.getAuthInstance().signOut();
    } else {
      gapi.auth2.getAuthInstance().signIn();
    }
  }

  createDoc() {
    gapi.client.docs.documents.create({
      title: "User Guide"
    }).then(response => {
      const doc = response.result;
      const documentId = doc["documentId"]
      this.setState({ documentId: documentId });
      updateGoogleDocId(this.props.moduleCode, documentId);
    });
  }

  componentDidMount() {
    gapi.load('client:auth2', () => { this.initClient() });
  }

  submitDocument(docJson) {
    console.log(docJson);
    submitGoogleDoc(this.props.moduleCode, docJson);
  }

  getJSONfromDoc() {
    gapi.client.docs.documents.get({
      documentId: this.state.documentId
    }).then((response) => {
      var doc = response.result.body;
      this.submitDocument(doc);
    });
  }

  updateGoogleDocId(documentId) {
    this.setState({documentId: documentId});
  }

  initClient() {
    gapi.client.init({
      apiKey: API_KEY,
      clientId: CLIENT_ID,
      discoveryDocs: DISCOVERY_DOCS,
      scope: SCOPES
    }).then(() => {
      this.setState({signedIn: gapi.auth2.getAuthInstance().isSignedIn});
      getGoogleDocId(this.props.moduleCode)(this.updateGoogleDocId, this.createDoc);
    });
  }

  render() {
    const gsigninState = this.state.signedIn ? "Log out" : "Authorize";
    const submitButton = this.state.documentId ? <Button variant="contained" size="large" onClick={this.getJSONfromDoc}>Submit</Button> : "";
    const gDocComponent = this.state.documentId ? <iframe className="editor" src={"https://docs.google.com/document/d/" + this.state.documentId + "/edit?usp=sharing"} /> : "";
    return (
      <div>
        <Button
          variant="contained"
          size="large"
          onClick={() => this.handleGoogleSigninClick(!this.state.signedIn)}
        >
          {gsigninState}
        </Button>
        {submitButton}
        {gDocComponent}
      </div>
    );
  }
}

export default GDocs;
