import * as vscode from 'vscode';

export class NoteTreeDataProvider implements vscode.TreeDataProvider<NoteTreeItem> {
  constructor() {}

  getTreeItem(element: NoteTreeItem): vscode.TreeItem {
    return element;
  }

  getChildren(element?: NoteTreeItem): Thenable<NoteTreeItem[]> {
    // For now, return an empty array. We will implement fetching notes later.
    return Promise.resolve([]);
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
