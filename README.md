# Project: ToDo App (Go + PostgreSQL + Redis)

## Overview:
   Basic ToDo application implemented in Go, using PostgreSQL  for persistent storage and Redis for caching.
   Demonstrates clean architecture, design patterns, and Docker-based deployment.
   
## Tech Stack:
   - Language: Go (Golang)
   - Database: PostgreSQL (via pgxpool)
   - Cache: Redis
   - Authentication: JWT (HMAC-SHA algorithm)
   - Containerization: Docker

## Architecture & Design Principles:
  - Built on Clean Code principles
  - Decorator Pattern for cross-cutting functionality
  - Clear separation of layers:
      * Handler (HTTP/API layer)
      * Service (business logic)
      * Repository/DAO (database layer)
     → No direct interaction between handler and repo
   - Enums in code for mapping values like status
   - Redis caching for user details (demo integration)
   - JWT-based authentication (HSA algorithm)

## Docker:
   Fully dockerized project.
   Build and run using:

    docker build -t todo-app .
    docker run -p 8080:8080 todo-app
 Key Features:
   - User registration & authentication (JWT)
   - Create / Update / Delete / List ToDo items
   - Enum-based status tracking
   - Redis cache integration for users
   - PostgreSQL with pgxpool connection pooling
   - Clean modular structure
     
 Future Improvements:
   - Add unit and integration tests
   - Expand caching strategy for ToDo lists

(refer [API documentation](https://documenter.getpostman.com/view/32823353/2sB3WnxNFC))
