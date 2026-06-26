import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import Navbar from '../components/Navbar';
import Footer from '../components/Footer';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();

  const [form, setForm] = useState({ username: '', password: '' });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const onChange = (e) => setForm({ ...form, [e.target.name]: e.target.value });

  const onSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);
    try {
      await login(form.username, form.password);
      navigate('/'); // login สำเร็จ กลับหน้าแรก
    } catch (err) {
      setError(err.response?.data?.message || 'เข้าสู่ระบบไม่สำเร็จ');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen flex-col">
      <Navbar />
      <main className="flex flex-1 items-center justify-center px-6 py-12">
        <div className="w-full max-w-sm rounded-2xl border border-line bg-white p-8 shadow-sm">
          <h1 className="text-2xl">เข้าสู่ระบบ</h1>
          <p className="mt-1 text-sm text-muted">ลงชื่อเข้าใช้เพื่อเขียนบล็อก</p>

          {error && (
            <div className="mt-4 rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-600">
              {error}
            </div>
          )}

          <form onSubmit={onSubmit} className="mt-6 space-y-4">
            <div>
              <label className="block text-sm font-medium text-ink">Username</label>
              <input
                name="username"
                value={form.username}
                onChange={onChange}
                autoComplete="username"
                required
                className="mt-1 w-full rounded-lg border border-line px-3 py-2 text-sm outline-none focus:border-accent focus:ring-1 focus:ring-accent"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-ink">Password</label>
              <input
                name="password"
                type="password"
                value={form.password}
                onChange={onChange}
                autoComplete="current-password"
                required
                className="mt-1 w-full rounded-lg border border-line px-3 py-2 text-sm outline-none focus:border-accent focus:ring-1 focus:ring-accent"
              />
            </div>
            <button
              type="submit"
              disabled={loading}
              className="w-full rounded-lg bg-ink px-4 py-2.5 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:opacity-50"
            >
              {loading ? 'กำลังเข้าสู่ระบบ...' : 'เข้าสู่ระบบ'}
            </button>
          </form>

          <p className="mt-6 text-center text-sm text-muted">
            ยังไม่มีบัญชี?{' '}
            <Link to="/register" className="font-medium text-accent hover:text-accentDark">
              สมัครสมาชิก
            </Link>
          </p>
        </div>
      </main>
      <Footer />
    </div>
  );
}
