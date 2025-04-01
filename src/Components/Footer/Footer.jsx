import React,{useEffect} from 'react';
import './footer.css';
import video from '../../Assets/video.mp4'
import { FiSend } from "react-icons/fi";
import { MdOutlineTravelExplore } from "react-icons/md";
import { AiOutlineTwitter } from "react-icons/ai";
import { AiFillInstagram } from "react-icons/ai";
import { AiFillYoutube } from "react-icons/ai";
import { FaTripadvisor } from "react-icons/fa";
import { FiChevronRight } from "react-icons/fi";

import Aos from 'aos';
import 'aos/dist/aos.css';



const Footer = () => {

  useEffect(() => {
      Aos.init({duration: 2000})
    }, [])

  return (
    <section className='footer'>
      <div className='videoDiv'>
      <video src={video} loop autoPlay muted type="video/mp4" />
      </div>

      <div className='secContent container'>
        <div className='contactDiv flex'>
          <div data-aos="fade-up" className="text">
            <small> KEEP IN TOUCH</small>
            <h2>Big things are coming your way! ...</h2>
          </div>

          <div data-aos="fade-up" className="inputDiv flex">
            <input type="text" placeholder='Enter Email Address' />
            <button className="btn flex" type='submit'>
              SEND <FiSend className='icon'/>
            </button>
          </div>
        </div>

        <div className="footerCard flex">
          <div className="footerIntro flex">
            <div className="logoDiv">
              <a href='#' className='logo flex'>
              <MdOutlineTravelExplore className='icon'/>
                Travel.
              </a>
            </div>

            <div data-aos="fade-up" className="footerParagraph">
              It had been six years since dinosaurs returned to Earth. 
              Dust clouded the atmosphere. In some parts of Asia, lava flowed 
              freely again. A scientist, Dr. Lydia Chen, woke up with a pterodactyl
               laying eggs atop her house. Ridiculous but scary. After all, animal
                maternal instincts were fierce.
            </div>

            <div data-aos="fade-up" className="footerSocials flex">
              <AiOutlineTwitter className='icon' />
              <AiFillInstagram className='icon' />
              <AiFillYoutube className='icon' />
              <FaTripadvisor className='icon' />
            </div>
          </div>
          <div  className="footerLinks grid">

              {/* Group one*/}
              <div data-aos="fade-up" className="linkGroup">
                <span className="groupTitle">
                  OUR AGENCY
                </span>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  Services
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  Insurance
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  Agency
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  Tourism
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  Payments
                </li>

              </div>

              {/* Group two*/}
              <div data-aos="fade-up" className="linkGroup">
                <span className="groupTitle">
                  PARTNERS
                </span>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  Bookings
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  Rent Cars
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  HostelWorld
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  Trivago
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  Campus Tour
                </li>

              </div>

              {/* Group three*/}
              <div data-aos="fade-up" className="linkGroup">
                <span className="groupTitle">
                  LAST MINUTE
                </span>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  VACATIONS
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  EVENTS
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  CULTURALS
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  HAPPENING GAINSVILLE
                </li>

                <li className='footerList flex'>
                  <FiChevronRight className='icon'/>
                  HIDDEN GEMS
                </li>

              </div>

          </div>

          <div className="footerDiv flex">
          <small>GATORNEST INN WEBSITE </small>
          <small>TORCHBEARERS</small>
          </div>

        </div>
      </div>

    </section>
  )
}

export default Footer
