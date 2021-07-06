const { db } = require("../util/admin");
const { splitByHeading } = require("../util/gDocParser");

exports.createModule = (req, res) => {
  if (req.user.status == "Module Admin") {
    db.doc(`/modules/${req.body.moduleCode}`).set({ admin: req.user.handle });
    return res.status(205).send("Module created");
  } else {
    return res.status(405).json({ error: "User not admin" });
  }
};

exports.listModules = (req, res) => {
  db.collection("modules")
    .listDocuments()
    .then((docs) => {
      const documentIds = docs.map((it) => it.id);
      return res.status(205).json({ moduleList: documentIds });
    })
    .catch((err) => {
      return res.status(405).json(err);
    });
};

exports.getModuleData = (req, res) => {
  db.doc(`/modules/${req.params.moduleCode}`)
    .get()
    .then((doc) => {
      if (doc.exists) {
        const moduleData = doc.data();
        return res.json(moduleData);
      } else {
        return res.status(503).send("Module doesn't exist");
      }
    })
    .catch((err) => {
      return res.status(403).json(err);
    });
};

exports.addUserToModule = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const userId = req.params.handle;
  const groupId = {};
  groupId[moduleCode] = "unassigned";
  db.doc(`/users/${userId}`).update(groupId);
  db.doc(`/modules/${moduleCode}`)
    .get()
    .then((doc) => {
      if (doc.exists) {
        let groupings = doc.data().groupings;
        if (!groupings) {
          groupings = {};
          groupings.unassigned = [];
        } else if (!groupings[groupId[moduleCode]]) {
          groupings[groupId[moduleCode]] = [];
        }
        groupings[groupId[moduleCode]].push(userId);
        db.collection("modules").doc(moduleCode).update({ groupings });
        return res.status(202).send("User has been added to the module");
      } else {
        return res.status(502).send("Module doesn't exist");
      }
    });
};

exports.updateGroupings = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const groupings = req.body;
  db.doc(`/modules/${moduleCode}`).update({ groupings });
  Object.keys(groupings).forEach((key) => {
    groupings[key].forEach((studentHandle) => {
      const groupId = {};
      groupId[moduleCode] = key;
      db.doc(`/users/${studentHandle}`).update(groupId);
    });
  });
  return res.status(202).send("Groupings have been updated");
};

exports.getGroupings = (req, res) => {
  const moduleCode = req.params.moduleCode;
  db.doc(`/modules/${moduleCode}`)
    .get()
    .then((doc) => {
      if (doc.exists) {
        const groupings = doc.data().groupings;
        return res.json(groupings);
      }
    });
};

exports.updateAssignments = (req, res) => {
  const assignments = req.body;
  const moduleCode = req.params.moduleCode;
  db.collection("modules").doc(moduleCode).update({ assignments });
  return res.status(203).send("Assignments have been updated");
};

exports.getAllAssignments = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const version = req.params.docVersion;
  db.doc(`/modules/${moduleCode}`)
    .get()
    .then((doc) => {
      if (doc.exists) {
        const assignments = doc.data().assignments[version];
        return res.json(assignments);
      }
    });
};

exports.getGroupGoogleDocId = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const groupId = req.user[moduleCode];
  if (groupId == "unassigned") {
    return res.status(403).send("You are unassigned");
  }
  db.doc(`/modules/${moduleCode}`)
    .get()
    .then((doc) => {
      const data = doc.data();
      if (data[groupId + "doc"]) {
        return res.status(203).json(data[groupId + "doc"]);
      } else {
        return res.status(504).send("Doc doesn't exist");
      }
    });
};

exports.updateGroupGoogleDocId = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const groupId = req.user[moduleCode];
  if (groupId == "unassigned") {
    return res.status(503).send("You are unassigned");
  }
  const documentId = req.body.documentId;
  const doc = {};
  doc[groupId + "doc"] = documentId;
  db.doc(`/modules/${moduleCode}`).update(doc);
  return res.status(203).send("Created a google doc");
};

exports.updateGroupGoogleDoc = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const groupId = req.user[moduleCode];
  const gDoc = req.body;
  db.doc(`/modules/${moduleCode}`)
    .get()
    .then((doc) => {
      const version = doc.data().docVersion;
      if (version) {
        const versionedDoc = {};
        versionedDoc[version] = gDoc;
        db.collection("modules")
          .doc(moduleCode)
          .collection(groupId)
          .doc("gDoc")
          .update(versionedDoc)
          .catch((err) => {
            db.collection("modules")
              .doc(moduleCode)
              .collection(groupId)
              .doc("gDoc")
              .set(versionedDoc);
          });
        return res.status(203).send("Google Doc content updated");
      } else {
        return res.status(503).send("Version not set");
      }
    });
};

exports.getOwnExtract = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const groupId = req.user[moduleCode];
  let version = req.params.docVersion;
  db.doc(`/modules/${moduleCode}`)
    .get()
    .then((doc) => {
      try {
        if (req.params.docVersion == "current") {
          version = doc.data().docVersion;
        }
        const assignments = doc.data().assignments[version];
        const assignee = assignments[groupId]["Assignee"];
        const extractHeading = assignments[groupId]["Extract"];

        db.doc(`/modules/${moduleCode}/${assignee}/gDoc`)
          .get()
          .then((doc) => {
            const gDoc = doc.data()[version];
            if (extractHeading == "ENTIRE") {
              return res.status(200).json(gDoc);
            } else {
              const extract = splitByHeading(gDoc)[
                extractHeading.toLowerCase()
              ];
              return res.status(200).json(extract);
            }
          })
          .catch((err) => {
            console.error(err);
            return res.status(500).send("Something went wrong");
          });
      } catch (err) {
        return res.status(500).send("Something went wrong");
      }
    });
};

exports.getExtract = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const groupId = req.params.groupId;
  let version = req.params.docVersion;
  console.log(moduleCode);
  console.log(groupId);
  console.log(version);
  console.log(req.user.status);
  if (req.user.status == "Student") {
    return res.status(403).send("Unauth");
  }
  db.doc(`/modules/${moduleCode}`)
    .get()
    .then((doc) => {
      try {
        if (req.params.docVersion == "current") {
          version = doc.data().docVersion;
        }
        const assignments = doc.data().assignments[version];
        const assignee = assignments[groupId]["Assignee"];
        const extractHeading = assignments[groupId]["Extract"];

        db.doc(`/modules/${moduleCode}/${assignee}/gDoc`)
          .get()
          .then((doc) => {
            const gDoc = doc.data()[version];
            if (extractHeading == "ENTIRE") {
              return res.status(200).json(gDoc);
            } else {
              const extract = splitByHeading(gDoc)[
                extractHeading.toLowerCase()
              ];
              return res.status(200).json(extract);
            }
          })
          .catch((err) => {
            console.error(err);
            return res.status(500).send("Something went wrong");
          });
      } catch (err) {
        return res.status(500).send("Something went wrong");
      }
    });
};

exports.updateVersion = (req, res) => {
  db.doc(`/modules/${req.params.moduleCode}`).update({
    docVersion: req.params.docVersion,
  });
  return res.status(203).send("Project version has been updated");
};

exports.getOwnDocument = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const groupId = req.user[moduleCode];
  let version = req.params.docVersion;
  db.doc(`/modules/${moduleCode}`)
    .get()
    .then((doc) => {
      try {
        if (req.params.docVersion == "current") {
          version = doc.data().docVersion;
        }
        db.doc(`/modules/${moduleCode}/${groupId}/gDoc`)
          .get()
          .then((doc) => {
            const gDoc = doc.data()[version];
            return res.status(200).json(gDoc);
          })
          .catch((err) => {
            console.error(err);
            return res.status(500).send("Something went wrong");
          });
      } catch (err) {
        return res.status(500).send("Something went wrong");
      }
    });
};
