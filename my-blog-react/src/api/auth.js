import axios from 'axios';

// userservice (ออก/ตรวจ JWT). override ได้ผ่าน VITE_USER_API ใน .env
// หมายเหตุ: route group ของ userservice คือ /v1/myblogs (ไม่ใช่ /v1/users)
const USER_API = import.meta.env.VITE_USER_API || 'http://localhost:8001/v1/myblogs';
const TOKEN_KEY = 'myblogs_token';

// ── token storage ──────────────────────────────
export const getToken = () => localStorage.getItem(TOKEN_KEY);
export const setToken = (token) => localStorage.setItem(TOKEN_KEY, token);
export const clearToken = () => localStorage.removeItem(TOKEN_KEY);

// decode payload ของ JWT (ไม่ verify — แค่อ่าน username/role/exp มาแสดงผล)
export function decodeToken(token) {
  try {
    const payload = token.split('.')[1];
    return JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
  } catch {
    return null;
  }
}

// token ใช้ได้อยู่ไหม (มีจริง + ยังไม่หมดอายุ)
export function isTokenValid(token) {
  const claims = decodeToken(token);
  if (!claims || !claims.exp) return false;
  return claims.exp * 1000 > Date.now();
}

// ── API calls ──────────────────────────────────
export async function loginRequest(username, password) {
  const res = await axios.post(`${USER_API}/login`, { username, password });
  return res.data; // { accessToken, tokenType, ... }
}

export async function registerRequest(payload) {
  // payload: { username, password, name, description, userImage }
  const res = await axios.post(`${USER_API}/register`, payload);
  return res.data;
}
