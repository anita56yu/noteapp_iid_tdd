import * as vscode from 'vscode';

export interface Note {
  id: string;
  title: string;
}

export class NoteService {
  constructor() {}

  async getNotesForUser(userId: string): Promise<Note[]> {
    // For now, return mock data
    return [
      { id: '1', title: 'My First Note' },
      { id: '2', title: 'Another Note' },
      { id: '3', title: 'A Third Note' },
    ];
  }
}
