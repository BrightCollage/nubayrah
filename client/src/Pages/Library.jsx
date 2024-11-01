import { DarkThemeToggle, Flowbite, Button } from "flowbite-react";
import { BookFileInput } from "../Components/BookFileInput";
import { NavBar } from "../Components/NavBar";
import { BookTable } from "../Components/BookTable";
import { DefaultFooter } from "../Components/DefaultFooter";

export default function Library() {
    return (
        <>
            <NavBar></NavBar>
            <BookTable></BookTable>
            <DefaultFooter></DefaultFooter>
        </>
    )
}
