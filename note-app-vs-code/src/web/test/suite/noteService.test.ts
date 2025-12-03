import * as assert from 'assert';
import axios from 'axios';
import * as sinon from 'sinon';
import { NoteService, Note } from '../../noteService';
import { AuthService } from '../../authService';
import * as vscode from 'vscode';

suite('NoteService Test Suite', () => {
  let axiosGetStub: sinon.SinonStub;
  let authServiceGetTokenStub: sinon.SinonStub;
  let mockContext: vscode.ExtensionContext;
  let secretStorage: { [key: string]: string };

  setup(() => {
    axiosGetStub = sinon.stub(axios, 'get');
    secretStorage = {};

    mockContext = {
      secrets: {
        get: async (key: string): Promise<string | undefined> => secretStorage[key],
        store: async (key: string, value: string): Promise<void> => {
          secretStorage[key] = value;
        },
        delete: async (key: string): Promise<void> => {
          delete secretStorage[key];
        },
        onDidChange: sinon.stub(),
      },
    } as any;
    AuthService.initialize(mockContext);
    authServiceGetTokenStub = sinon.stub(AuthService.getInstance(), 'getToken');
  });

  teardown(() => {
    axiosGetStub.restore();
    authServiceGetTokenStub.restore();
  });

  test('getNotes should fetch notes for an authenticated user', async () => {
    const noteService = NoteService.getInstance();
    const token = 'fake-jwt-token';
    const expectedNotes: Note[] = [
      { id: '1', title: 'Note 1', content_ids: [], version: 1 },
      { id: '2', title: 'Note 2', content_ids: [], version: 1 },
    ];

    authServiceGetTokenStub.resolves(token);
    axiosGetStub.withArgs('http://localhost:8080/notes/accessible-notes', sinon.match.any).resolves({ data: expectedNotes });

    const notes = await noteService.getNotes();

    assert.deepStrictEqual(notes, expectedNotes);
    assert.ok(axiosGetStub.calledOnceWith('http://localhost:8080/notes/accessible-notes', {
      headers: { Authorization: `Bearer ${token}` },
    }));
  });

  test('getNotes should return an empty array on error', async () => {
    const noteService = NoteService.getInstance();
    const token = 'fake-jwt-token';

    authServiceGetTokenStub.resolves(token);
    axiosGetStub.withArgs('http://localhost:8080/notes/accessible-notes', sinon.match.any).rejects(new Error('Network Error'));

    const notes = await noteService.getNotes();

    assert.deepStrictEqual(notes, []);
    assert.ok(axiosGetStub.calledOnce);
  });

  test('getNotes should return an empty array if no token is available', async () => {
    const noteService = NoteService.getInstance();
    authServiceGetTokenStub.resolves(undefined);

    const notes = await noteService.getNotes();

    assert.deepStrictEqual(notes, []);
    assert.ok(axiosGetStub.notCalled);
  });
});
