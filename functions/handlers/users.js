const { db } = require("../util/admin");
const firebase = require("firebase");
const config = require("../util/config");

firebase.initializeApp(config);
//firebase.auth().useEmulator("http://localhost:9099/");

exports.signup = (req, res) => {
  console.log(req.body)
  const newUser = {
    email: req.body.email,
    password: req.body.password,
    handle: req.body.handle,
    status: req.body.status,
    moduleCode: req.body.moduleCode
  };

  let token, userId;
  db.doc(`/users/${newUser.handle}`).get().then((doc) => {
    if (doc.exists) {
      return res.status(400).json({ handle: "this handle is already taken" });
    } else {
      return firebase.auth().createUserWithEmailAndPassword(newUser.email, newUser.password);
    }
  }).then(data => {
    userId = data.user.uid;
    return data.user.getIdToken(false);
  }).then(idToken => {
    token = idToken;
    const userCredentials = {
      handle: newUser.handle,
      email: newUser.email,
      createdAt: new Date().toISOString(),
      status: newUser.status,
      moduleCode: newUser.moduleCode,
      userId
    }
    return db.doc(`/users/${newUser.handle}`).set(userCredentials);
  }).then(() => {
    return res.status(201).json({ token });
  }).catch((err) => {
    console.error(err);
    if (err.code === "auth/email-already-in-use") {
      return res.status(400).json({ email: "Email is already is use" });
    } else {
      return res
        .status(500)
        .json({ general: "Something went wrong, please try again" });
    }
  });
};

exports.login = (req, res) => {
  const user = {
    email: req.body.email,
    password: req.body.password
  };
  console.log(user);
  firebase.auth().signInWithEmailAndPassword(user.email, user.password).then(data => {
    return data.user.getIdToken(false);
  }).then(token => {
    return res.json({ token });
  }).catch(err => {
    console.error(err)
    return res.status(403).json(err);
  });
}

exports.getUserDetails = (req, res) => {
  let userData = {};
  db.doc(`/users/${req.params.handle}`).get().then(doc => {
    if (doc.exists) {
      userData.user = doc.data();
      return res.json(userData);
    } else {
      return res.status(404).send("User data not found");
    }
  }).catch(err => {
    return res.status(404).json(err);
  });
}

exports.getAuthenticatedUser = (req, res) => {
  let userData = {};
  db.doc(`/users/${req.user.handle}`).get().then(doc => {
    if (doc.exists) {
      if (doc.exists) {
        userData.credentials = doc.data();
        return res.json(userData);
      } else {
        console.error(res.status(404).json({ error: "User not found" }));
      }
    } else {
      return res.status(404).send("User details not found");
    }
  }).catch(err => {
    return res.status(404).json(err);
  });
}
