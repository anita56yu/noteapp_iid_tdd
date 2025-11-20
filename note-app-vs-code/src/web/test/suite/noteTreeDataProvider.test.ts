import * as assert from 'assert';
import * as vscode from 'vscode';
import { NoteTreeDataProvider, NoteTreeItem } from '../../noteTreeDataProvider';
import { NoteService, Note } from '../../noteService';
import { SinonStub, stub } from 'sinon';

suite('NoteTreeDataProvider Test Suite', () => {
  let noteService: NoteService;
  let noteTreeDataProvider: NoteTreeDataProvider;
  let getNotesStub: SinonStub;

  setup(() => {
    // Create a new instance of the service for each test
    noteService = NoteService.getInstance();
    // Stub the getNotes method before creating the provider
    getNotesStub = stub(noteService, 'getNotes');
    noteTreeDataProvider = new NoteTreeDataProvider(noteService);
  });

  teardown(() => {
    // Restore the original method after each test
    getNotesStub.restore();
  });

  test('getTreeItem should return the provided element', () => {
    const noteItem = new NoteTreeItem('1', 'Test Note', vscode.TreeItemCollapsibleState.None);
    const treeItem = noteTreeDataProvider.getTreeItem(noteItem);
    assert.deepStrictEqual(treeItem, noteItem);
  });

  test('getChildren should return notes from NoteService when no element is provided', async () => {
    const mockNotes: Note[] = [
      { id: '1', title: 'Note 1', content_ids: [], version: 1 },
      { id: '2', title: 'Note 2', content_ids: [], version: 1 },
    ];
    // Configure the stub to return mock data
    getNotesStub.resolves(mockNotes);

    const children = await noteTreeDataProvider.getChildren();

    assert.strictEqual(children.length, 2, 'Should return two children');
    assert.strictEqual(children[0].noteId, '1');
    assert.strictEqual(children[0].label, 'Note 1');
    assert.strictEqual(children[1].noteId, '2');
    assert.strictEqual(children[1].label, 'Note 2');
    assert.ok(getNotesStub.calledOnceWith("testUser1"), 'getNotes should be called once with the correct user ID');
  });

  test('getChildren should return an empty array when an element is provided', async () => {
    const noteItem = new NoteTreeItem('1', 'Test Note', vscode.TreeItemCollapsibleState.None);
    const children = await noteTreeDataProvider.getChildren(noteItem);
    assert.strictEqual(children.length, 0, 'Should return no children for a note item');
  });

  test('refresh should fire onDidChangeTreeData event', (done) => {
    const disposable = noteTreeDataProvider.onDidChangeTreeData(() => {
      // This callback being called is the success condition
      disposable.dispose(); // Clean up the listener
      done(); // Signal that the test is complete
    });

    noteTreeDataProvider.refresh();
  });
});
