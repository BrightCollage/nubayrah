import { DarkThemeToggle, Flowbite, Button } from "flowbite-react";
import { BookFileInput } from "../Components/BookFileInput";
import { NavBar } from "../Components/NavBar";
import { DefaultFooter } from "../Components/DefaultFooter";
import { UploadModal } from "../Components/UploadModal";

export default function Upload() {
    return (
        <>
            <NavBar></NavBar>
            <UploadModal></UploadModal>
            <DefaultFooter></DefaultFooter>
        </>
    )
}
