# Roots of Flavor
## Project Description

"Roots of Flavor" is a web application designed to help food enthusiasts discover, explore, and learn about a variety of foods from around the world. The platform allows users to search for food items by name and access detailed information, including descriptions, images, and the countries where the foods originate. The project aims to connect people with diverse cuisines, promoting cultural understanding and culinary exploration.

Target Audience:
This application is targeted at food enthusiasts, travelers, home cooks, or anyone looking to discover and learn about different foods from around the world.
### Team Members:
- Kantai Daulet
- Sanzhar Vaisov
- Nurdaulet Kolbai
  ## Screenshot
  ![Main Page Screenshot](screenshot/screenshot1.png)

  ## How to Start the Project

Follow these steps to run the project on your local machine:

### 1. Clone the repository

Clone the project repository to your local machine:

```bash
git clone https://github.com/da3goo/Roots-of-Flavor.git
```
### 2. Set up the Backend (Go)
- Clone the repository.
- Navigate to the `Backend` directory.
- Install Go (version >= 1.18).
- Run `go mod tidy` to install dependencies.
- Set up your PostgreSQL database and configure the connection string in the `init` function or use existing online database
- Start the server by running:
  ```bash
  go run main.go




  ## Tools and Resources

1. **Go (Golang)**:
   - Programming language used to build the backend.
   - [Official Go website](https://golang.org/)

2. **PostgreSQL**:
   - Relational database used to store food information.
   - [Official PostgreSQL website](https://www.postgresql.org/)

3. **JavaScript**:
   - Programming language used to create the frontend interactivity.

4. **HTML/CSS**:
   - Markup and styling languages for building and designing the user interface.

5. **Postman**:
   - Tool for testing and documenting the API.
   - [Official Postman website](https://www.postman.com/)

6. **Supabase**:
   - Backend-as-a-Service platform providing database, authentication, and file storage services.
   - [Official Supabase website](https://supabase.com/)

7. **Visual Studio Code**:
   - Popular code editor supporting Go, JavaScript, HTML, CSS, and many extensions.
   - [Download Visual Studio Code](https://code.visualstudio.com/)

8. **Git and GitHub**:
   - Version control system for tracking changes in the project.
   - [Official Git website](https://git-scm.com/)
   - [GitHub](https://github.com/) for hosting and collaboration on the project.

9. **Logrus**:
   - Structured logger for Go, used for advanced logging features.
   - [Official Logrus GitHub](https://github.com/sirupsen/logrus)

10. **SMTP**:
   - Protocol used for sending emails in the project, especially for the email functionalities.
   
11. **Bcrypt**:
   - A Go package for hashing passwords securely.
   - [Official Bcrypt GitHub](https://github.com/golang/crypto/tree/master/bcrypt)

12. **Base64**:
   - Encoding/decoding method used for handling data in textual form, often used for encoding binary data in email communication.



## Update history

- **v1.0**: Realese
- **v1.0.1**: Bugs fixed
- **v1.1**: Added Registrationm and profile handler
- **v1.1.1**: Added delete function and fixed some bugs
- **v1.1.2**: More optimized
- **v1.2**: Realese!
- **v1.2.1**: Redirection added, and fixed bugs
- **v1.3**: Added new Admin Page, and functions for it.Admin can view the data by sorting, and filtering
- **v1.3.1**: More optimized, and fixed bugs
- **v1.4.**: Added pagination.
- **v1.4.1**: Added logrus support for methods
- **v1.4.2**: Rate limiting for login method. Secure from broot force attacks!
- **v1.4.3**: More secure! Added hashing passwords. Also now , you can change your password and emails
- **v1.4.4**: Error handlings added. More logrus support
- **v1.5**: Added email sending feature. You can now contact us with real google account email!
- **v1.5.1**: More optimized












