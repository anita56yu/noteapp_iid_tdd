import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Content {
  id: string;
  noteID: string;
  data: string;
  type: string;
  version: number;
  position: number;
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

  addContent(noteId: string, content: Content, noteVersion: number): Observable<{ id: string }> {
    const { data, type, position } = content;
    return this.http.post<{ id: string }>(`${this.notesApiUrl}/${noteId}/contents`, { data, type, index: position, note_version: noteVersion });
  }

  updateContent(content: Content): Observable<void> {
    const { data, version } = content;
    return this.http.put<void>(`${this.notesApiUrl}/${content.noteID}/contents/${content.id}`, { data, content_version: version });
  }
}
