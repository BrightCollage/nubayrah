
import { DarkThemeToggle, Flowbite, Button } from "flowbite-react";
import { BookFileInput } from "../Components/BookFileInput";
import { NavBar } from "../Components/NavBar";
import { DefaultFooter } from "../Components/DefaultFooter";

export default function () {
    return (
        <>
            <NavBar></NavBar>
            <BookFileInput></BookFileInput>
            <DefaultFooter></DefaultFooter>
        </>
    )
}
