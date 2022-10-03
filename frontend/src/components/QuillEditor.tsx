import React, {FC, useState, useEffect} from "react";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import { Button } from "@material-ui/core";
import "./QuillEditor.css"
import EditorToolbar, { modules, formats } from "./EditorToolBar";
import {getOnlineSubmission, onlineSubmit, onlineUpdate} from "../actions/moduleActions";


const QuillEditor: FC<{
    studentId: number
    questionId: number;
}> = (props) => {
    const [submission, setSubmission] = useState({content: ''});
    const [submitted, setSubmitted] = useState(false);
    const { studentId, questionId } = props;
    const onChangeContent = (value) => {
        setSubmission({
            content: value
        });
    }

    const onSubmit = async (e) => {
        e.preventDefault();
        e.persist();
        try {
            await onlineSubmit(studentId, questionId, submission.content);
        } catch(e) {
            alert("Fail to submit");
        }
        setSubmitted(true);
    }

    const onUpdate = async (e) => {
        e.preventDefault();
        e.persist();
        try {
            await onlineUpdate(studentId, questionId, submission.content);
        } catch(e) {
            alert("Fail to update");
        }
    }

    useEffect(() => {
        getOnlineSubmission(studentId, questionId, setSubmission, setSubmitted);
    }, []);

    return (
        <div>
            <EditorToolbar toolbarId={'t1'} style={{ height: '300px' }}/>
            <ReactQuill
                theme="snow"
                value={submission.content || ''}
                onChange={onChangeContent}
                placeholder={"Enter..."}
                modules={modules('t1')}
                formats={formats}
            />

            {!submitted && <Button
            style={{ marginTop: "10px", marginBottom: "10px", float:"right"}}
            variant="contained"
            onClick={onSubmit}
        >
            SUBMIT
            </Button>}

            {submitted && <Button
                style={{ marginTop: "10px", marginBottom: "10px", float:"right"}}
                variant="contained"
                onClick={onUpdate}
            >
                UPDATE
            </Button>}
        </div>
    );
};

export default QuillEditor;