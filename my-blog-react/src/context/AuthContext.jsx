import { createContext, useContext, useEffect, useState } from 'react';
import {
  getToken, setToken, clearToken,
  decodeToken, isTokenValid,
  loginRequest, registerRequest,
} from '../api/auth';

const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null); // { user_id, username, role }

  // โหลด token เดิมตอนเปิดแอป — ถ้ายัง valid ถือว่ายัง login อยู่
  useEffect(() => {
    const token = getToken();
    if (token && isTokenValid(token)) {
      setUser(decodeToken(token));
    } else {
      clearToken();
    }
  }, []);

  const login = async (username, password) => {
    const data = await loginRequest(username, password);
    setToken(data.accessToken);
    setUser(decodeToken(data.accessToken));
  };

  const register = async (payload) => {
    return registerRequest(payload);
  };

  const logout = () => {
    clearToken();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, isAuthenticated: !!user, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => useContext(AuthContext);
