import React, {FC, useContext, useState} from "react";
import ReactQuill from "react-quill";
import { Button } from "@material-ui/core";
import "react-quill/dist/quill.snow.css";
import "./QuillEditor.css"
import EditorToolbar, { modules, formats } from "./EditorToolBar";
import {onlineSubmit} from "../actions/moduleActions";
import {AuthContext} from "../context/context";


const QuillEditor: FC<{
    studentId: number
    questionId: number;
}> = (props) => {
    const [submission, setSubmission] = useState({content: ''});
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

        console.log(submission.content)
    }

    return (
        <div>
            <EditorToolbar toolbarId={'t1'} style={{ height: '300px' }}/>
            <ReactQuill
                theme="snow"
                value={submission.content}
                onChange={onChangeContent}
                placeholder={"Enter..."}
                modules={modules('t1')}
                formats={formats}
            />
            <Button
                style={{ marginTop: "10px", marginBottom: "10px", float:"right"}}
                variant="contained"
                onClick={onSubmit}
            >
                Submit
            </Button>
        </div>
    );
};

export default QuillEditor;