import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import Home from './Pages/Home'
import NoPage from './Pages/NoPage'

const App = () => {
  return (
    <>
      <Router>
        <Routes>
          <Route index element={<Home />} />
          <Route path='/home' element={<Home />} />
          <Route path="*" element={<NoPage />} />
        </Routes>
      </Router>
    </>
  );
};

export default App;