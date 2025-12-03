import axios from 'axios';
import * as vscode from 'vscode';

const AUTH_TOKEN_KEY = 'note-app-jwt';

export class AuthService {
  private static instance: AuthService;
  private baseUrl = 'http://localhost:8080';

  private constructor(private readonly context: vscode.ExtensionContext) {}

  public static initialize(context: vscode.ExtensionContext): void {
    if (!AuthService.instance) {
      AuthService.instance = new AuthService(context);
    }
  }

  public static getInstance(): AuthService {
    if (!AuthService.instance) {
      throw new Error('AuthService has not been initialized. Call AuthService.initialize() first.');
    }
    return AuthService.instance;
  }

  async register(username: string, password: string): Promise<any> {
    try {
      const response = await axios.post(`${this.baseUrl}/register`, { username, password });
      return response.data;
    } catch (error) {
      console.error('Error during registration:', error);
      throw error;
    }
  }

  async login(username: string, password: string): Promise<string | undefined> {
    try {
      const response = await axios.post<{ token: string }>(`${this.baseUrl}/login`, { username, password });
      const token = response.data.token;
      if (token) {
        await this.setToken(token);
        return;
      }
      return undefined;
    } catch (error) {
      console.error('Error during login:', error);
      throw error;
    }
  }

  async setToken(token: string): Promise<void> {
    await this.context.secrets.store(AUTH_TOKEN_KEY, token);
  }

  async getToken(): Promise<string | undefined> {
    return this.context.secrets.get(AUTH_TOKEN_KEY);
  }
}
