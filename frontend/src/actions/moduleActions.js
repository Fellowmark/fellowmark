import axios from "axios";

export const createModule = (moduleCode) => {
  axios.post(`/module/create`, { moduleCode: moduleCode }).catch((err) => {
    console.error(err);
  });
};

export const listModules = (updateModuleList) => {
  axios
    .get("/module/list")
    .then((res) => {
      updateModuleList(res.data);
    })
    .catch((err) => {
      console.error(err);
    });
};

export const getModuleInfo = (moduleCode, updateModuleInfo) => {
  axios
    .get(`/module/${moduleCode}`)
    .then((res) => {
      updateModuleInfo(res.data);
    })
    .catch((err) => {
      if (err.response.status === 503) {
        createModule(moduleCode);
        getModuleInfo(moduleCode, updateModuleInfo);
      }
    });
};

export const addUserToModule = (moduleCode, userHandle) => {
  axios.get(`/module/${moduleCode}/add/${userHandle}`).catch((err) => {
    console.error(err);
  });
};

export const updateGroupings = (moduleCode, groupings) => {
  axios.post(`/module/${moduleCode}/groups/update`, groupings).catch((err) => {
    console.error(err);
  });
};

export const updateAssignments = (moduleCode, assignments) => {
  axios
    .post(`/module/${moduleCode}/assignments/update`, assignments)
    .catch((err) => {
      console.error(err);
    });
};

export const updateVersion = (moduleCode, newVersion) => {
  axios.post(`/module/${moduleCode}/version/${newVersion}`).catch((err) => {
    console.error(err);
  });
};

export const getGoogleDocId = (moduleCode) => (
  updateGoogleDocIdState,
  createDoc
) => {
  axios
    .get(`/module/${moduleCode}/doc/id`)
    .then((res) => {
      updateGoogleDocIdState(res.data);
    })
    .catch((err) => {
      if (err.response.status === 504) {
        createDoc();
      }
    });
};

export const updateGoogleDocId = (moduleCode, documentId) => {
  axios
    .post(`/module/${moduleCode}/doc/id`, { documentId: documentId })
    .catch((err) => {
      throw new Error(err);
    });
};

export const submitGoogleDoc = (moduleCode, docJson) => {
  axios.post(`/module/${moduleCode}/doc`, docJson).catch((err) => {
    console.error(err);
  });
};

export const getAssignedExtract = (moduleCode, version) => (
  updateExtractJson
) => {
  axios
    .get(`/module/${moduleCode}/assignment/${version}`)
    .then((res) => {
      updateExtractJson(res.data);
    })
    .catch((err) => {
      console.error(err);
    });
};

export const getOwnDocument = (moduleCode, version) => (updateOwnDoc) => {
  axios
    .get(`/module/${moduleCode}/doc/${version}`)
    .then((res) => {
      updateOwnDoc(res.data.content);
    })
    .catch((err) => {
      console.error(err);
    });
};

export const splitByHeading = (doc) => {
    let currentHeading = null;
    let headingObject = {};
    doc["content"].forEach(object => {
        if (object["paragraph"]) {
            if (object["paragraph"]["paragraphStyle"]["headingId"]) {
                currentHeading = object["paragraph"]["elements"][0]["textRun"]["content"];
                currentHeading = currentHeading.replace("\n", "").toLowerCase();
                headingObject[currentHeading] = [];
                headingObject[currentHeading].push(object);
            } else {
                headingObject[currentHeading].push(object);
            }
        }
    });
    return headingObject;
}

export const getGroupAssignedExtract = (moduleCode, groupId, version) => (
  updateExtract
) => {
  axios
    .get(`/module/${moduleCode}/assignment/${groupId}/${version}`)
    .then((res) => {
      updateExtract(res.data);
    })
    .catch((err) => {
      console.error(err);
    });
};
