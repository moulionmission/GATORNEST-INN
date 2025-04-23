GATORNEST-INN PROJECT


Description
GatorNest-Inn is a modern, cloud-based solution crafted to simplify and optimize operations for large hotels. Developed with agile methodologies, this system brings together key features like real-time reservation management, guest profiles, staff scheduling, and secure billing, all within an easy-to-use web platform.

By utilizing a scalable cloud infrastructure and a reliable MySQL database via GO Backend API, GatorNest-Inn ensures top-notch performance, security, and dependability. Designed to be intuitive and efficient, this platform not only enhances operational workflows but also elevates the overall guest experience, making it an ideal choice for hotels aiming to grow and thrive in a competitive industry.


# Contributors
Chandra Mouli Dasari,
Gurudeep Paleti,
Harsha Vardhan Reddy Palagiri,
Tarunkrishna Dasari

To run the application, you must run the backend first before running backend

To run backend:
1) You must have a MySQL server to set up the SQL schema, which you will find in [this SQL code](https://github.com/moulionmission/GATORNEST-INN/blob/main/hotel_db_setup.sql)
2) To set up the connection, make sure to use the same connection credentials as the one in [.env file](https://github.com/moulionmission/GATORNEST-INN/blob/main/.env)
3) After setting up MySQL connection, open the connection, then copy and paste the SQL code linked in step 1. Finally, run it to ensure that the necessary data is in SQL
4) Run these commands in the project folder to install GO packages:
   ```
   go mod init hotel-module
   go mod tidy
   ```
5) Run this command to begin running API server locally on port 3000:
   ```
   go run main.go
   ```

To run frontend:
1) Run this command in the project folder to install dependencies:
    ```
    npm install
    ```
2) Then, run this command to launch frontend. If the program asks to run in a different port (typically 3001), type 'Y' to make it proceed:
   ```
   npm start
   ```

Finally, you can now run this program on the given link by nodeJS (which should open automatically and will be deployed on http://localhost:3001/). To login:
 - Use the following credentials to run a regular user account (but you can also register yourself before logging in):
   ```
   Email: tarund2302@gmail.com
   Password: test
   ```
 - Use the following credentails to run the admin account:
   ```
   Email: admin@gmail.com
   Password: test
   ```
