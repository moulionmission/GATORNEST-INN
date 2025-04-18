import React,{useState} from 'react'
import './navbar.css';
import { FaHotel } from "react-icons/fa6";
import { AiFillCloseCircle } from "react-icons/ai";
import { TbGridDots } from "react-icons/tb";

const Navbar = () => {
  const [active,setActive] = useState('navBar')
  
  const showNav = () => {
    setActive('navBar activeNavbar')
  }

  const removeNavbar = () => {
    setActive('navBar')
  }

  return (
    <section className='navBarSection'>
      <header className='header flex'>

        <div className='logoDiv'>
          <a href="#" className='logo flex'>
            <h1> <FaHotel className='icon'/> GatorNest - INN</h1>
          </a>
        </div>

        <div className={active}>
          <ul className="navLists flex">

            <li className="navItem">
              <a href="#" className="navLink">Home</a>
            </li>

            <li className="navItem">
              <a href="#" className="navLink">Staff</a>
            </li>

            <li className="navItem">
              <a href="https://dineoncampus.com/UF/transact-mobile-ordering" className="navLink">Order</a>
            </li>

            <li className="navItem">
              <a href="https://union.ufl.edu/hotel/" className="navLink">About</a>
            </li>

            

            <li className="navItem">
              <a href="#" className="navLink">News</a>
            </li>

            <button className='btn'>
              <a href="#">AVATAR....</a>
            </button>
          </ul>

            <div onClick={removeNavbar} className='closeNavbar'>
            <AiFillCloseCircle />
            </div>

        </div>

        <div onClick={showNav} className="togglerNavbar">
        <TbGridDots className='icon'/>
        </div>
      </header>
    </section>
  )
}

export default Navbar
