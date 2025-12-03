import * as assert from 'assert';
import * as vscode from 'vscode';
import * as sinon from 'sinon';
import { openNoteEditor } from '../../noteEditor';
import { Note, NoteService } from '../../noteService';

suite('NoteEditor Test Suite', () => {
  let createWebviewPanelStub: sinon.SinonStub;
  let postMessageStub: sinon.SinonStub;
  let showErrorMessageStub: sinon.SinonStub;
  let getNoteByIdStub: sinon.SinonStub;

  setup(() => {
    // Mock vscode.window.createWebviewPanel
    postMessageStub = sinon.stub();
    createWebviewPanelStub = sinon.stub(vscode.window, 'createWebviewPanel').returns({
      webview: {
        html: '',
        postMessage: postMessageStub,
      },
      dispose: sinon.stub(),
    } as any);

    // Mock vscode.window.showErrorMessage
    showErrorMessageStub = sinon.stub(vscode.window, 'showErrorMessage');

    // Mock NoteService.getInstance().getNoteById
    getNoteByIdStub = sinon.stub(NoteService.getInstance(), 'getNoteById');
  });

  teardown(() => {
    createWebviewPanelStub.restore();
    showErrorMessageStub.restore();
    getNoteByIdStub.restore();
  });

  test('openNoteEditor should create a webview panel and load note data', async () => {
    const noteId = 'test-note-id';
    const noteTitle = 'Test Note Title';
    const mockNote: Note & { contents: any[] } = {
      id: noteId,
      title: noteTitle,
      content_ids: ['content-1'],
      version: 1,
      contents: [{ id: 'content-1', note_id: noteId, content: 'Test content', content_type: 'text', version: 1 }],
    };

    getNoteByIdStub.withArgs(noteId).resolves(mockNote);

    openNoteEditor(noteId, noteTitle);

    // Ensure webview panel is created
    assert.ok(createWebviewPanelStub.calledOnce, 'createWebviewPanel should be called once');
    assert.strictEqual(createWebviewPanelStub.firstCall.args[1], noteTitle, 'Webview panel title should match note title');

    // Ensure getNoteById is called
    assert.ok(getNoteByIdStub.calledOnceWith(noteId), 'getNoteById should be called once with the correct noteId');

    // Wait for async operations to complete
    await new Promise(resolve => setTimeout(resolve, 0));

    // Ensure message is posted to webview
    assert.ok(postMessageStub.calledOnce, 'postMessage should be called once');
    assert.deepStrictEqual(postMessageStub.firstCall.args[0], {
      type: 'loadNote',
      payload: mockNote,
    }, 'postMessage should send loadNote event with correct payload');

    // Ensure no error message is shown
    assert.ok(showErrorMessageStub.notCalled, 'showErrorMessage should not be called');
  });

  test('openNoteEditor should show an error message if note fetching fails', async () => {
    const noteId = 'test-note-id';
    const noteTitle = 'Test Note Title';
    const errorMessage = 'Failed to fetch note';

    getNoteByIdStub.withArgs(noteId).rejects(new Error(errorMessage));

    openNoteEditor(noteId, noteTitle);

    // Wait for async operations to complete
    await new Promise(resolve => setTimeout(resolve, 0));

    // Ensure error message is shown
    assert.ok(showErrorMessageStub.calledOnce, 'showErrorMessage should be called once');
    assert.ok(showErrorMessageStub.firstCall.args[0].includes(errorMessage), 'Error message should contain the failure reason');

    // Ensure no message is posted to webview
    assert.ok(postMessageStub.notCalled, 'postMessage should not be called');
  });
});
