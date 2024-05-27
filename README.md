# Chirpy

Chirpy is a simple RESTful API built with Go that allows users to create, read, update, and delete "chirps" (short messages) and manage their user accounts.

## Features

- User authentication and authorization with JSON Web Tokens (JWT)
- Create, read, update, and delete chirps
- Filter chirps by author
- Sort chirps by ID in ascending or descending order
- Create and manage user accounts
- Upgrade users to "Chirpy Red" membership
- Webhook integration for handling user upgrades from payment providers

## Getting Started

### Prerequisites

- Go (version 1.16 or later)
- An environment with the following environment variables set:
  - `JWT_SECRET`: A secret key used for signing and verifying JSON Web Tokens
  - `POLKA_API_KEY`: An API key provided by the Polka payment provider for handling webhooks

### Installation

1. Clone the repository:

```
git clone https://github.com/your-username/chirpy.git
```

2. Navigate to the project directory:

```
cd chirpy
```

3. Build and run the server:

```
go build -o out && ./out
```

To clear the database before starting the server, you can use the `--debug` flag

The server should now be running on `http://localhost:8080`.

### Usage

The Chirpy API provides the following endpoints:

- `POST /api/users`: Create a new user account
- `POST /api/login`: Authenticate a user and obtain a JWT
- `PUT /api/users`: Update a user's email or password
- `POST /api/chirps`: Create a new chirp
- `GET /api/chirps`: Retrieve all chirps or filter by author
- `GET /api/chirps/{chirpID}`: Retrieve a single chirp by ID
- `DELETE /api/chirps/{chirpID}`: Delete a chirp (requires authentication)
- `POST /api/polka/webhooks`: Handle webhooks from the Polka payment provider for user upgrades
