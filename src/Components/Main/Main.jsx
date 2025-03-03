import React from 'react'
import './main.css';
import img1 from '../../Assets/image1.jpg';
import img2 from '../../Assets/image2.jpg';
import img3 from '../../Assets/image3.jpg';
import img4 from '../../Assets/image4.jpg';
import img5 from '../../Assets/image5.jpg';
import img6 from '../../Assets/image6.jpg';
import img7 from '../../Assets/image7.jpg';
import img8 from '../../Assets/image8.jpg';
import img9 from '../../Assets/image9.jpg';
import { HiOutlineLocationMarker } from "react-icons/hi";
import { HiOutlineClipboardCheck } from "react-icons/hi";





const Data = [
  {
    id: 1,
    imgSrc: img1,
    roomTitle: "Deluxe Suite",
    location: "First Floor",
    grade: "Luxury",
    fee: 250,
    description: "A spacious suite featuring modern amenities and a panoramic city view."
  },
  {
    id: 2,
    imgSrc: img2,
    roomTitle: "Executive Room",
    location: "Second Floor",
    grade: "Premium",
    fee: 180,
    description: "A well-appointed room ideal for business travelers with high-speed WiFi."
  },
  {
    id: 3,
    imgSrc: img3,
    roomTitle: "Standard Room",
    location: "Third Floor",
    grade: "Standard",
    fee: 120,
    description: "A comfortable room with all the basic amenities for a pleasant stay."
  },
  {
    id: 4,
    imgSrc: img8,
    roomTitle: "Family Room",
    location: "Ground Floor",
    grade: "Family",
    fee: 200,
    description: "A large room perfect for families, offering extra space and kid-friendly facilities."
  },
  {
    id: 5,
    imgSrc: img5,
    roomTitle: "Single Room",
    location: "Fourth Floor",
    grade: "Economy",
    fee: 90,
    description: "An economical option ideal for solo travelers with all essential amenities."
  },
  {
    id: 6,
    imgSrc: img6,
    roomTitle: "King Suite",
    location: "Fifth Floor",
    grade: "Luxury",
    fee: 300,
    description: "An opulent suite featuring a king-size bed, luxury bath, and a private balcony."
  },
  {
    id: 7,
    imgSrc: img7,
    roomTitle: "Junior Suite",
    location: "Second Floor",
    grade: "Premium",
    fee: 210,
    description: "A modern suite offering extra living space and contemporary design."
  },
  {
    id: 8,
    imgSrc: img4,
    roomTitle: "VIP Room",
    location: "Penthouse",
    grade: "Exclusive",
    fee: 450,
    description: "An exclusive room with premium services, a private lounge, and stunning views."
  },
  {
    id: 9,
    imgSrc: img9,
    roomTitle: "Budget Room",
    location: "Third Floor",
    grade: "Economy",
    fee: 75,
    description: "A cost-effective room with essential comforts and a cozy ambiance."
  },
  
];



const Main = () => {
  return (
    <section className='main container section'>

      <div className="secTitle">
        <h3 className='title'>
          Popular Rooms
        </h3>
      </div>

      <div className="secContent grid">
        {
          Data.map(({id, imgSrc, roomTitle, location, grade, fee, description}) => {
            return(
              <div key={id} className="singleRoom">
              
                <div className="imageDiv">
                  <img src={imgSrc} alt={roomTitle}/>
                </div>

                <div className="cardInfo">
                  <h4 className="roomTitle">{roomTitle}</h4>
                  <span className='continent flex'>
                  <HiOutlineLocationMarker className='icon'/>
                  <span className="name">{location}</span>
                  </span>

                  <div className="fee flex">
                    <div className="grade">
                      <span>{grade}<small>+1</small></span>
                    </div>
                    <div className="price">
                      <h5>{'$' + fee}</h5>
                    </div>
                  </div>

                  <div className="desc">
                    <p>{description}</p>
                  </div>

                  <button className='btn flex'>
                    DETAILS <HiOutlineClipboardCheck className='icon'/>

                  </button>
                </div>
              </div>
            )
          })
        }
      </div>

    </section>
  )
}

export default Main
