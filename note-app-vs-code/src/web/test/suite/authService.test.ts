import * as assert from 'assert';
import axios from 'axios';
import * as sinon from 'sinon';
import { AuthService } from '../../authService';

suite('AuthService Test Suite', () => {
  let axiosPostStub: sinon.SinonStub;

  setup(() => {
    axiosPostStub = sinon.stub(axios, 'post');
  });

  teardown(() => {
    axiosPostStub.restore();
  });

  test('register should call the correct endpoint with correct data', async () => {
    const authService = AuthService.getInstance();
    const username = 'testuser';
    const password = 'password';
    const expectedUrl = 'http://localhost:8080/register';
    const expectedPayload = { username, password };

    axiosPostStub.resolves({ data: { id: '1', username: 'testuser' } });

    await authService.register(username, password);

    assert.ok(axiosPostStub.calledOnceWith(expectedUrl, expectedPayload));
  });

  test('register should throw an error when API call fails', async () => {
    const authService = AuthService.getInstance();
    const username = 'testuser';
    const password = 'password';
    const errorMessage = 'Network Error';

    axiosPostStub.rejects(new Error(errorMessage));

    try {
      await authService.register(username, password);
      assert.fail('The register method should have thrown an error.');
    } catch (error: any) {
      assert.strictEqual(error.message, errorMessage);
    }
  });
});
