exports.splitByHeading = (doc) => {
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
