# Requirements
This file keeps the Problem statement, features, and tasks of the development project. It is updated on a revolving basis.
Features and tasks have IDs prefixed with F and T, respectively, and are arranged as checklist items. A checked box indicates that the feature or task is done. A feature is checked only of all tasks under it are checked. The following example shows a snapshot:

## Example of Feature checklist
    - [ ] **F1:** Description of feature 1 
        - [x] **T1.1:** the normal case: inner product computed for two vectors of legal and equal dimensions
        - [ ] **T1.2:** exception case: when a vector has an illegal dimension
    - [x] **F2:** Description of feature 2
        - [x] **T2.1:** Create a `Vector` class that can be constructed and stores vector data.

## Problem statement
A note app is needed to let users create, organize, and share notes.
A note begins with a title, followed by a sequence of contents that are texts, pictures, and so on.
Contents can be deleted.
A text content can be edited.
A text content ends with a new line '\n'.
Pictures cannot be edited.
A note should have an unique ID. 
A note can be tagged with any number of keywords. 
A note can be found through keyword filtering. 
A note can be shared and co-edited among multiple users. Two users opening a shared note should see the same contents.
A user should see all the notes they own.
An owner can read, write, share, and delete a note.
A user should see all the notes they were shared with.
An owner can share a note as read or read/write to another user.
A user cannot delete a note not owned by them.
A user cannot share a note not owned by them.
A user should see which content another user is on in a shared note.   
Two users can change two different contents in a shared note simultaneously.
Only one of multiple users can change the content if they are on the same content in a shared note. 
Keywords are user-specific and are private to the users. A shared note can have a different set of keywords tagged to it among its users.
A user can have multiple devices. These devices can access all the notes and keywords owned by the user. The notes and keywords should be synchronized among the devices owned by a user.
User facing APIs should guard against illegal parameters

## Features
- [ ] **F1:** Note Lifecycle Management. Users can create, read, update, and delete their own notes.
    - [x] **T1.1:** Create a `Note` model with attributes like ID, title, and content.
    - [x] **T1.2:** Create a `NoteUsecase` with a `CreateNote` method that handles the business logic for creating a new note.
    - [x] **T1.3:** Implement the `POST /notes` API endpoint in the `api` layer, which will call the `CreateNote` usecase.
    - [x] **T1.4:** Update `main.go` to initialize all dependencies and start the HTTP server.
    - [x] **T1.5:** Add a `GetNoteByID` method to the `NoteUsecase` that returns a `NoteDTO`, and define the `NoteDTO` in the `usecase` package.
    - [x] **T1.6:** Implement the `GET /notes/{id}` API endpoint, which will receive the `NoteDTO`.
    - [ ] **T1.7:** Add an `UpdateNote` method to the `NoteUsecase`.
    - [ ] **T1.8:** Implement the `PUT /notes/{id}` API endpoint.
    - [x] **T1.9:** Add a `DeleteNote` method to the `NoteUsecase`.
    - [x] **T1.10:** Implement the `DELETE /notes/{id}` API endpoint.
    - [x] **T1.11:** Refactor `NoteUsecase` to translate repository-specific errors into use case-level errors.
- [ ] **F2:** Note Content Management. A note is composed of text and pictures. Users can add, edit, and delete note contents.
    - [x] **T2.1:** Redefine the `Note` and `Content` models in the `domain` layer.
    - [x] **T2.2:** Implement `AddContent` and `Contents` methods on the `Note` entity in the domain layer.
    - [x] **T2.3:** Create an `AddContent` method in `NoteUsecase`.
    - [x] **T2.4:** Update `NoteDTO` to handle a slice of contents.
    - [x] **T2.5:** Implement the `POST /notes/{id}/contents` API endpoint.
    - [ ] **T2.6:** Implement an `UpdateContent` method on the `Note` entity in the domain layer.
    - [ ] **T2.7:** Create an `UpdateContent` method in `NoteUsecase`.
    - [ ] **T2.8:** Implement the `PUT /notes/{id}/contents/{contentId}` API endpoint.
    - [ ] **T2.9:** Implement a `DeleteContent` method on the `Note` entity in the domain layer.
    - [ ] **T2.10:** Create a `DeleteContent` method in `NoteUsecase`.
    - [ ] **T2.11:** Implement the `DELETE /notes/{id}/contents/{contentId}` API endpoint.
    - [ ] **T2.12:** Implement logic to handle `ImageContentType`.
- [ ] **F3:** Note Tagging and Searching. Users can tag notes with keywords and search for notes using these keywords.
- [ ] **F4:** Note Sharing and Collaboration. Users can share notes with others, enabling co-editing and real-time content visibility.
- [ ] **F5:** Multi-Device Synchronization. User's notes and keywords are synchronized across all their devices.
- [ ] **F6:** API Security. APIs validate input to prevent errors and misuse.
- [ ] **F7:** Decouple Data Persistence with a Repository Layer.
    - [x] **T7.1:** Define a `NoteRepository` interface with methods for note persistence (e.g., `Save`, `GetByID`).
    - [x] **T7.2:** Create an `InMemoryNoteRepository` implementation that satisfies the `NoteRepository` interface.
    - [ ] **T7.3:** Make the `InMemoryNoteRepository` thread-safe.
- [ ] **F8:** API and Codebase Polish.
    - [ ] **T8.1:** Refactor: Standardize API error responses to return JSON objects.
    - [ ] **T8.2:** Refactor: Move router setup out of `main.go` to improve modularity.
    - [ ] **T8.3:** Refactor: Centralize API error handling in a helper function.
- [ ] **F9:** Decouple Domain and Persistence Layers.
    - [x] **T9.1:** Create `NotePO` and `ContentPO` in the `repository` layer, and implement a `NoteMapper` in the `usecase` layer to map between `domain.Note` and `repository.NotePO`.
    - [ ] **T9.2:** Update the `NoteRepository` interface, `InMemoryNoteRepository`, and `NoteUsecase` to use the new `NotePO` and `NoteMapper`.
