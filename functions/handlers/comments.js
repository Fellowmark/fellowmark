const { admin, db } = require("../util/admin");

exports.commentOnExtract = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const commentDetails = req.body;
  db.doc(`/modules/${moduleCode}/`)
    .get()
    .then((doc) => {
      const version = doc.data().docVersion;
      if (version) {
        const versionedComment = {};
        versionedComment[version] = admin.firestore.FieldValue.arrayUnion(
          commentDetails
        );
        db.doc(
          `/modules/${moduleCode}/${req.user[moduleCode]}/${req.user.handle}`
        )
          .update(versionedComment)
          .then(() => {
            return res.status(203).send("Comment added");
          })
          .catch((_) => {
            db.doc(
              `/modules/${moduleCode}/${req.user[moduleCode]}/${req.user.handle}`
            ).set(versionedComment);
            return res.status(203).send("Comment added");
          });
      } else {
        throw new Error("Version is not set");
      }
    })
    .catch((err) => {
      return res.status(400).json(err);
    });
};

exports.removeCommentByIndex = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const commentIndex = req.params.commentIndex;
  console.log(req.user[moduleCode])
  console.log(req.user.handle)
  db.doc(`/modules/${moduleCode}/`)
    .get()
    .then((doc) => {
      const version = doc.data().docVersion;
      if (version) {
        db.doc(
          `/modules/${moduleCode}/${req.user[moduleCode]}/${req.user.handle}`
        )
          .get()
          .then((commentsDoc) => {
            const comments = commentsDoc.data();
            comments[version].splice(commentIndex, 1);
            console.log(comments);
            db.doc(
              `/modules/${moduleCode}/${req.user[moduleCode]}/${req.user.handle}`
            ).set(comments);
            return res.status(203).send("Comment removed");
          })
          .catch((_) => {
            return res.status(403).send("Couldn't access doc");
          });
      } else {
        throw new Error("Version is not set");
      }
    })
    .catch((err) => {
      return res.status(400).json(err);
    });
};

exports.getAllGroupVersionComments = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const version = req.params.docVersion;
  const groupId = req.params.groupId;
  const comments = {};

  db.collection(`/modules/${moduleCode}/${groupId}`)
    .get()
    .then((docs) => {
      docs.forEach((doc) => {
        if (doc.id != "gDoc") {
          try {
            comments[doc.id].push(...doc.data()[version]);
          } catch (err) {
            comments[doc.id] = [];
            comments[doc.id].push(...doc.data()[version]);
          }
        }
      });
      return res.json(comments);
    })
    .catch((err) => {
      return res.status(404).json(err);
    });
};

exports.getAllMarkers = (req, res) => {
  const moduleCode = req.params.moduleCode;
  const groupId = req.params.groupId;
  const version = req.params.docVersion;
  const markers = {};
  db.doc(`/modules/${moduleCode}`)
    .get()
    .then((doc) => {
      const assignments = doc.data().assignments[version];
      Object.keys(assignments).forEach((key) => {
        const assignment = assignments[key];
        if (assignment["Assignee"] === groupId) {
          try {
            markers[assignment["Extract"]].push(key);
          } catch (err) {
            markers[assignment["Extract"]] = [];
            markers[assignment["Extract"]].push(key);
          }
        }
      });
      return res.json(markers);
    })
    .catch((err) => {
      return res.status(404).json(err);
    });
};
