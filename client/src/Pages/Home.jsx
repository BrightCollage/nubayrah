'use client'

import { DarkThemeToggle, Flowbite, Button } from "flowbite-react";
import { BookFileInput } from "../Components/BookFileInput";

export default function () {
    return (
        <>
            <main className="flex flex-col min-h-screen dark:bg-gray-800">
                <DarkThemeToggle className="flex self-end" />
                <BookFileInput></BookFileInput>
            </main>
        </>
    )
}
