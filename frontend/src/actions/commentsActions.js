import axios from "axios";

export const commentOnExtract = (moduleCode, comment) => {
  axios.post(`/module/${moduleCode}/comment`, comment).catch((err) => {
    console.error(err);
  });
};

export const removeCommentByIndex = (moduleCode, commentIndex) => {
  axios
    .post(`/module/${moduleCode}/comment/remove/${commentIndex}`)
    .catch((err) => {
      console.error(err);
    });
};

export const getAllGroupVersionComments = (moduleCode, groupId, version) => (
  updateCommentHistory
) => {
  axios
    .get(`/module/${moduleCode}/${groupId}/comments/${version}`)
    .then((req) => {
      updateCommentHistory(req.data);
    })
    .catch((err) => {
      console.error(err);
    });
};

export const getAllMarkers = (moduleCode, groupId, version) => (
  updateMarkersList
) => {
  axios
    .get(`/module/${moduleCode}/${groupId}/markers/${version}`)
    .then((res) => {
      updateMarkersList(res.data);
    })
    .catch((err) => {
      console.error(err);
    });
};

export const stringToColour = (str) => {
  var hash = 0;
  for (var i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  var colour = "#";
  for (var i = 0; i < 3; i++) {
    var value = (hash >> (i * 8)) & 0xff;
    colour += ("00" + value.toString(16)).substr(-2);
  }
  return colour;
};
