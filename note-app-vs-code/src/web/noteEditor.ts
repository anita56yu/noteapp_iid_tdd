import * as vscode from 'vscode';

function getHtmlForWebview(webview: vscode.Webview): string {
  return `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Note Editor</title>
    </head>
    <body>
        <h1>Note Editor</h1>
        <p>This is a placeholder for the note editor UI.</p>
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
}
