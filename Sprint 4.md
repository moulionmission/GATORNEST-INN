Project Documentation - Sprint 4 Updates

Demo Video: https://youtu.be/lao8h6Ym9OY 

___FRONTEND___

Sprint 4 centered on developing the admin panel functionality for GATOR NEST INN. The primary objective was to introduce advanced staff scheduling and room management capabilities so that hotel administrators could control operations effectively via a dynamic interface.

Key Features Implemented:
=========================

1\. Day-Wise Staff Scheduling System:
----------------------------

Description: The admin can allocate more than one staff for any weekday.
 - Time slot reservation available for each task (Morning, Afternoon, Night).
 - Staff names are filled with sample data.
 - Permits mass assignment of staff.

2\. Edit Assigned Staff:
---------------------------------

Description: Any pre-defined staff contribution can be edited inline.
 - Admin can edit staff name and time slot directly from the schedule.
 - Offers save option that directly saves changes.



3\. Remove Assigned Staff:
----------------------------

Description: Admin can remove any assigned staff entry from the weekly schedule.
 - Clean interface with Edit and Delete buttons together for each item.


4\. Room Availability Manager:
-------------------------------

Description: Admin can choose among 9 actual hotel room types (e.g., Deluxe Suite, VIP Room, Budget Room).
 - For every room type, the admin can input the number of available rooms.
 - A dummy "Save Availability" button is included for backend integration.
 - Availability is shown in a readable list beneath the input form.


 5\. Weekly Schedule Summary:
-------------------------------

Description: Complete week breakdown displays which personnel are rostered on each day
 - Clustered and vertically arranged for clarity and responsiveness.

 6\. "Who's Working Today" View:
-------------------------------

Description: A section that dynamically shows staff on shift for the day
 - Automatically shows changes in real time from the schedule


 7\. UI/UX Improvements:
-------------------------------

Description: Clean vertical design with columns for Room Management, Assignment Form, Weekly Schedule, and Today's Staff
 - Styled to prevent overlap with the sticky navigation bar (padding-top: 7rem).
 - Color-coded buttons Edit (yellow), Save (green), and Delete (red).
 - Enhanced readability and spacing for desktop and mobile designs.


Conclusion:
===========
Sprint 4 Benefits Realized Simplified hotel admin tasks from manual entry to interactive scheduling. Enhanced visibility into daily staff management. Simpler and more maintainable codebase with modular styling and states.

*Unit test case for Navbar component*

toggles navbar active class on open and close

renders Navbar with all expected elements

Unit tests screen shots:
![](s3_frontend_test_ss1.png)

=================================================

___BACKEND___

In sprint 4, we have made the following progress for the backend:
- Verified that previous integration between frontend and backend is working and is intact
- Created new endpoints for faciliating staff scheduling system in GO SQL
- Ran Unit Tests on Backend API endpoints
- Continued integrating the frontend for staff scheduling with  GO SQL backend API

## BACKEND API Documentation

### Authentication Endpoints
- **POST /register** - Registers a new user
- **POST /login** - Authenticates a user and returns a JWT token

### Guest Endpoints
- **GET /guests** - Fetch all guests
- **GET /guests/{id}** - Fetch a specific guest by ID
- **POST /guests** - Create a new guest
- **PUT /guests/{id}** - Update guest details
- **DELETE /guests/{id}** - Delete a guest
- **GET /profile** - Retrieves details of the authenticated user

## Payment Endpoints
- **GET /payments** - Fetch all payments
- **GET /payments/{id}** - Fetch a specific payment by ID
- **POST /payments** - Create a new payment
- **PUT /payments/{id}** - Update payment details
- **DELETE /payments/{id}** - Delete a payment

## Room Endpoints
- **POST /reservations/** - Creates room reservation for user
- **GET /reservations/{id}** - get room reservation made by user based on user ID

## Staff Endpoints
- **GET /staffs** - Fetch all staff members
- **GET /staff/{id}** - Fetch a specific staff member by their ID

## Staff Scheduling Endpoints
- **GET /schedules** - Retrieve all staff schedules
- **POST /schedule** - Add a new schedule entry for a staff member
- **PUT /schedule/{id}** - Update an existing schedule entry by staff ID
- **DELETE /schedule/{id}** - Remove a schedule entry by staff ID

# Backend Unit testing

We used GO Unit Testing to test each endpoint and verify actual results with expected results.

**Guests Endpoint Unit Test Cases**

 1. **Get All Guests**
    - **Method:** `GET`
    - **URL:** `/guests`
    - **Expected Status:** `200 OK`
    - **Expected Response:** A list of all guests in JSON format.

2. **Get Guest by ID**
    - **Method:** `GET`
    - **URL:** `/guests/1`
    - **Expected Status:** `200 OK`
    - **Expected Response:** JSON object of the guest with ID `1`.

3. **Create Guest**
    - **Method:** `POST`
    - **URL:** `/guests`
    - **Request Body:**
    ```json
    {
        "name": "John Doe",
        "email": "john@example.com"
    }
    ```
    - **Expected Status:** `201 Created`
    - **Expected Response:** JSON object of the created guest with assigned ID.

4. **Update Guest**
    - **Method:** `PUT`
    - **URL:** `/guests/1`
    - **Request Body:**
    ```json
    {
        "name": "John Smith"
    }
    ```
    - **Expected Status:** `200 OK`
    - **Expected Response:** Updated guest details in JSON format.

5. **Delete Guest**
    - **Method:** `DELETE`
    - **URL:** `/guests/1`
    - **Expected Status:** `204 No Content`
    - **Expected Response:** No response body.

---

**Payments Endpoint Unit Test Cases**

6. **Get All Payments**
    - **Method:** `GET`
    - **URL:** `/payments`
    - **Expected Status:** `200 OK`
    - **Expected Response:** A list of all payments in JSON format.

7. **Get Payment by ID**
    - **Method:** `GET`
    - **URL:** `/payments/1`
    - **Expected Status:** `200 OK`
    - **Expected Response:** JSON object of the payment with ID `1`.

8. **Create Payment**
    - **Method:** `POST`
    - **URL:** `/payments`
    - **Request Body:**
    ```json
    {
        "amount": 100.0,
        "guest_id": 1
    }
    ```
    - **Expected Status:** `201 Created`
    - **Expected Response:** JSON object of the created payment with assigned ID.

9. **Update Payment**
    - **Method:** `PUT`
    - **URL:** `/payments/1`
    - **Request Body:**
    ```json
    {
        "amount": 150.0
    }
    ```
    - **Expected Status:** `200 OK`
    - **Expected Response:** Updated payment details in JSON format.

10. **Delete Payment**
    - **Method:** `DELETE`
    - **URL:** `/payments/1`
    - **Expected Status:** `204 No Content`
    - **Expected Response:** No response body.

---

**Rooms Endpoint Unit Test Cases**

11. **Make Reservation**
    - **Method:** `POST`
    - **URL:** `/reservations/`
    - **Request Body:**
    ```json
    {
        "check_in_date": "2025-03-27",
        "check_out_date": "2025-03-28",
        "email": "jr203@gmail.com",
        "first_name": "Jason",
        "guest_id": 2,
        "last_name": "Robb",
        "reservation_id": 3,
        "room_id": 1,
        "status": "Pending",
        "total_price": 90
    }
    ```
    - **Expected Status:** `200 OK`
    - **Expected Response:** JSON object of the reservation (same as request body) alongside user information.

12. **Get Reservations by ID**
    - **Method:** `GET`
    - **URL:** `/payments/2`
    - **Expected Status:** `200 OK`
    - **Expected Response:** JSON object of all reservations made by the user with User ID `2`.
---

**Staff Endpoints Unit Test Cases**

13. **Get All Housekeeping Staff**
- **Method:** `GET`
- **URL:** `/staffs`
- **Expected Status:** `200 OK`
- **Expected Response:**
  ```json
  [
    {
      "staff_id": 1,
      "first_name": "John"
    },
    {
      "staff_id": 2,
      "first_name": "Maria"
    }
  ]
  ```

14. **Get Staff by ID**
- **Method:** `GET`
- **URL:** `/staff/1`
- **Expected Status:** `200 OK`
- **Expected Response:**
  ```json
  {
    "staff_id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "role": "Housekeeping"
  }
  ```

- **Error Case (Staff Not Found):**
  - **Expected Status:** `404 Not Found`
  - **Expected Response:**
    ```json
    {
      "error": "Staff not found"
    }
    ```
**Staff Scheduling Endpoints Unit Test Cases**

15. **Get All Schedules**
- **Method:** `GET`
- **URL:** `/schedules`
- **Expected Status:** `200 OK`
- **Expected Response:**
  ```json
  [
    {
      "staff_id": 1,
      "shift_date": "Tuesday",
      "shift_time": "Morning"
    },
    {
      "staff_id": 2,
      "shift_date": "Wednesday",
      "shift_time": "Evening"
    }
  ]
  ```

16. **Add to Schedule**
- **Method:** `POST`
- **URL:** `/schedule`
- **Request Body:**
  ```json
  {
    "staff_id": 3,
    "shift_date": "Tuesday",
    "shift_time": "Night"
  }
  ```
- **Expected Status:** `201 Created`
- **Expected Response:**
  ```json
  {
    "staff_id": 3,
    "shift_date": "Tuesday",
    "shift_time": "Night",
    "schedule_id": 10
  }
  ```

- **Error Case (Bad Input):**
  - **Expected Status:** `400 Bad Request`
  - **Expected Response:**
    ```json
    {
      "error": "Invalid input"
    }
    ```

17. **Remove from Schedule**
- **Method:** `DELETE`
- **URL:** `/schedule/3?shift_date=2025-04-24&shift_time=Night`
- **Expected Status:** `200 OK`
- **Expected Response:**
  ```json
  {
    "message": "Schedule removed"
  }
  ```

- **Error Case (Missing Parameters):**
  - **Expected Status:** `400 Bad Request`
  - **Expected Response:**
    ```json
    {
      "error": "Missing shift_date or shift_time query parameter"
    }
    ```

- **Error Case (Schedule Not Found):**
  - **Expected Status:** `404 Not Found`
  - **Expected Response:**
    ```json
    {
      "error": "Schedule not found"
    }
    ```

18. **Update Schedule**
- **Method:** `PUT`
- **URL:** `/schedule/3?shift_date=2025-04-24&new_shift_time=Evening`
- **Expected Status:** `200 OK`
- **Expected Response:**
  ```json
  {
    "message": "Shift time updated"
  }
  ```

- **Error Case (Missing Parameters):**
  - **Expected Status:** `400 Bad Request`
  - **Expected Response:**
    ```json
    {
      "error": "Missing shift_date or new_shift_time in query parameters"
    }
    ```

- **Error Case (Schedule Not Found):**
  - **Expected Status:** `404 Not Found`
  - **Expected Response:**
    ```json
    {
      "error": "Schedule not found"
    }
    ```

These unit tests ensure that each API endpoint is functioning correctly and returning the expected status codes.
