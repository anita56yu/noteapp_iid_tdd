import * as assert from 'assert';
import axios from 'axios';
import * as sinon from 'sinon';
import * as vscode from 'vscode';
import { AuthService } from '../../authService';

suite('AuthService Test Suite', () => {
  let axiosPostStub: sinon.SinonStub;
  let mockContext: vscode.ExtensionContext;
  let secretStorage: { [key: string]: string };

  setup(() => {
    axiosPostStub = sinon.stub(axios, 'post');
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

  test('login should call the correct endpoint and store the token', async () => {
    const authService = AuthService.getInstance();
    const username = 'testuser';
    const password = 'password';
    const token = 'fake-jwt-token';
    const expectedUrl = 'http://localhost:8080/login';
    const expectedPayload = { username, password };

    axiosPostStub.resolves({ data: { token } });

    await authService.login(username, password);

    assert.ok(axiosPostStub.calledOnceWith(expectedUrl, expectedPayload));
    const storedToken = await authService.getToken();
    assert.strictEqual(storedToken, token);
  });
});
