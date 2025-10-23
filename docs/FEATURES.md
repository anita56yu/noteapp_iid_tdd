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

## Design Decisions
### Concurrency of Notes and Contents
    **Solution #1:** For concurrency of notes and contents, we will use optimistic locking on the repository side. On the client side(frontend), it will keep 2 copies of the note it is working on: the real copy and the working copy. The real copy is the closest to its counter part in the repository, and the working copy as the current content that is changed. Everytime the repository is updated, the client side receives an event, and updates the real copy accordingly. Simultaneously, the client side also need to merge the update into the working copy. If the update doesn't confict with the working content, it should be successfully merged, and update the version number. If the update conflict with the working content, the client side discard or manually merged the working content(stale.)
    
    **Solution #2:** Treats content and note as respective aggregates which are built association with ids(notes will hold a slice of content ids, and content will hold a note id.) Each content and note will hold its version number and use optimistic lock for consistency. A content change will only trigger a content event; a deletion or insertion of content will trigger a note event. The client side will only keep a working copy and apply events received. In case of a conflict, which will be reflected by version number, the local change is stashed. 

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
    - [x] **T2.6:** Implement an `UpdateContent` method on the `Note` entity in the domain layer.
    - [x] **T2.7:** Create an `UpdateContent` method in `NoteUsecase`.
    - [x] **T2.8:** Implement the `PUT /notes/{id}/contents/{contentId}` API endpoint.
    - [x] **T2.9:** Implement a `DeleteContent` method on the `Note` entity in the domain layer.
    - [x] **T2.10:** Create a `DeleteContent` method in `NoteUsecase`.
    - [x] **T2.11:** Implement the `DELETE /notes/{id}/contents/{contentId}` API endpoint.
    - [ ] **T2.12:** Implement logic to handle `ImageContentType`.
- [x] **F3:** Note Tagging and Searching. Users can tag notes with keywords and search for notes using these keywords.
    - [x] **T3.1:** Define `Keyword` as a value object in the `domain` layer.
    - [x] **T3.2:** Enhance the `Note` domain model to store user-specific keywords and add an `AddKeyword` method.
    - [x] **T3.3:** Update `NotePO` in the repository and the `NoteMapper` in the usecase to handle the new tag data.
    - [x] **T3.4:** Implement the `TagNote` method in `NoteUsecase` to fetch, update, and save the note.
    - [x] **T3.5:** Implement the `POST /users/{userID}/notes/{noteID}/keywords` API endpoint.
    - [x] **T3.6:** Add a `FindByKeywordForUser` method to the `NoteRepository` interface and `InMemoryNoteRepository`.
    - [x] **T3.7:** Implement the `FindNotesByKeyword` method in `NoteUsecase` to filter notes in memory.
    - [x] **T3.8:** Implement the `GET /users/{userID}/notes?keyword={keyword}` API endpoint.
    - [x] **T3.9:** Implement an `UntagNote` method in `NoteUsecase`.
    - [x] **T3.10:** Implement a `DELETE /users/{userID}/notes/{noteID}/keywords/{keyword}` API endpoint.
- [x] **F4:** Note Sharing and Access. An owner can share a note with other users, specifying their permissions (read-only or read-write), and users can view and access notes that have been shared with them.
    - [x] **T4.1:** In the `domain` layer, update the `Note` entity to include a list of collaborators and their permissions (e.g., read, read-write).
    - [x] **T4.2:** Create a `ShareNote` method in the `NoteUsecase` that allows a note owner to share a note with another user and set their permissions.
    - [x] **T4.3:** Implement a `POST /users/{ownerID}/notes/{noteID}/shares` API endpoint to expose the `ShareNote` functionality. This endpoint will take a user ID and permission level in the request body.
    - [x] **T4.4:** Add a `GetAccessibleNotesByUserID` method to the `NoteRepository` to retrieve all notes shared with or owned by a specific user.
    - [x] **T4.5:** Create a `GetAccessibleNotesForUser` method in `NoteUsecase`.
    - [x] **T4.6:** Implement a `GET /users/{userID}/accessible-notes` API endpoint to allow users to see all the accessible notes.
    - [x] **T4.7:** Add a `RemoveCollaborator` method to the `Note` entity in the domain layer.
    - [x] **T4.8:** Create a `RevokeAccess` method in the `NoteUsecase`.
    - [x] **T4.9:** Implement a `DELETE /users/{ownerID}/notes/{noteID}/shares` API endpoint.
- [ ] **F5:** Real-time Collaboration and Concurrent Editing. Users can see who is currently editing a content block and view changes made by others in real-time. The system will manage simultaneous edits to prevent conflicts while allowing users to work on different parts of a note at the same time. Use **Solution #2** in the Design Decisions section above.
    - [x] **T5.1:** In the `domain` layer, make the `Content` entity a standalone aggregate by adding `ID`, `NoteID`, and `Version` fields.
    - [x] **T5.2:** In the `domain` layer, update the `Note` entity to hold a slice of content IDs and a `Version` field.
    - [x] **T5.3:** In the `repository` layer, define a `ContentRepository` interface and create a `ContentPO`.
    - [x] **T5.4:** In the `repository` layer, create an `InMemoryContentRepository` that implements the `ContentRepository` interface, including optimistic locking.
    - [ ] **T5.5:** In the `repository` layer, update `NotePO` and `InMemoryNoteRepository` to support optimistic locking using the `Version` field.
    - [ ] **T5.6:** In the `usecase` layer, create a `ContentDTO`, a `ContentMapper`, and a `ContentUsecase` with `CreateContent`, `UpdateContent`, and `DeleteContent` methods.
    - [ ] **T5.7:** In `NoteUsecase`, refactor the `AddContent` method to accept a `contentID` and add it to the note's list.
    - [ ] **T5.8:** In `NoteUsecase`, refactor the `DeleteContent` method to accept a `contentID` and remove it from the note's list.
    - [ ] **T5.9:** Update `NoteMapper` to only map content IDs between `NotePO` and `domain.Note`.
    - [ ] **T5.10:** Create a `content_handler` for `Content` CRUD operations and update `note_handler` to orchestrate the relationship between notes and content.
- [ ] **F6:** Multi-Device Synchronization. User's notes and keywords are synchronized across all their devices.
- [ ] **F7:** API Security. APIs validate input to prevent errors and misuse.
- [ ] **F8:** Decouple Data Persistence with a Repository Layer.
    - [x] **T8.1:** Define a `NoteRepository` interface with methods for note persistence (e.g., `Save`, `GetByID`).
    - [x] **T8.2:** Create an `InMemoryNoteRepository` implementation that satisfies the `NoteRepository` interface.
    - [x] **T8.3:** Make the `InMemoryNoteRepository` thread-safe, and anti-racing. When two users are changing the same note, the PO in the repository should be anti-racing.
    - [ ] **T8.4:** (Potential) Refactor the repository to use a transactional callback pattern for thread-safe updates, moving locking logic out of the usecase layer.
- [ ] **F9:** API and Codebase Polish.
    - [ ] **T9.1:** Refactor: Standardize API error responses to return JSON objects.
    - [ ] **T9.2:** Refactor: Move router setup out of `main.go` to improve modularity.
    - [ ] **T9.3:** Refactor: Centralize API error handling in a helper function.
- [x] **F10:** Decouple Domain and Persistence Layers.
    - [x] **T10.1:** Create `NotePO` and `ContentPO` in the `repository` layer, and implement a `NoteMapper` in the `usecase` layer to map between `domain.Note` and `repository.NotePO`.
    - [x] **T10.2:** Update the `NoteRepository` interface, `InMemoryNoteRepository`, and `NoteUsecase` to use the new `NotePO` and `NoteMapper`.
