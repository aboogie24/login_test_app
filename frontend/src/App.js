import React from 'react';
import { BrowserRouter as Router, Route, Routes, Link } from 'react-router-dom';
import { RegisterForm, LoginForm } from './components/auth'
import './App.css';

function App() {
  return (
    <Router> 
      <div className="min-h-screen bg-gray-100 flex flex-col items-center justify-center">
        <nav className="mb-8">
          <ul className="flex space-x-4">
            <li>
              <Link to="/login" className="text-blue-500 hover:text-blue-700">Login</Link>
            </li>
            <li>
              <Link to="/register" className="text-blue-500 hover:text-blue-700">Register</Link>
            </li>
          </ul>
        </nav>

        <Routes>
          <Route path="/login" element={<LoginForm />} />
          <Route path="/register" element={<RegisterForm />} />
          <Route path="/" element={<h1 className="text-2xl font-bold">Welcome to our App!</h1>} />
        </Routes>
      </div>
    </Router> 
      
  );
}

export default App;
