# golang-db

This is a Go project with a database setup that demonstrates basic CRUD operations.

## Prerequisites

- Go 1.16 or later
- MySQL database
- `go-sql-driver/mysql` package
- `joho/godotenv` package
- Docker
- Docker Compose

## Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/golang-db.git
    cd golang-db
    ```

2. Create a `.env` file in the project root with the following content:
    ```env
    DBUSER=your_db_user
    DBPASS=your_db_password
    ```

3. Create the database and tables:


4. Install dependencies:
    ```sh
    go mod tidy
    ```

## Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

2. The server will start on `http://localhost:8080`.

## API Endpoints

### Get Albums by Artist

- **URL:** `/albums`
- **Method:** `GET`
- **Query Parameter:** `artist` (e.g., `Pink Floyd`)
- **Example Request:**
    ```
    http://localhost:8080/albums?artist=Pink%20Floyd
    ```

### Get Album by ID

- **URL:** `/album`
- **Method:** `GET`
- **Query Parameter:** `id` (e.g., `1`)
- **Example Request:**
    ```
    http://localhost:8080/album?id=1
    ```

### Add New Album

- **URL:** `/addAlbum`
- **Method:** `POST`
- **Request Body:**
    ```json
    {
        "title": "Album Title",
        "artist": "Artist Name",
        "price": 19.99
    }
    ```

### Delete Album by ID

- **URL:** `/deleteById`
- **Method:** `DELETE`
- **Query Parameter:** `id` (e.g., `1`)
- **Example Request:**
    ```
    http://localhost:8080/deleteById?id=1
    ```
