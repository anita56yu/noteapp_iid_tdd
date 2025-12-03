import axios from 'axios';

export class AuthService {
  private static instance: AuthService;
  private baseUrl = 'http://localhost:8080'; // Adjust if your backend URL is different

  private constructor() {}

  public static getInstance(): AuthService {
    if (!AuthService.instance) {
      AuthService.instance = new AuthService();
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
}
