import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Content {
  id: string;
  noteID: string;
  data: string;
  type: string;
  version: number;
}

export interface Note {
  id: string;
  title: string;
  contents: Content[];
  version: number;
  collaborators: { [key: string]: string }; // userId: permission
  keywords: string[];
}

@Injectable({
  providedIn: 'root',
})
export class NoteService {
  private usersApiUrl = 'http://localhost:8080/users'; // Assuming backend runs on 8080
  private notesApiUrl = 'http://localhost:8080/notes'; // Assuming backend runs on 8080

  constructor(private http: HttpClient) {}

  getAccessibleNotes(userId: string): Observable<Note[]> {
    return this.http.get<Note[]>(`${this.usersApiUrl}/${userId}/accessible-notes`);
  }

  getNoteById(noteId: string): Observable<Note> {
    return this.http.get<Note>(`${this.notesApiUrl}/${noteId}`);
  }
}
