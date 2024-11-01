
import { Button, Modal, Select } from "flowbite-react";
import { useState } from "react";
import { BookFileInput } from "./BookFileInput";

export function UploadModal() {
    const [openModal, setOpenModal] = useState(false);

    return (
        <>
            <Button onClick={() => setOpenModal(true)}>Upload</Button>
            <Modal show={openModal} onClose={() => setOpenModal(false)}>
                <Modal.Header>Select File to Upload</Modal.Header>
                <Modal.Body>
                    <BookFileInput></BookFileInput>
                </Modal.Body>
            </Modal>
        </>
    );
}
