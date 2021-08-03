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

export const updateQuestion = (moduleCode, newQuestion) => {
  axios.post(`/module/${moduleCode}/version/${newQuestion}`).catch((err) => {
    console.error(err);
  });
};
