import React, {useEffect} from 'react'
import './home.css';
import video from '../../Assets/footerVideo.mp4'
import { GrLocation } from "react-icons/gr";
import { HiFilter } from "react-icons/hi";
import { FiFacebook } from "react-icons/fi";
import { AiOutlineInstagram } from "react-icons/ai";
import { SiTripadvisor } from "react-icons/si";
import { BsListTask } from "react-icons/bs";
import { TbApps } from "react-icons/tb";

import Aos from 'aos';
import 'aos/dist/aos.css';

const Home = () => {

  useEffect(() => {
    Aos.init({duration: 2000})
  }, [])


  return (
    <section className='home'>
      <div className="overlay"></div>
      <video src={video} muted autoPlay loop type="video/mp4" className='video'></video>

      <div className='homeContent container'>
        <div className="textDiv">

          <span data-aos="fade-up" className='smallText'>
            Our Packages
          </span>

          <h1 data-aos="fade-up" className='homeTitle'>
            Search your Hotels
          </h1>

        </div>

        <div data-aos="fade-up" className="cardDiv grid">
          <div className='locationInput'>
            <label htmlFor='type'>Name:</label>
            <div className="input flex">
              <input type='text' placeholder='Enter your first and last name here...' />
              {/* <GrLocation className='icon'/> */}
            </div>
          </div>

          <div className="dateInput">
            <label htmlFor='date'>Select your date:</label>
            <div className="input flex">
              <input type='date'/>
            </div>
          </div>

          <div className="priceInput">
            <div className="label_total flex">
              <label htmlFor='price'>Max price:</label>
              <h3 className='total'>$1000</h3>
            </div>
            <div className="input flex">
              <input type='range'max="1000" min="100"/>
            </div>
          </div>

          <div className="searchOptions flex">
          <HiFilter className='icon'/>
          <span className='span'>CHECK AVAILABILITY</span>

          </div>
        </div>

        <div  className="homeFooterIcons flex">
          <div className="rightIcons">
          <FiFacebook className='icon' href="https://www.facebook.com/reitzunion/"/>
          <AiOutlineInstagram className='icon'/>
          <SiTripadvisor className='icon'/>
          </div>

          <div className="leftIcons">
            <BsListTask className='icon'/>
            <TbApps className='icon'/>
          </div>
        </div>
      </div>


    </section>
  )
}

export default Home
