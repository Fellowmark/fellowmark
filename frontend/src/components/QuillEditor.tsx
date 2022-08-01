import React, {useState} from "react";
import ReactQuill from "react-quill";
import { Button } from "@material-ui/core";
import "react-quill/dist/quill.snow.css";
import "./QuillEditor.css"
import EditorToolbar, { modules, formats } from "./EditorToolBar";


const QuillEditor = () => {
    const [submission, setSubmission] = useState({
        content: ''
    });

    const onChangeContent = (value) => {
        setSubmission({
            ...submission,
            content: value
        });
    }

    const onSubmit = async (e) => {
        e.preventDefault();
        e.persist();

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