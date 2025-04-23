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

INSERT INTO Users (email, password_hash) VALUES
('tarund2302@gmail.com', '$2a$10$eMWMtRpqmmW.Csp6sSQYSOeZaunKfDdvL0lcXjqQc5pHPZEq5xLpm'),
('admin@gmail.com', '$2a$10$wMPF8JFdm3cSej97z3k92.X9fykHqN//e87wMC.9bsmlm7r6gtpJ.');

INSERT INTO Staff (first_name, last_name, email, role) VALUES
('Alice', 'Smith', 'alice@gmail.com', 'Housekeeping'),
('Bob', 'Johnson', 'bob@gmail.com', 'Receptionist'),
('Charlie', 'Williams', 'charlie@gmail.com', 'Housekeeping'),
('Diana', 'Brown', 'diana@gmail.com', 'Housekeeping'),
('Ethan', 'Jones', 'ethan@gmail.com', 'Manager'),
('Fiona', 'Garcia', 'fiona@gmail.com', 'Housekeeping'),
('George', 'Miller', 'george@gmail.com', 'Housekeeping'),
('Hannah', 'Davis', 'hannah@gmail.com', 'Security');

INSERT INTO Rooms (room_id, room_number, room_type, price_per_night, status) VALUES
('101', 'Deluxe', 250.00, 'Available'),
('102', 'Suite', 250.00, 'Available'),
('501', 'Suite', 300.00, 'Available'),
('201', 'Suite', 210.00, 'Available'),
('103', 'Deluxe', 260.00, 'Available'),
('104', 'Deluxe', 275.00, 'Booked'),
('105', 'Deluxe', 255.00, 'Available'),
('106', 'Deluxe', 265.00, 'Maintenance'),
('107', 'Deluxe', 270.00, 'Available'),
('402', 'Single', 95.00, 'Available'),
('403', 'Single', 85.00, 'Booked'),
('404', 'Single', 92.00, 'Available'),
('405', 'Single', 88.00, 'Maintenance'),
('406', 'Single', 100.00, 'Available'),
('202', 'Suite', 220.00, 'Available'),
('203', 'Suite', 230.00, 'Booked'),
('204', 'Suite', 240.00, 'Available'),
('205', 'Suite', 215.00, 'Maintenance'),
('206', 'Suite', 235.00, 'Available'),
('502', 'Suite', 310.00, 'Available'),
('503', 'Suite', 320.00, 'Booked'),
('504', 'Suite', 290.00, 'Available'),
('505', 'Suite', 305.00, 'Maintenance'),
('506', 'Suite', 330.00, 'Available');

