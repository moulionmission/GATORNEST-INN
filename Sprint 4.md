Project Documentation - Sprint 4 Updates

Demo Video: https://drive.google.com/file/d/1EYXNnet02pYhrV-Vih9lPD5HEBTSYKe1/view?usp=sharing 
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


 5\. Weekly Schedule:
-------------------------------

Description: Summary Complete week breakdown displays which personnel are rostered on each day
 - Clustered and vertically arranged for clarity and responsiveness.



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
- Paused development on payment system and reallocated efforts into new login/register feature
- Implemented Basic Reservation model and token-based login/registration functionality in SQL
- Link data to authenticated users


- Ran Unit Tests on Backend API endpoints
- Continued the GO SQL Backend API with the frontend, although more of a challenge now

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

These unit tests ensure that each API endpoint is functioning correctly and returning the expected status codes.
