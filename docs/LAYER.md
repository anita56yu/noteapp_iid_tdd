# Clean Architecture Layers

This project adheres to the principles of Clean Architecture (CA). The layers are organized to enforce a one-way dependency flow, pointing inwards towards the domain. At no time should a domain layer object show up at the api layer.

## 1. Domain Layer
- **Location:** `backend/internal/domain`
- **Description:** This is the core of the application. It contains the enterprise-wide business logic and entities (e.g., `Note`, `Content`). This layer is the most stable and has zero dependencies on any other layer in the project.

## 2. Usecase Layer
- **Location:** `backend/internal/usecase`
- **Description:** This layer contains the application-specific business rules. It orchestrates the flow of data to and from the domain entities to achieve the goals of a particular use case (e.g., `CreateNote`). It depends only on the Domain layer. It also defines the interfaces for repositories that the outer layers must implement.

## 3. API Layer (Interface Adapters)
- **Location:** `backend/internal/api`
- **Description:** This layer is responsible for handling all external interactions, primarily the HTTP requests from the frontend. It acts as a translator, converting external data formats (like JSON) into formats that the Usecase layer can understand, and vice versa. It depends on the Usecase layer.

## 4. Repository Layer (Frameworks & Drivers)
- **Location:** `backend/internal/repository`
- **Description:** This layer provides the concrete implementations for the data persistence interfaces defined in the Usecase layer (e.g., `InMemoryNoteRepository`). It handles all communication with the database or other data sources.

## The Dependency Rule

The core principle of our architecture is that **dependencies must only point inwards**.

The dependency chain is as follows:

`API Layer` -> `Usecase Layer` -> `Domain Layer`

The Repository layer adheres to this rule through **Dependency Inversion**. The Usecase layer defines a `NoteRepository` *interface*, and the concrete `InMemoryNoteRepository` *implements* it. This means the Usecase layer does not depend on the concrete repository, maintaining the inward-pointing dependency flow.
