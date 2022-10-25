import React from 'react'
import ReactDOM from 'react-dom/client'
import 'bootstrap/dist/css/bootstrap.css';

import {
  BrowserRouter,
  Routes,
  Route,
} from "react-router-dom";


import { NavbarComponent } from './components/NavbarComponent';

import App from './App'
import {About}  from './views/About'
import {Login}  from './views/Login'

import './index.css';

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <BrowserRouter>
      <NavbarComponent />
      <Routes>
        <Route path='/' element={<App />} />
        <Route path='/about' element={<About />} />
        <Route path='/login' element={<Login />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
)
