
import { useState } from "react";
import axios from 'axios';
import { Alert, Button, FileInput, Label } from "flowbite-react";

export function BookFileInput() {

    const [file, setFile] = useState()
    const [requestResponse, setRequestResponse] = useState()
    function handleChange(event) {
        setFile(event.target.files[0])
    }

    function handleSubmit(event) {
        event.preventDefault()
        const formData = new FormData();
        formData.append("epub", file);
        axios
            .post("http://localhost:5050/books", formData, {
                headers: {
                    "Content-Type": "multipart/form-data",
                },
            })
            .then((response) => {
                // handle the response
                console.log(response);
                setRequestResponse(response);
            })
            .catch((error) => {
                // handle errors
                console.log(error);
                setRequestResponse(error);
            });
    }


    return (
        <div className="grid gap-2">
            <form onSubmit={handleSubmit}>
                <div className="grid grid-cols-5 gap-4">
                    <div className="col-span-4"> <FileInput type="file" id="file-upload" onChange={handleChange} /></div>
                    <div className="grow"><Button type="submit">Submit</Button></div>
                </div>
            </form>
            <div>{Boolean(requestResponse) &&
                <Alert color="info">
                    <span className="font-medium">{requestResponse.status}</span> {requestResponse.statusText}
                </Alert>
            }
            </div>
        </div>
    );
}
