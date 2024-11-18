import { DarkThemeToggle, Flowbite, Button } from "flowbite-react";
import { BookFileInput } from "components/BookFileInput";
import { NavBar } from "components/NavBar";
import { BookTable } from "components/BookTable";
import { DefaultFooter } from "components/DefaultFooter";

export default function Library() {
    return (
        <>
            <NavBar></NavBar>
            <BookTable></BookTable>
            <DefaultFooter></DefaultFooter>
        </>
    )
}
