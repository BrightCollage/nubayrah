import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Home from './Pages/Home'
import NoPage from './Pages/NoPage'
import Upload from "./Pages/Upload";
import Library from "./Pages/Library";

const App = () => {
  return (
    <>
      <main className="flex flex-col min-h-screen dark:bg-gray-800">
        <Router>
          <Routes>
            <Route index element={<Home />} />
            <Route path='/home' element={<Home />} />
            <Route path='/library' element={<Library />} />
            <Route path='/upload' element={<Upload />} />
            <Route path="*" element={<NoPage />} />
          </Routes>
        </Router>
      </main>
    </>
  );
};

export default App;