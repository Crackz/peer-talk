export interface AuthResponse {
  accessToken: string;
  user: {
    id: string;
    name: string;
    username: string;
  };
}
