// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from 'vscode';
import { NoteTreeDataProvider } from './noteTreeDataProvider';
import { NoteService } from './noteService';
import { openNoteEditor } from './noteEditor';
import { AuthService } from './authService';

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed
export function activate(context: vscode.ExtensionContext) {

	// Use the console to output diagnostic information (console.log) and errors (console.error)
	// This line of code will only be executed once when your extension is activated
	console.log('Congratulations, your extension "note-app-vs-code" is now active in the web extension host!');

	// Initialize services
	AuthService.initialize(context);
	const authService = AuthService.getInstance();
	const noteService = NoteService.getInstance();

	const noteTreeDataProvider = new NoteTreeDataProvider(noteService);
	vscode.window.createTreeView('noteapp.notes', { treeDataProvider: noteTreeDataProvider });

	// Automatic login for test user
	authService.login("testuser", "password").then(() => {
		console.log("Test user logged in automatically.");
		noteTreeDataProvider.refresh();
	}).catch(error => {
		console.error("Automatic login failed:", error);
		vscode.window.showErrorMessage('Automatic login failed. Please check backend server.');
	});


	// The command has been defined in the package.json file
	// Now provide the implementation of the command with registerCommand
	// The commandId parameter must match the command field in package.json
	const disposable = vscode.commands.registerCommand('note-app-vs-code.helloWorld', async () => {
		// The code you place here will be executed every time your command is executed

		// Display a message box to the user
		// const notes = await noteTreeDataProvider.getChildren();
		// const noteTitles = notes.map(note => note.noteTitle).join(', ');
		// vscode.window.showInformationMessage(`Fetched Notes: ${noteTitles}`);
	});

	const openNoteCommand = vscode.commands.registerCommand('note-app-vs-code.openNote', (noteId: string, noteTitle: string) => {
		openNoteEditor(noteId, noteTitle);
	});

	context.subscriptions.push(disposable, openNoteCommand);
}

// This method is called when your extension is deactivated
export function deactivate() {}
