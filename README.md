# Cloud File storage and Sharing service API


## File Storage and Sharing Service

### Overview

Welcome to the File Storage and Sharing Service! This project is designed to streamline file management by leveraging Amazon S3 for secure file storage and signed URLs for controlled access. Whether you're handling documents, images, or PDFs, this service provides a seamless way to upload, share, and retrieve files.

### Key Features

- **Upload Files:** Efficiently upload files to S3 using a straightforward API.
- **Share Files:** Generate signed URLs to securely share files with others.
- **File Retrieval:** Fetch and manage your files with flexible search options.
- **User Management:** Register and authenticate users to manage their files.

### Setup Instructions

#### Prerequisites

- Go (1.23+)
- AWS SDK and Cloud Console (for managing S3)
- Docker (optional, for containerized deployment)

#### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/yourusername/your-repo.git
   cd your-repo

2. **Install Dependencies**
   Ensure you have Go installed. Install the required Go modules:
   ```bash
    go mod download 

3. **Configure AWS**

Make sure your AWS credentials are configured properly. Place your AWS credentials file (~/.aws/credentials) into the Docker container if you are using Docker.

4. **Run the server**
    ```bash
    go run main.go

5. **To containerize and run the server with docker**
    ```bash
    docker build -t file-service .
    docker run -p 8080:8080 -e AWS_ACCESS_KEY_ID=your_access_key -e AWS_SECRET_ACCESS_KEY=your_secret_key file-service

## API Routes

### Ping Service

- **GET** `/ping`
  - Returns a simple "pong" response to check if the service is running.

### Upload File

- **POST** `/upload`
  - Uploads a file to S3 and creates a record in the database.
  - **Form Data:**
    - `file` (file): The file to upload.
    - `userId` (string): The ID of the user uploading the file.
    - `bucketKey` (string): The key to identify the file in S3.
    - `contentType` (string): The MIME type of the file (e.g., `application/pdf`, `image/jpeg`).

### Get Files by User ID

- **GET** `/files/:userId`
  - Retrieves a list of files for a specific user.
  - **Path Parameter:**
    - `userId` (string): The ID of the user.

### Share File

- **GET** `/share/:file_id`
  - Generates a signed URL for sharing a specific file.
  - **Path Parameter:**
    - `file_id` (string): The ID of the file.

### User Registration

- **POST** `/register`
  - Registers a new user.
  - **Form Data:**
    - `username` (string): The username of the new user.
    - `password` (string): The password for the new user.

### User Login

- **POST** `/login`
  - Authenticates a user and returns a token.
  - **Form Data:**
    - `username` (string): The username of the user.
    - `password` (string): The password of the user.

### Search Files

- **GET** `/search`
  - Searches for files based on various criteria.
  - **Query Parameters:**
    - `userId` (string): The ID of the user.
    - `fileName` (string): The name of the file.
    - `startDate` (string): The start date for the search range (format: `YYYY-MM-DD`).
    - `endDate` (string): The end date for the search range (format: `YYYY-MM-DD`).
    - `fileType` (string): The type of the file (e.g., `pdf`, `jpg`).
