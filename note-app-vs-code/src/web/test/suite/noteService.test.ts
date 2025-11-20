import * as assert from 'assert';
import axios from 'axios';
import * as sinon from 'sinon';
import { NoteService, Note } from '../../noteService';

suite('NoteService Test Suite', () => {
  let axiosGetStub: sinon.SinonStub;

  setup(() => {
    axiosGetStub = sinon.stub(axios, 'get');
  });

  teardown(() => {
    axiosGetStub.restore();
  });

  test('getNotes should fetch notes for a user', async () => {
    const noteService = NoteService.getInstance();
    const userId = 'test-user';
    const expectedNotes: Note[] = [
      { id: '1', title: 'Note 1', content_ids: [], version: 1 },
      { id: '2', title: 'Note 2', content_ids: [], version: 1 },
    ];

    axiosGetStub.withArgs(`http://localhost:8080/users/${userId}/accessible-notes`).resolves({ data: expectedNotes });

    const notes = await noteService.getNotes(userId);

    assert.deepStrictEqual(notes, expectedNotes);
    assert.ok(axiosGetStub.calledOnce);
  });

  test('getNotes should return an empty array on error', async () => {
    const noteService = NoteService.getInstance();
    const userId = 'test-user';

    axiosGetStub.withArgs(`http://localhost:8080/users/${userId}/accessible-notes`).rejects(new Error('Network Error'));

    const notes = await noteService.getNotes(userId);

    assert.deepStrictEqual(notes, []);
    assert.ok(axiosGetStub.calledOnce);
  });
});
