import { DarkThemeToggle, Flowbite, Button } from "flowbite-react";
import { BookFileInput } from "components/BookFileInput";
import { NavBar } from "components/NavBar";
import { DefaultFooter } from "components/DefaultFooter";
import { UploadModal } from "components/UploadModal";

export default function Upload() {
    return (
        <>
            <NavBar></NavBar>
            <UploadModal></UploadModal>
            <DefaultFooter></DefaultFooter>
        </>
    )
}
