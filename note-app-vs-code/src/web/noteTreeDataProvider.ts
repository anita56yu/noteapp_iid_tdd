import * as vscode from 'vscode';
import { Note, NoteService } from './noteService';

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
      const notes = await this.noteService.getNotes();
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
    this.command = {
      command: 'note-app-vs-code.openNote',
      title: 'Open Note',
      arguments: [this.noteId, this.noteTitle],
    };
  }
}
