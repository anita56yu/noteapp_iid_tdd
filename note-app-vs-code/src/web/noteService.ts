import axios from 'axios';

export interface Note {
  id: string;
  title: string;
  content_ids: string[];
  version: number;
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

  async getNotes(userId: string): Promise<Note[]> {
    try {
      const response = await axios.get<Note[]>(`${this.baseUrl}/users/${userId}/accessible-notes`);
      console.log('Fetched notes:', response);
      return response.data;
    } catch (error) {
      console.error('Error fetching notes:', error);
      return [];
    }
  }
}