import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Content {
  id: string;
  noteId: string;
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
    console.log('Updating content ID:', content.id, 'note ID:', content.noteId, 'with new text:', data);
    return this.http.put<void>(`${this.notesApiUrl}/${content.noteId}/contents/${content.id}`, { data, content_version: version });
  }

  deleteContent(noteId: string, contentId: string, noteVersion: number, contentVersion: number): Observable<void> {
    return this.http.request<void>('DELETE', `${this.notesApiUrl}/${noteId}/contents/${contentId}`, { body: { note_version: noteVersion, content_version: contentVersion } });
  }

  createNote(userId: string): Observable<{ id: string }> {
    const defaultTitle = 'New Note';
    return this.http.post<{ id: string }>(this.notesApiUrl, { title: defaultTitle, owner_id: userId });
  }

  updateNote(noteId: string, title: string, noteVersion: number): Observable<void> {
    return this.http.put<void>(`${this.notesApiUrl}/${noteId}`, { title, note_version: noteVersion });
  }
}
