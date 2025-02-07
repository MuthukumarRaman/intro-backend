# Project
This project demonstrates a simple web application using the Angular framework for the frontend, Go Lang for the backend, and MongoDB as the database. The application allows users to perform basic CRUD (Create, Read, Update, Delete) operations on a collection of data stored in the database.

Prerequisites
Before running the project, ensure you have the following installed:

Node.js and npm (Node Package Manager)
Angular CLI
Go Lang
MongoDB
Installation
Frontend (Angular)
Clone the repository or download the source code.

Open a terminal or command prompt and navigate to the frontend directory.

Run the following command to install the project dependencies: npm i
Clone the repository or download the source code.

Open a terminal or command prompt and navigate to the backend directory.

Ensure you have MongoDB installed and running.

Run the following command to install the required Go packages:go mod download
Configuration
Frontend
In the frontend/src/environments directory, you'll find two environment files: environment.ts and environment.prod.ts. Update the apiUrl property in both files to match the URL where the backend server is running. By default, it is set to http://localhost:8080/api.

Backend
In the backend directory, you'll find a config.go file. Update the mongoURL constant to match your MongoDB connection string.

Running the Application
Frontend
Open a terminal or command prompt and navigate to the frontend directory.

Run the following command to start the Angular development server:ng s
Open your web browser and navigate to http://localhost:4200 to access the application.
Backend
Open a terminal or command prompt and navigate to the backend directory.

Run the following command to start the Go server:go run main.go

Angular and Go Lang Project with MongoDB Backend
This project demonstrates a simple web application using the Angular framework for the frontend, Go Lang for the backend, and MongoDB as the database. The application allows users to perform basic CRUD (Create, Read, Update, Delete) operations on a collection of data stored in the database.

Prerequisites
Before running the project, ensure you have the following installed:

Node.js and npm (Node Package Manager)
Angular CLI
Go Lang
MongoDB
Installation
Frontend (Angular)
Clone the repository or download the source code.

Open a terminal or command prompt and navigate to the frontend directory.

Run the following command to install the project dependencies:

Copy code
npm install
Backend (Go Lang)
Clone the repository or download the source code.

Open a terminal or command prompt and navigate to the backend directory.

Ensure you have MongoDB installed and running.

Run the following command to install the required Go packages:

go
Copy code
go mod download
Configuration
Frontend
In the frontend/src/environments directory, you'll find two environment files: environment.ts and environment.prod.ts. Update the apiUrl property in both files to match the URL where the backend server is running. By default, it is set to http://localhost:8080/api.

Backend
In the backend directory, you'll find a config.go file. Update the mongoURL constant to match your MongoDB connection string.

Running the Application
Frontend
Open a terminal or command prompt and navigate to the frontend directory.

Run the following command to start the Angular development server:

Copy code
ng serve
Open your web browser and navigate to http://localhost:4200 to access the application.

Backend
Open a terminal or command prompt and navigate to the backend directory.

Run the following command to start the Go server:

go
Copy code
go run main.go
Usage
Once the application is running, you can use the web interface to perform CRUD operations on the data stored in the MongoDB database. The application provides forms to create new records, display existing records, update records, and delete records.

Contributing
Contributions are welcome! If you find any issues or want to add new features, please submit an issue or a pull request.

Acknowledgements
This project was created as a simple demonstration of building a web application using Angular, Go Lang, and MongoDB. Feel free to expand upon it and use it as a starting point for your own projects.
