import * as vscode from 'vscode';
import { Note, Content, NoteService } from './noteService';

function getWebviewScript(): string {
  return `
    const vscode = acquireVsCodeApi();
    window.addEventListener('message', event => {
        const message = event.data;
        switch (message.type) {
            case 'loadNote':
                const note = message.payload;
                document.getElementById('note-title').textContent = note.title;
                const contentsDiv = document.getElementById('note-contents');
                contentsDiv.innerHTML = ''; // Clear existing content
                if (note.contents && note.contents.length > 0) {
                    note.contents.forEach(content => {
                        const p = document.createElement('p');
                        p.textContent = content.data;
                        contentsDiv.appendChild(p);
                    });
                } else {
                    const p = document.createElement('p');
                    p.textContent = 'This note is empty.';
                    contentsDiv.appendChild(p);
                }
                break;
        }
    });
  `;
}

function getHtmlForWebview(webview: vscode.Webview): string {
  const script = getWebviewScript();
  return `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Note Editor</title>
    </head>
    <body>
        <h1 id="note-title">Loading Note...</h1>
        <div id="note-contents"></div>
        <script>${script}</script>
    </body>
    </html>`;
}

export function openNoteEditor(noteId: string, noteTitle: string) {
	const panel = vscode.window.createWebviewPanel(
		'note-app-vs-code.noteEditor',
		noteTitle,
		vscode.ViewColumn.One,
		{
			enableScripts: true
		}
	);

	panel.webview.html = getHtmlForWebview(panel.webview);

	// Fetch note data and send it to the webview
	NoteService.getInstance().getNoteById(noteId)
		.then(note => {
			if (note) {
				panel.webview.postMessage({ type: 'loadNote', payload: note });
			}
            console.log('Note loaded:', note);
		})
		.catch(error => {
			vscode.window.showErrorMessage(`Failed to load note: ${error.message}`);
		});
}
