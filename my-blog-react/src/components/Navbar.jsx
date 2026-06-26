import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

// ── แบรนด์เว็บ (ชื่อ product ไม่ใช่ชื่อเจ้าของ) ──
const BRAND_NAME = 'MyBlog';
const GITHUB_URL = 'https://github.com/thantam-tumnat';

export default function Navbar() {
  const { isAuthenticated, user, logout } = useAuth();

  return (
    <header className="sticky top-0 z-50 border-b border-line bg-white/80 backdrop-blur">
      <nav className="wrap flex h-16 items-center justify-between">
        {/* แบรนด์ — wordmark ของ product ลิงก์กลับหน้าแรก */}
        <Link to="/" className="flex items-center gap-2">
          <span className="grid h-8 w-8 place-items-center rounded-lg bg-ink text-sm font-bold text-white">
            M
          </span>
          <span className="text-lg font-bold tracking-tight text-ink">{BRAND_NAME}</span>
        </Link>

        {/* เมนู */}
        <div className="flex items-center gap-1 sm:gap-2">
          <Link to="/" className="btn-ghost">Home</Link>
          <Link to="/pdf-viewer" className="btn-ghost">Resume</Link>

          {isAuthenticated ? (
            <>
              <span className="hidden px-2 text-sm text-muted ss:inline">
                สวัสดี, <span className="font-medium text-ink">{user?.username}</span>
              </span>
              <button onClick={logout} className="btn-ghost">Logout</button>
            </>
          ) : (
            <>
              <Link to="/login" className="btn-ghost">Login</Link>
              <Link to="/register" className="btn-ghost">Sign up</Link>
            </>
          )}

          <a
            href={GITHUB_URL}
            target="_blank"
            rel="noreferrer"
            className="rounded-lg bg-ink px-3 py-2 text-sm font-medium text-white transition-opacity hover:opacity-80"
          >
            GitHub
          </a>
        </div>
      </nav>
    </header>
  );
}
