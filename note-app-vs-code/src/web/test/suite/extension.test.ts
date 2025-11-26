import * as assert from 'assert';
import * as vscode from 'vscode';
import * as sinon from 'sinon';
import { activate } from '../../extension';

/*
suite('Web Extension Test Suite', () => {
  let executeCommandStub: sinon.SinonStub;
  let registerCustomEditorProviderStub: sinon.SinonStub;
  let context: vscode.ExtensionContext;

  setup(() => {
    // Stub vscode.commands.executeCommand generally
    executeCommandStub = sinon.stub(vscode.commands, 'executeCommand');
    registerCustomEditorProviderStub = sinon.stub(vscode.window, 'registerCustomEditorProvider');

    context = {
      subscriptions: [],
      workspaceState: { get: () => {}, update: () => Promise.resolve(), keys: () => [] },
      globalState: { get: () => {}, update: () => Promise.resolve(), setKeysForSync: () => {}, keys: () => [] },
      extensionPath: '',
      storagePath: '',
      globalStoragePath: '',
      logPath: '',
      asAbsolutePath: (relativePath: string) => relativePath,
      extensionUri: vscode.Uri.parse('file:///mock'),
      storageUri: vscode.Uri.parse('file:///mock'),
      globalStorageUri: vscode.Uri.parse('file:///mock'),
      logUri: vscode.Uri.parse('file:///mock'),
      extensionMode: vscode.ExtensionMode.Test,
      environmentVariableCollection: {} as any,
      secrets: { get: () => Promise.resolve(undefined), store: () => Promise.resolve(), delete: () => Promise.resolve(), onDidChange: new vscode.EventEmitter<vscode.SecretStorageChangeEvent>().event, keys: () => Promise.resolve([]) },
      extension: {} as any,
      languageModelAccessInformation: {} as any,
    } as vscode.ExtensionContext;
    activate(context);
  });

  teardown(() => {
    executeCommandStub.restore();
    registerCustomEditorProviderStub.restore();
    context.subscriptions.forEach(s => s.dispose());
  });

  test('note-app-vs-code.openNote command should open a custom editor', async () => {
    const testNoteId = 'test-note-id-123';
    
    // Execute the command, which should trigger the handler registered in extension.ts
    await vscode.commands.executeCommand('note-app-vs-code.openNote', testNoteId);

    // Find the call to 'vscode.openWith'
    const openWithCall = executeCommandStub.getCalls().find(call => call.args[0] === 'vscode.openWith');

    assert.ok(openWithCall, 'vscode.openWith should have been called');
    assert.strictEqual(openWithCall.args[0], 'vscode.openWith', 'The command should be vscode.openWith');
    assert.ok(openWithCall.args[1] instanceof vscode.Uri, 'The first argument should be a Uri');
    assert.strictEqual(openWithCall.args[1].scheme, 'note', 'The Uri scheme should be \'note\'');
    assert.strictEqual(openWithCall.args[1].path, testNoteId, 'The Uri path should be the noteId');
    assert.strictEqual(openWithCall.args[2], NoteEditorProvider.viewType, 'The viewType should be correct');
  });

	test('Sample test', () => {
		assert.strictEqual(-1, [1, 2, 3].indexOf(5));
		assert.strictEqual(-1, [1, 2, 3].indexOf(0));
	});
});
*/