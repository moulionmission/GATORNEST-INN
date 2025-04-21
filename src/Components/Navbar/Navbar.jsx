import React, { useState } from 'react';
import './navbar.css';
import { FaHotel } from "react-icons/fa6";
import { AiFillCloseCircle } from "react-icons/ai";
import { TbGridDots } from "react-icons/tb";
import { Link } from 'react-router-dom'; // ✅ React Router Link

const Navbar = () => {
  const [active, setActive] = useState('navBar');

  const showNav = () => {
    setActive('navBar activeNavbar');
  };

  const removeNavbar = () => {
    setActive('navBar');
  };

  return (
    <section className='navBarSection'>
      <header className='header flex'>
        <div className='logoDiv'>
          <Link to="/" className='logo flex'>
            <h1><FaHotel className='icon' /> GatorNest - INN</h1>
          </Link>
        </div>

        <div className={active}>
          <ul className="navLists flex">
            <li className="navItem">
              <Link to="/" className="navLink">Home</Link>
            </li>

            <li className="navItem">
              <Link to="/staff" className="navLink">Staff</Link> {/* ✅ Routed to AdminScheduler */}
            </li>

            <li className="navItem">
              <a
                href="https://dineoncampus.com/UF/transact-mobile-ordering"
                className="navLink"
                target="_blank"
                rel="noopener noreferrer"
              >
                Order
              </a>
            </li>

            <li className="navItem">
              <a
                href="https://union.ufl.edu/hotel/"
                className="navLink"
                target="_blank"
                rel="noopener noreferrer"
              >
                About
              </a>
            </li>

            <li className="navItem">
              <Link to="/news" className="navLink">News</Link> {/* Optional route */}
            </li>

            <button className='btn'>
              <Link to="/profile">AVATAR....</Link> {/* Example Avatar/Profile */}
            </button>
          </ul>

          <div onClick={removeNavbar} className='closeNavbar'>
            <AiFillCloseCircle />
          </div>
        </div>

        <div onClick={showNav} className="togglerNavbar">
          <TbGridDots className='icon' />
        </div>
      </header>
    </section>
  );
};

export default Navbar;
