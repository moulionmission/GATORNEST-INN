import React, { useState } from 'react'
import "./app.css"
import Navbar from './Components/Navbar/Navbar'
import Home from './Components/Home/Home'
import Main from './Components/Main/Main'
import Footer from './Components/Footer/Footer'

import Login from './Components/Login/Login';
import Register from './Components/Register/Register';

const App = () => {
  const [user, setUser] = useState(null);
  const [isRegistered, setIsRegistered] = useState(true);

  const handleLogin = (data) => {
    console.log('User logged in:', data);
    setUser(data);
  };

  const handleRegister = (data) => {
    console.log('User registered:', data);
    setIsRegistered(true);
  };

  const handleLogout = () => {
    setUser(null);
  };

  return (
    <>
      <Navbar />
      {user ? (
        <>
          <Home />
          <Main />
          <button onClick={handleLogout} className="logout-btn">
            Logout
          </button>
          {/* <Footer /> */}
        </>
      ) : (
        <div className="auth-wrapper">
          {isRegistered ? (
            <Login onToggle={() => setIsRegistered(false)} onLogin={handleLogin} />
          ) : (
            <Register onToggle={() => setIsRegistered(true)} onRegister={handleRegister} />
          )}

        </div>
      )}
    </>
  );
};

export default App
