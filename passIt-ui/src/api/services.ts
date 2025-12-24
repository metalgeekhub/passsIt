import authAxios, { publicAxios } from './axios';
import type { User } from '@/models/User';

/**
 * User API Service
 * All user-related API calls
 */
export const userApi = {
  /**
   * Get all users (requires authentication)
   */
  async getAll(): Promise<User[]> {
    const response = await authAxios.get<User[]>('/api/users');
    return response.data;
  },

  /**
   * Get user by ID (requires authentication)
   */
  async getById(id: string): Promise<User> {
    const response = await authAxios.get<User>(`/api/users/find?id=${id}`);
    return response.data;
  },

  /**
   * Get user by email (requires authentication)
   */
  async getByEmail(email: string): Promise<User> {
    const response = await authAxios.get<User>(`/api/users/by-email?email=${encodeURIComponent(email)}`);
    return response.data;
  },

  /**
   * Create a new user (requires authentication)
   */
  async create(userData: {
    email: string;
    username: string;
    firstName: string;
    lastName: string;
    password: string;
    isAdmin?: boolean;
  }): Promise<User> {
    // Transform camelCase to snake_case for backend
    const requestData = {
      email: userData.email,
      username: userData.username,
      first_name: userData.firstName,
      last_name: userData.lastName,
      password: userData.password,
      is_admin: userData.isAdmin || false,
    };
    const response = await authAxios.post<User>('/api/users', requestData);
    return response.data;
  },

  /**
   * Update user by ID (requires authentication)
   */
  async update(id: string, userData: {
    email?: string;
    username?: string;
    firstName?: string;
    lastName?: string;
    dob?: string;
    phoneNumber?: string;
    address?: string;
    isAdmin?: boolean;
    password?: string;
  }): Promise<User> {
    // Transform camelCase to snake_case for backend
    const requestData: any = {};
    if (userData.email) requestData.email = userData.email;
    if (userData.username) requestData.username = userData.username;
    if (userData.firstName) requestData.first_name = userData.firstName;
    if (userData.lastName) requestData.last_name = userData.lastName;
    if (userData.dob) requestData.dob = userData.dob;
    if (userData.phoneNumber) requestData.phone_number = userData.phoneNumber;
    if (userData.address) requestData.address = userData.address;
    if (userData.isAdmin !== undefined) requestData.is_admin = userData.isAdmin;
    if (userData.password) requestData.password = userData.password;
    
    const response = await authAxios.put<User>(`/api/users/${id}`, requestData);
    return response.data;
  },

  /**
   * Delete user by ID - soft delete (requires admin authentication)
   */
  async delete(id: string): Promise<void> {
    await authAxios.delete(`/api/users/${id}`);
  },

  /**
   * Get all inactive/deleted users (requires admin authentication)
   */
  async getInactive(): Promise<User[]> {
    const response = await authAxios.get<User[]>('/api/users/inactive');
    return response.data;
  },
};

/**
 * Authentication API Service
 */
export const authApi = {
  /**
   * Check if user is authenticated
   */
  async checkAuth(): Promise<boolean> {
    try {
      await authAxios.get('/api/users');
      return true;
    } catch {
      return false;
    }
  },

  /**
   * Get current authenticated user
   */
  async getCurrentUser(): Promise<User> {
    const response = await authAxios.get<User>('/api/users/me');
    return response.data;
  },

  /**
   * Public signup (no authentication required)
   */
  async signup(userData: {
    username: string;
    email: string;
    firstName?: string;
    lastName?: string;
    password: string;
  }): Promise<{ message: string; user: any }> {
    // Transform camelCase to snake_case for backend
    const requestData = {
      username: userData.username,
      email: userData.email,
      first_name: userData.firstName || '',
      last_name: userData.lastName || '',
      password: userData.password,
    };
    const response = await publicAxios.post('/auth/signup', requestData);
    return response.data;
  },

  /**
   * Redirect to login page
   */
  login(): void {
    window.location.href = `${authAxios.defaults.baseURL}/auth/login`;
  },

  /**
   * Redirect to logout page
   */
  logout(): void {
    window.location.href = `${authAxios.defaults.baseURL}/auth/logout`;
  },
};

/**
 * Health check API
 */
export const healthApi = {
  /**
   * Check backend health status
   */
  async check(): Promise<any> {
    const response = await publicAxios.get('/health');
    return response.data;
  },
};

// Export all APIs
export default {
  user: userApi,
  auth: authApi,
  health: healthApi,
};
