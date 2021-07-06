const functions = require("firebase-functions");
const app = require("express")();
const FBAuth = require("./util/fbAuth");
const cors = require("cors");
const helmet = require("helmet");
const compression = require("compression");

app.use(helmet());
app.use(cors());
app.use(compression());

const {
  signup,
  login,
  getAuthenticatedUser,
  getUserDetails,
} = require("./handlers/users");

const {
  listModules,
  createModule,
  getModuleData,
  addUserToModule,
  updateGroupings,
  getGroupings,
  updateAssignments,
  getAllAssignments,
  getGroupGoogleDocId,
  updateGroupGoogleDocId,
  updateGroupGoogleDoc,
  updateVersion,
  getOwnExtract, 
  getExtract,
  getOwnDocument,
} = require("./handlers/modules");

const {
  commentOnExtract,
  getAllGroupVersionComments,
  getAllMarkers,
  removeCommentByIndex,
} = require("./handlers/comments");

app.post("/signup", signup);
app.post("/login", login);
app.get("/user/:handle", getUserDetails);
app.get("/user", FBAuth, getAuthenticatedUser);
app.post("/module/create", FBAuth, createModule);
app.get("/module/list", listModules);
app.get("/module/:moduleCode", getModuleData);
app.post("/module/:moduleCode/add/:handle", addUserToModule);
app.post("/module/:moduleCode/groups/update", updateGroupings);
app.post("/module/:moduleCode/assignments/update", updateAssignments);
app.get("/module/:moduleCode/assignments/:docVersion", getAllAssignments);
app.get("/module/:moduleCode/groups", getGroupings);
app.get("/module/:moduleCode/doc/id", FBAuth, getGroupGoogleDocId);
app.post("/module/:moduleCode/doc/id", FBAuth, updateGroupGoogleDocId);
app.post("/module/:moduleCode/doc", FBAuth, updateGroupGoogleDoc);
app.post("/module/:moduleCode/comment", FBAuth, commentOnExtract);
app.post(
  "/module/:moduleCode/comment/remove/:commentIndex",
  FBAuth,
  removeCommentByIndex
);
app.post("/module/:moduleCode/version/:docVersion", updateVersion);
app.get("/module/:moduleCode/assignment/:docVersion", FBAuth, getOwnExtract);
app.get("/module/:moduleCode/assignment/:groupId/:docVersion", FBAuth, getExtract);
app.get(
  "/module/:moduleCode/:groupId/comments/:docVersion",
  getAllGroupVersionComments
);
app.get("/module/:moduleCode/:groupId/markers/:docVersion", getAllMarkers);
app.get("/module/:moduleCode/doc/:docVersion", FBAuth, getOwnDocument);

exports.api = functions.region("asia-southeast2").https.onRequest(app);
