Below is a **comprehensive README** for this project. Feel free to modify any sections as needed for your specific requirements.

---

# Go-Based CDK Project: Serverless Authentication API

This repository contains a **Go-based AWS CDK project** that deploys:
- A DynamoDB table for user data
- A Lambda function (written in Go) providing a simple user registration and login flow
- An API Gateway exposing three routes:
  - **POST /register**  
  - **POST /login**  
  - **GET /protected**  

Users can **register**, **log in**, and then call the protected endpoint by including a **JWT** in the `Authorization` header.

---

## Table of Contents

1. [Project Structure](#project-structure)  
2. [Architecture Overview](#architecture-overview)  
3. [Prerequisites](#prerequisites)  
4. [Setup & Installation](#setup--installation)  
5. [Build & Deploy](#build--deploy)  
6. [Testing the Application](#testing-the-application)  
7. [API Endpoints](#api-endpoints)  
8. [Implementation Details](#implementation-details)  
   1. [Lambda Handlers](#lambda-handlers)  
   2. [DynamoDB Database Operations](#dynamodb-database-operations)  
   3. [JWT Middleware](#jwt-middleware)  
9. [Local Development Notes](#local-development-notes)  
10. [Contributing](#contributing)  
11. [License](#license)

---

## Project Structure

```
└── sibteali786-goCdk/
    ├── go.mod                # Root Go module definition
    ├── go.sum
    ├── cdk.json              # CDK configuration
    ├── go_cdk.go             # Main CDK stack file
    ├── go_cdk_test.go        # CDK tests (optional, commented out for reference)
    ├── README.md             # This README
    ├── lambda/
    │   ├── main.go           # Lambda entry point (switches by path)
    │   ├── app/
    │   │   └── app.go        # Constructs the app with DB + API handler
    │   ├── api/
    │   │   └── api.go        # Handlers for /register and /login
    │   ├── middleware/
    │   │   └── middleware.go # JWT validation middleware
    │   ├── database/
    │   │   └── database.go   # DynamoDB-based user store
    │   ├── types/
    │   │   └── types.go      # RegisterUser, User, JWT creation, password hash check
    │   ├── function.zip      # Zipped Lambda artifact after `make build`
    │   ├── makefile          # Makefile to build and zip the Lambda code
    │   ├── go.mod            # Lambda-specific Go module definition
    │   └── go.sum
    └── ...
```

---

## Architecture Overview

1. **AWS Lambda**  
   - Handles requests from API Gateway.  
   - Written in Go.  
   - Handles user registration, login (JWT creation), and protected route authorization checks.

2. **Amazon DynamoDB**  
   - Stores user information in a table named `userTable`.  
   - Partition key: `username` (string).  

3. **Amazon API Gateway**  
   - Routes traffic to the Lambda function.  
   - Three main paths:
     1. `/register` (POST)
     2. `/login` (POST)
     3. `/protected` (GET)

4. **JWT Security**  
   - The `/protected` route validates a JWT sent in the `Authorization` header.  
   - Tokens expire after **1 hour**.  
   - Token secret is currently hard-coded as `"secret"` (for demonstration only!).

---

## Prerequisites

1. **AWS Account**: You need an AWS account with permissions to create and manage:
   - AWS Lambda
   - DynamoDB
   - API Gateway
   - IAM Roles/Policies

2. **AWS CLI**:  
   - [Install AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) if you haven’t already.  
   - Configure it (`aws configure`) with valid credentials.

3. **CDK CLI**:  
   - [Install AWS CDK](https://docs.aws.amazon.com/cdk/latest/guide/cli.html).
   - Confirm installation with `cdk --version`.

4. **Go (1.18+ recommended)**:  
   - [Download and install Go](https://go.dev/dl/).  
   - Verify with `go version`.

5. **GNU Make**:  
   - Needed to run the `make build` command in the `lambda/` folder.

---

## Setup & Installation

1. **Clone the Repository**  
   ```bash
   git clone https://github.com/your-repo/sibteali786-goCdk.git
   cd sibteali786-goCdk
   ```

2. **Install Dependencies**  
   - At the root:  
     ```bash
     go mod download
     ```
   - In the `lambda/` directory:  
     ```bash
     cd lambda
     go mod download
     cd ..
     ```

---

## Build & Deploy

1. **Build the Lambda**  
   From the `lambda/` directory, run:
   ```bash
   make build
   ```
   This will:
   - Compile a Linux AMD64 Go binary named `bootstrap`.
   - Zip it into `function.zip`.

2. **Synthesize the CDK Stack**  
   Return to the project root, run:
   ```bash
   cdk synth
   ```
   This generates the CloudFormation template in the `cdk.out/` folder.

3. **Deploy to AWS**  
   ```bash
   cdk deploy
   ```
   When prompted, type **`y`** or **`yes`** to allow IAM roles/permissions creation.  

> **Note**: If you have multiple AWS profiles, you may specify `cdk deploy --profile <profile>`.

---

## Testing the Application

After deployment, CDK will output your **API Gateway endpoint**, for example:
```
https://<random>.execute-api.<region>.amazonaws.com/prod/
```
You can test using **curl** or **Postman**:

1. **Register**  
   ```bash
   curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"username":"alice", "password":"mypassword"}' \
     https://<api-url>/register
   ```
   **Responses**:
   - `200 OK` – Successfully Registered
   - `400 Bad Request` – Invalid/empty fields
   - `409 Conflict` – User already exists

2. **Login**  
   ```bash
   curl -X POST \
     -H "Content-Type: application/json" \
     -d '{"username":"alice", "password":"mypassword"}' \
     https://<api-url>/login
   ```
   **Responses**:
   - `200 OK` – Returns JSON like `{"access_token": "<JWT_TOKEN>"}`
   - `400 Bad Request` – Invalid user credentials
   - `500 Internal Server Error` – DynamoDB or other internal error

3. **Protected**  
   ```bash
   curl -X GET \
     -H "Authorization: Bearer <JWT_TOKEN>" \
     https://<api-url>/protected
   ```
   **Responses**:
   - `200 OK` – `"This is protected path"`
   - `401 Unauthorized` – Missing or invalid token, or token expired

---

## API Endpoints

| Method | Path        | Description                                                  |
|--------|------------|--------------------------------------------------------------|
| POST   | `/register` | Register a new user. Body: `{"username": "...", "password": "..."}` |
| POST   | `/login`    | Authenticate an existing user and receive a JWT. Body: `{"username": "...", "password": "..."}` |
| GET    | `/protected`| Access a protected resource, must include `Bearer <JWT>` in `Authorization` header. |

---

## Implementation Details

### Lambda Handlers

- **`main.go`**  
  - The entry point for the Lambda.  
  - Imports `app.NewApp()`, which creates an `ApiHandler` with DB references.  
  - Uses a `switch` statement to route incoming requests based on `request.Path`:  
    - `/register` → `RegisterUserHandler`  
    - `/login` → `LoginUser`  
    - `/protected` → `ValidateJWTMiddleware(ProtectedHandler)`

- **`api/api.go`**  
  - **`RegisterUserHandler`**:  
    1. Parses incoming JSON into `RegisterUser` struct.  
    2. Checks if user exists; if not, hashes password and inserts in DynamoDB.  
    3. Returns success or relevant error.  
  - **`LoginUser`**:  
    1. Parses credentials.  
    2. Verifies password.  
    3. Creates a JWT token and returns to client.

### DynamoDB Database Operations

- **`database/database.go`**  
  - Manages a DynamoDB client via the AWS SDK for Go.  
  - **`DoesUserExist(username string) (bool, error)`**: Returns whether a user item exists.  
  - **`InsertUser(user types.User) error`**: Inserts a new user record (hashed password).  
  - **`GetUser(username string) (types.User, error)`**: Retrieves the user record.

### JWT Middleware

- **`middleware/middleware.go`**  
  - **`ValidateJWTMiddleware`**:  
    - Extracts the token from the `Authorization: Bearer <token>` header.  
    - Parses and validates the token signature and expiry.  
    - If invalid/expired, returns `401 Unauthorized`.  
    - Otherwise, calls the next handler (`ProtectedHandler`).

---

## Local Development Notes

- **Offline Testing**  
  - While AWS SAM or LocalStack can be used, this project is currently set up to run **directly in AWS**.  
  - If you need local development, consider refactoring to use SAM or a local DynamoDB instance.

- **Secrets**  
  - The JWT secret is hard-coded as `"secret"` in `types/types.go`. In production, store secrets in:
    - AWS Secrets Manager, or  
    - AWS SSM Parameter Store  
  - Then retrieve the secret at runtime instead of hard-coding.

- **Logging**  
  - The project uses standard logs in CloudWatch (via the Lambda function logs).  
  - Expand logging as needed to meet your requirements.

---

## Contributing

1. **Fork** the repository.
2. **Create** a new feature branch.
3. **Commit** your changes.
4. **Push** to your branch.
5. Submit a **Pull Request**.

All contributions are welcome!

---

## License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT). See the [LICENSE](LICENSE) file for details.

---

**Enjoy building and deploying your serverless authentication API with Go and AWS CDK!**
