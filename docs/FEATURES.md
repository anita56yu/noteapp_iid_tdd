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
    - [ ] **T1.2:** Implement an API endpoint to create a new note.
    - [ ] **T1.3:** Implement an API endpoint to retrieve a note by its ID.
    - [ ] **T1.4:** Implement an API endpoint to update an existing note.
    - [ ] **T1.5:** Implement an API endpoint to delete a note by its ID.
    - [ ] **T1.6:** Write unit tests for the Note model and API endpoints.
- [ ] **F2:** Note Content Management. A note is composed of text and pictures. Users can add, edit, and delete note contents.
- [ ] **F3:** Note Tagging and Searching. Users can tag notes with keywords and search for notes using these keywords.
- [ ] **F4:** Note Sharing and Collaboration. Users can share notes with others, enabling co-editing and real-time content visibility.
- [ ] **F5:** Multi-Device Synchronization. User's notes and keywords are synchronized across all their devices.
- [ ] **F6:** API Security. APIs validate input to prevent errors and misuse.
