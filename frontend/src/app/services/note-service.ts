import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Note {
  id: string;
  title: string;
  contentIDs: string[];
  version: number;
  collaborators: { [key: string]: string }; // userId: permission
  keywords: string[];
}

@Injectable({
  providedIn: 'root',
})
export class NoteService {
  private apiUrl = 'http://localhost:8080/users'; // Assuming backend runs on 8080

  constructor(private http: HttpClient) {}

  getAccessibleNotes(userId: string): Observable<Note[]> {
    return this.http.get<Note[]>(`${this.apiUrl}/${userId}/accessible-notes`);
  }
}
