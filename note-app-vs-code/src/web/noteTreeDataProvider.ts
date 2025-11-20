import * as vscode from 'vscode';
import { Note, NoteService } from './noteService';
var userId = 'testUser1';

export class NoteTreeDataProvider implements vscode.TreeDataProvider<NoteTreeItem> {
  private _onDidChangeTreeData: vscode.EventEmitter<NoteTreeItem | undefined | void> = new vscode.EventEmitter<NoteTreeItem | undefined | void>();
  readonly onDidChangeTreeData: vscode.Event<NoteTreeItem | undefined | void> = this._onDidChangeTreeData.event;

  constructor(private noteService: NoteService) {}

  getTreeItem(element: NoteTreeItem): vscode.TreeItem {
    return element;
  }

  async getChildren(element?: NoteTreeItem): Promise<NoteTreeItem[]> {
    if (element) {
      return Promise.resolve([]); // No children for notes yet
    } else {
      // Fetch notes for a dummy user ID for now
      const notes = await this.noteService.getNotes(userId);
      // const notes: Note[] = [
      //   { id: '1', title: 'Note 1', content_ids: [], version: 1 },
      //   { id: '2', title: 'Note 2', content_ids: [], version: 1 },
      // ]; 
      console.log('Fetched notes:', notes);
      return notes.map(note => new NoteTreeItem(note.id, note.title, vscode.TreeItemCollapsibleState.None));
    }
  }

  refresh(): void {
    this._onDidChangeTreeData.fire();
  }
}

export class NoteTreeItem extends vscode.TreeItem {
  constructor(
    public readonly noteId: string,
    public readonly noteTitle: string,
    public readonly collapsibleState: vscode.TreeItemCollapsibleState
  ) {
    super(noteTitle, collapsibleState);
    this.tooltip = this.noteTitle;
    this.description = this.noteTitle;
  }
}
