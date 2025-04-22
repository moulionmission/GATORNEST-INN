-- Create the database (if it doesn't exist) and select it.
CREATE DATABASE IF NOT EXISTS hotel;
USE hotel;

-- Drop tables if they already exist to ensure a clean setup.
DROP TABLE IF EXISTS Staff_Schedule;
DROP TABLE IF EXISTS Reviews;
DROP TABLE IF EXISTS Room_Availability;
DROP TABLE IF EXISTS Payments;
DROP TABLE IF EXISTS Reservations;
DROP TABLE IF EXISTS Staff;
DROP TABLE IF EXISTS Rooms;
DROP TABLE IF EXISTS Guests;
DROP TABLE IF EXISTS Users;

-- Users table for storing login credentials
CREATE TABLE Users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Guests Table (linked to Users)
CREATE TABLE Guests (
    guest_id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE NOT NULL,
    user_id INT,
    FOREIGN KEY (user_id) REFERENCES Users(user_id) ON DELETE SET NULL
);

-- Create Rooms Table
CREATE TABLE Rooms (
    room_id INT AUTO_INCREMENT PRIMARY KEY,
    room_number VARCHAR(10) UNIQUE NOT NULL,
    room_type ENUM('Single', 'Double', 'Suite', 'Deluxe') NOT NULL,
    price_per_night DECIMAL(10,2) NOT NULL,
    status ENUM('Available', 'Booked', 'Maintenance') DEFAULT 'Available'
);

-- Create Reservations Table
CREATE TABLE Reservations (
    reservation_id INT AUTO_INCREMENT PRIMARY KEY,
    guest_id INT NOT NULL,
    room_id INT NOT NULL,
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    status ENUM('Pending', 'Confirmed', 'Checked-in', 'Checked-out', 'Cancelled') DEFAULT 'Pending',
    total_price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (guest_id) REFERENCES Guests(guest_id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES Rooms(room_id) ON DELETE CASCADE
);

-- Create Payments Table
CREATE TABLE Payments (
    payment_id INT AUTO_INCREMENT PRIMARY KEY,
    reservation_id INT NOT NULL,
    payment_method ENUM('Credit Card', 'Debit Card', 'PayPal', 'Cash') NOT NULL,
    payment_status ENUM('Pending', 'Completed', 'Failed', 'Refunded') DEFAULT 'Pending',
    amount DECIMAL(10,2) NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (reservation_id) REFERENCES Reservations(reservation_id) ON DELETE CASCADE
);

-- Create Room Availability Table
CREATE TABLE Room_Availability (
    availability_id INT AUTO_INCREMENT PRIMARY KEY,
    room_id INT NOT NULL,
    date DATE NOT NULL,
    status ENUM('Available', 'Booked', 'Maintenance') DEFAULT 'Available',
    FOREIGN KEY (room_id) REFERENCES Rooms(room_id) ON DELETE CASCADE
);

-- Create Reviews Table
CREATE TABLE Reviews (
    review_id INT AUTO_INCREMENT PRIMARY KEY,
    guest_id INT NOT NULL,
    reservation_id INT NOT NULL,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    review_text TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (guest_id) REFERENCES Guests(guest_id) ON DELETE CASCADE,
    FOREIGN KEY (reservation_id) REFERENCES Reservations(reservation_id) ON DELETE CASCADE
);

-- Create Staff Scheduling Table
CREATE TABLE Staff_Schedule (
    schedule_id INT AUTO_INCREMENT PRIMARY KEY,
    staff_id INT NOT NULL,
    shift_date VARCHAR(50) NOT NULL,
    shift_time ENUM('Morning', 'Afternoon', 'Night') NOT NULL,
    shift_date VARCHAR(50) NOT NULL,
    shift_time ENUM('Morning', 'Afternoon', 'Night') NOT NULL,
    FOREIGN KEY (staff_id) REFERENCES Staff(staff_id) ON DELETE CASCADE
);

-- Create Staff Table
CREATE TABLE Staff (
    staff_id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    role ENUM('Manager', 'Receptionist', 'Housekeeping', 'Security') NOT NULL
);

INSERT INTO Staff (first_name, last_name, email, role) VALUES
('Alice', 'Smith', 'alice@gmail.com', 'Housekeeping'),
('Bob', 'Johnson', 'bob@gmail.com', 'Receptionist'),
('Charlie', 'Williams', 'charlie@gmail.com', 'Housekeeping'),
('Diana', 'Brown', 'diana@gmail.com', 'Housekeeping'),
('Ethan', 'Jones', 'ethan@gmail.com', 'Manager'),
('Fiona', 'Garcia', 'fiona@gmail.com', 'Housekeeping'),
('George', 'Miller', 'george@gmail.com', 'Housekeeping'),
('Hannah', 'Davis', 'hannah@gmail.com', 'Security');

