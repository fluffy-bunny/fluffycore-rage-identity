import React, { useState } from 'react';
import { CssBaseline } from '@mui/material';
import { Route, Routes } from 'react-router-dom';
import logo from './logo.svg';

import './App.css';
import LoginPage from './pages/login';
function App() {
  return (
    <div className="container">
          <LoginPage />
       </div>
  );
}

export default App;
