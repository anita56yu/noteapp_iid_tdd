import axios from 'axios';
import { AuthService } from './authService';

export interface Content {
  id: string;
  note_id: string;
  data: string;
  content_type: string;
  version: number;
}

export interface Note {
  id: string;
  title: string;
  content_ids: string[];
  version: number;
  contents?: Content[];
}

export class NoteService {
  private static instance: NoteService;
  private baseUrl = 'http://localhost:8080'; // Adjust if your backend URL is different

  private constructor() {}

  public static getInstance(): NoteService {
    if (!NoteService.instance) {
      NoteService.instance = new NoteService();
    }
    return NoteService.instance;
  }

  private async getAuthHeaders(): Promise<{ Authorization?: string }> {
    const token = await AuthService.getInstance().getToken();
    return token ? { Authorization: `Bearer ${token}` } : {};
  }

  async getNotes(): Promise<Note[]> {
    try {
      const token = await AuthService.getInstance().getToken();
      if (!token) {
        console.log('No authentication token found. Returning empty notes array.');
        return [];
      }
      const headers = await this.getAuthHeaders();
      const response = await axios.get<Note[]>(`${this.baseUrl}/notes/accessible-notes`, { headers });
      console.log('Fetched notes:', response);
      return response.data;
    } catch (error) {
      console.error('Error fetching notes:', error);
      return [];
    }
  }

  async getNoteById(noteId: string): Promise<Note | undefined> {
    try {
      const headers = await this.getAuthHeaders();
      const response = await axios.get<Note>(`${this.baseUrl}/notes/${noteId}`, { headers });
      return response.data;
    } catch (error) {
      console.error(`Error fetching note ${noteId}:`, error);
      return undefined;
    }
  }
}