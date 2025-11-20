# 🍺 Beer Microservices

<p align="center">
  Beer management system based on microservices, developed with Go + Gin and MongoDB,<br>
  orchestrated with Docker and Docker Compose.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25.x-00ADD8?logo=go&logoColor=white" alt="Go Badge">
  <img src="https://img.shields.io/badge/MongoDB-Database-47A248?logo=mongodb&logoColor=white" alt="MongoDB Badge">
  <img src="https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white" alt="Docker Badge">
  <img src="https://img.shields.io/badge/Architecture-Microservices-blueviolet" alt="Microservices Badge">
</p>

---

## 📋 Project Description

**Beer API** is a microservices system designed to manage a beer catalog. Each microservice handles a specific operation (create, read, update, delete) and communicates with a centralized **MongoDB** database.

- 🔧 Backend implemented with **Go** using the **Gin** framework
- 💾 Data persistence in **MongoDB**
- 🐳 Containers managed using **Docker** and **Docker Compose**
- 📦 Microservices-based architecture

---

## ✅ Features

- **create-service** microservice (Port 8080)
- **read-service** microservice (Port 8081)
- **update-service** microservice (Port 8082)
- **delete-service** microservice (Port 8083)
- `docker-compose.yml` configured to orchestrate services and database
- Tests implemented for all services
- CI/CD workflows with GitHub Actions
- Standardized architecture with adapters, repositories, services, and controllers

---

## 🏗️ Infrastructure Diagram

<p align="center">
  <img src="./assets/diagram.svg" width="95%" alt="Beer API Architecture">
</p>

---

## ⚙️ Prerequisites

- 🐳 **Docker** and **Docker Compose** installed
- 🔧 **Go 1.25.x** (optional for local development)

---

## 🔐 Environment Variables

`.env` file (located at the same level as `docker-compose.yml`):

```env
MONGO_URI=mongodb://beers-mongodb:27017
DATABASE=beersdb
COLLECTION=beers
```

These variables are used by the microservices to connect to the MongoDB database.

> [!NOTE]
> The hostname must match the MongoDB container name defined in `docker-compose.yml` (`beers-mongodb`).

---

## 🚀 How to Run the Project

1. **Clone the repository from GitHub**

   ```bash
   git clone https://github.com/N3VERS4YDIE/beer-microservices.git
   ```

2. **Create the `.env` file** in the project root with the following content:

   ```env
   MONGO_URI=mongodb://beers-mongodb:27017
   DATABASE=beersdb
   COLLECTION=beers
   ```

3. **Start containers with Docker Compose**

   ```bash
   docker compose up -d
   ```

   This will start MongoDB and the microservices.

4. **Verify running containers**

   ```bash
   docker compose ps
   ```

   You should see five active containers:
   - `beers-mongodb` - MongoDB database (Port 27017)
   - `beers-create-service` - Creation service (Port 8080)
   - `beers-read-service` - Query service (Port 8081)
   - `beers-update-service` - Update service (Port 8082)
   - `beers-delete-service` - Deletion service (Port 8083)

5. **Check microservice logs**

   ```bash
   docker logs beers-create-service
   docker logs beers-read-service
   docker logs beers-update-service
   docker logs beers-delete-service
   ```

---

## 🔌 Endpoints

### Create Service (Port 8080)

- `POST /beers` – Create a new beer

  ```bash
  curl -X POST http://localhost:8080/beers \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Heineken",
      "brand": "Heineken",
      "alcohol": 5.0,
      "year": 2024
    }'
  ```

### Read Service (Port 8081)

- `GET /beers` – Get all beers

  ```bash
  curl http://localhost:8081/beers
  ```

- `GET /beers/:id` – Get a specific beer by ID

  ```bash
  curl http://localhost:8081/beers/BEER_ID
  ```

### Update Service (Port 8082)

- `PUT /beers/:id` – Update an existing beer

  ```bash
  curl -X PUT http://localhost:8082/beers/BEER_ID \
    -H "Content-Type: application/json" \
    -d '{
      "name": "Heineken Premium",
      "brand": "Heineken International",
      "alcohol": 5.2,
      "year": 2025
    }'
  ```

### Delete Service (Port 8083)

- `DELETE /beers/:id` – Delete a beer by ID

  ```bash
  curl -X DELETE http://localhost:8083/beers/BEER_ID
  ```

### Beer Data Model

```json
{
  "name": "string",      // Beer name
  "brand": "string",     // Beer brand
  "alcohol": "float",    // Alcohol content
  "year": "int"          // Year
}
```

---

## 🔀 Git Workflow

1. **Create working branch from `develop`:**

   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/new-functionality
   ```

2. **Make changes, commit, and push:**

   ```bash
   git add .
   git commit -m "Clear description of changes"
   git push origin feature/new-functionality
   ```

3. **Create Pull Request to `develop`**.

---

## 🧪 Run Tests

**Create Service:**

```bash
cd create-service && go test ./... -v
```

**Read Service:**

```bash
cd read-service && go test ./... -v
```

**Update Service:**

```bash
cd update-service && go test ./... -v
```

**Delete Service:**

```bash
cd delete-service && go test ./... -v
```
