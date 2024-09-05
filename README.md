# Quiz-Management-System

Welcome to the **Quiz Management System** project! This project is a web-based application designed to facilitate quiz management for teachers and students. Teachers can set questions, and students can participate in quizzes. The project uses Go for the backend, MongoDB as the database, and HTML/CSS for the frontend.

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Setup Instructions](#setup-instructions)
- [Usage](#usage)
- [File Structure](#file-structure)
- [Screenshots](#screenshots)
- [Contributing](#contributing)
- [License](#license)

## Features

- **User Registration & Authentication**: Users can register and log in as either a teacher or a student.
- **Role-Based Access**: Teachers can set quiz questions, and students can participate in quizzes.
- **CSV Upload**: Teachers can upload quiz questions in bulk using CSV files.
- **Quiz Participation**: Students can attempt quizzes, and their answers are stored for later review.
- **Results Viewing**: Teachers can view student quiz results.

## Tech Stack

- **Backend**: Go
- **Frontend**: HTML, CSS, JavaScript
- **Database**: MongoDB
- **Template Engine**: Go's `html/template`
- **Email Service**: SMTP for sending verification codes

## Setup Instructions

### Prerequisites

- [Go](https://golang.org/dl/) (1.16 or higher)
- [MongoDB Atlas](https://www.mongodb.com/cloud/atlas)
- [Git](https://git-scm.com/)
- [Gmail Account](https://mail.google.com/)

### Installation

1. **Clone the repository**:
    ```bash
    git clone https://github.com/abbas0703/quiz-management-system.git
    cd quiz-management-system
    ```

2. **Setup MongoDB**:
    - Create a MongoDB cluster on MongoDB Atlas.
    - Create a database named `quizdb`.
    - Replace the MongoDB URI in `main.go` with your MongoDB connection string.

3. **Environment Setup**:
    - Create a `.env` file in the root directory with your SMTP email configuration.
    ```env
    SMTP_EMAIL=userstest323@gmail.com
    SMTP_PASSWORD=yourpassword
    ```

4. **Run the Application**:
    ```bash
    go run main.go
    ```

5. **Access the Application**:
    - Visit `http://localhost:8080` in your web browser.

## Usage

### For Teachers

1. **Login** as a teacher.
2. Navigate to the **Set Questions** section.
3. Upload a **CSV file** containing questions and answers.
4. View student quiz results under the **View Results** section.

### For Students

1. **Login** as a student.
2. Start the quiz by clicking on **Take Quiz**.
3. Submit your answers and view the results.

## File Structure
   ```bash
quiz-management-system/
│
├── static/ # Static assets (CSS, JS)
├── templates/ # HTML templates
│ ├── login.html # Login page
│ ├── register.html # Registration page
│ ├── student.html # Student dashboard
│ ├── teacher.html # Teacher dashboard
│ ├── setquestions.html # Question upload page
│ └── viewresults.html # Results viewing page
├── main.go # Main application file
├── README.md # Project documentation
└── .env # Environment variables
   ```

## Contributing

Contributions are welcome! Please follow the steps below to contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Commit your changes (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.


