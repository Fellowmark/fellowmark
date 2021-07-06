// Client ID and API key from the Developer Console
var CLIENT_ID = '744114636968-c4kcndfhfejccc3ijsqghv6n7fe5f5id.apps.googleusercontent.com'
var API_KEY = 'AIzaSyC-VDkWcebILsSu89pwoKk69TQm7W6jB_I';

// Array of API discovery doc URLs for APIs used by the sample
var DISCOVERY_DOCS = [
    'https://docs.googleapis.com/$discovery/rest?version=v1'
];

// Authorization scopes required by the API; multiple scopes can be
// included, separated by spaces.
var SCOPES = "https://www.googleapis.com/auth/documents https://www.googleapis.com/auth/drive";

var authorizeButton = document.getElementById('authorize-button');
var signoutButton = document.getElementById('signout-button');

/**
 *  On load, called to load the auth2 library and API client library.
 */
function handleClientLoad(moduleCode, groupNumber) {
    gapi.load('client:auth2', () => { initClient(moduleCode, groupNumber) });
}

/**
 *  Initializes the API client library and sets up sign-in state
 *  listeners.
 */
function initClient(moduleCode, groupNumber) {
    gapi.client.init({
        apiKey: API_KEY,
        clientId: CLIENT_ID,
        discoveryDocs: DISCOVERY_DOCS,
        scope: SCOPES
    }).then(function () {
        // Listen for sign-in state changes.
        gapi.auth2.getAuthInstance().isSignedIn.listen((isSignedIn) => {
            updateSigninStatus(isSignedIn, moduleCode, groupNumber);
        });

        // Handle the initial sign-in state.
        updateSigninStatus(gapi.auth2.getAuthInstance().isSignedIn.get(), moduleCode, groupNumber);
        authorizeButton.onclick = handleAuthClick;
        signoutButton.onclick = handleSignoutClick;
    });
}

/**
 *  Called when the signed in status changes, to update the UI
 *  appropriately. After a sign-in, the API is called.
 */
function updateSigninStatus(isSignedIn, moduleCode, groupNumber) {
    if (isSignedIn) {
        let id;
        firebase.firestore().collection("modules").doc(moduleCode).get().then(doc => {
            const documentId = doc.data()[groupNumber + "doc"];
            if (documentId) {
                id = documentId;
                document.getElementById("editor").src = "https://docs.google.com/document/d/" + documentId + "/edit?usp=sharing";
                console.log("Document id already present");
            } else {
                createDoc().then(documentId => {
                    id = documentId;
                    docId = {};
                    docId[groupNumber + "doc"] = documentId
                    firebase.firestore().collection("modules").doc(moduleCode).update(docId);
                });
            }
        });
        console.log("group number: " + groupNumber);
        authorizeButton.style.display = 'none';
        document.getElementById("submit").style.display = "inline";
        signoutButton.style.display = 'inline';
        document.getElementById("submit").addEventListener("click", () => {
            submitDoc(id, groupNumber);
        });


        viewExtract(moduleCode, groupNumber);

        document.getElementById("view").addEventListener("click", () => {
            document.getElementById("review").hidden = false;
            document.getElementById("docComponents").hidden = true;
            document.getElementById("view").classList.add("active");
            document.getElementById("edit").classList.remove("active");
        });

        document.getElementById("edit").addEventListener("click", () => {
            document.getElementById("docComponents").hidden = false;
            document.getElementById("review").hidden = true;
            document.getElementById("edit").classList.add("active");
            document.getElementById("view").classList.remove("active");
        });

    } else {
        authorizeButton.style.display = 'inline';
        signoutButton.style.display = 'none';
        document.getElementById("submit").style.display = "none";
    }
}

function viewExtract(moduleCode, groupNumber) {
    firebase.firestore().collection("modules").doc(moduleCode).get().then((doc) => {
        const assignee = doc.data().assignments[groupNumber + "Assignee"];
        const extract = doc.data().assignments[groupNumber + "Extract"];
        console.log(assignee);
        console.log(extract);
        firebase.firestore().collection("modules").doc(moduleCode).collection("" + assignee).doc("gDoc").get().then(result => {
            appendPre("\n\n", false);
            parseJSON(splitByHeading(result.data())[extract]);
        });

        document.getElementById("extract").onmouseup = () => {
            getSelection();
        }

        document.getElementById("extract").addEventListener("select", () => {
            getSelection();
        })
    });
}

/**
 *  Sign in the user upon button click.
 */
function handleAuthClick(event) {
    gapi.auth2.getAuthInstance().signIn();
}

/**
 *  Sign out the user upon button click.
 */
function handleSignoutClick(event) {
    gapi.auth2.getAuthInstance().signOut();
}

/**
 * Append a pre element to the body containing the given message
 * as its text node. Used to display the results of the API call.
 *
 * @param {string} message Text to be placed in pre element.
 * @param {boolean} isHeading True if appended text is bold
 */
function appendPre(message, isHeading) {
    let pre = document.getElementById('extract');
    let textContent = null;
    if (isHeading) {
        textContent = document.createElement("B");
        textContent.appendChild(document.createTextNode(message));
    } else {
        textContent = document.createTextNode(message);
    }
    pre.appendChild(textContent);
}

/**
 * Gets the JSON object with a document id and displays it by calling displayDocument
 *
 * @param {string} documentId Id of the document to read
 */
function getJSONfromDoc(documentId) {
    gapi.client.docs.documents.get({
        documentId: documentId
    }).then(function(response) {
        var doc = response.result;
        displayDocument(doc.body);
    },function(response) {
        appendPre('Error: ' + response.result.error.message);
    });
}

function submitDoc(documentId, groupNumber) {
    gapi.client.docs.documents.get({
        documentId: documentId
    }).then(function(response) {
        var doc = response.result;
        console.log(doc)
        firebase.firestore().collection("modules").doc("CS2113T").collection(groupNumber).doc("gDoc").set(doc.body);
        alert("Document submitted");
    },function(response) {
        appendPre('Error: ' + response.result.error.message);
    });
}

/**
 * Creates Google Doc in the users Google drive
 */
async function createDoc() {
    const response = await gapi.client.docs.documents.create({
        title: "User Guide"
    });
    const doc = response.result;
    const documentId = doc["documentId"]
    document.getElementById("editor").src = "https://docs.google.com/document/d/" + documentId + "/edit?usp=sharing";
    return documentId;
}

/**
 * Displays the JSON object
 *
 * @param {object} doc JSON extraction of document
 */
function displayDocument(doc) {
    content = doc["content"];
    parseJSON(content);
}

/**
 * Splits a JSON Doc by heading and returns an object
 *
 * @param {object} doc JSON Doc
 */
function splitByHeading(doc) {
    let currentHeading = null;
    let headingObject = {};
    doc["content"].forEach(object => {
        if (object["paragraph"]) {
            if (object["paragraph"]["paragraphStyle"]["headingId"]) {
                currentHeading = object["paragraph"]["elements"][0]["textRun"]["content"];
                currentHeading = currentHeading.replace("\n", "");
                headingObject[currentHeading] = [];
                headingObject[currentHeading].push(object);
            } else {
                headingObject[currentHeading].push(object);
            }
        }
    });

    return headingObject;
}

/**
 * Parses JSON Doc and displays content using appendPre
 *
 * @param {object} content JSON Doc content to be parsed
 */
function parseJSON(content) {
    content.forEach(object => {
        if (object["sectionBreak"]) {
            appendPre("\n\n", false);
        } else if (object["paragraph"]) {
            if (object["paragraph"]["paragraphStyle"]["headingId"]) {
                object["paragraph"]["elements"].forEach(element => {
                    if (element["textRun"]) {
                        appendPre(element["textRun"]["content"], true);
                    }
                });
            } else {
                object["paragraph"]["elements"].forEach(element => {
                    if (element["textRun"]) {
                        appendPre(element["textRun"]["content"], false);
                    } else if (element["inlineObjectElement"]) {
                        let pre = document.getElementById('extract');
                        let imageElement = document.createElement("iframe");
                        imageElement.src = element["inlineObjectElement"]["textStyle"]["link"]["url"];
                        pre.appendChild(imageElement);
                    }
                });
            }
        } 
    });
}

/**
 *
 *
 */
function getCurrentLine() {
    var selection = window.document.selection;
    var range = window.document.range;
    alert("hi");
    console.log("Selection details are as follows:")
    console.log(selection)
    console.log(range)
}

function sendComment(moduleCode, groupNumber, currGroup) {
    const db = firebase.firestore();
    db.collection("modules").doc(moduleCode).collection(groupNumber).doc(currGroup).set(
        
    );
}
