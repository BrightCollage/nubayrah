import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import { DarkThemeToggle, Flowbite, Button } from "flowbite-react";
import { BookFileInput } from "./components/BookFileInput";

function App() {

  return (
    <main className="flex flex-col min-h-screen dark:bg-gray-800">
      <DarkThemeToggle className="flex self-end" />
      <BookFileInput></BookFileInput>
    </main>
  )
}

export default App
