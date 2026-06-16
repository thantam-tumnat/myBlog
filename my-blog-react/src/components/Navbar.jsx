import { Link } from 'react-router-dom';

// ── ข้อมูลโปรไฟล์ (แก้ค่าตรงนี้ที่เดียว) ──
const PROFILE_NAME = 'Thantam Tumnat';
const PROFILE_IMG = 'https://i.ibb.co/99wyY0vR/0-BADD6-F0-D2-C0-4151-A45-A-F3-F5-F7443-B2-D.jpg';
const GITHUB_URL = 'https://github.com/thantam-tumnat';

export default function Navbar() {
  return (
    <header className="sticky top-0 z-50 border-b border-line bg-white/80 backdrop-blur">
      <nav className="wrap flex h-16 items-center justify-between">
        {/* แบรนด์ — ชื่อ + รูปโปรไฟล์ ลิงก์กลับหน้าแรก */}
        <Link to="/" className="flex items-center gap-3">
          <img
            src={PROFILE_IMG}
            alt={PROFILE_NAME}
            className="h-9 w-9 rounded-full object-cover ring-1 ring-line"
          />
          <span className="font-semibold tracking-tight text-ink">{PROFILE_NAME}</span>
        </Link>

        {/* เมนู */}
        <div className="flex items-center gap-1 sm:gap-2">
          <Link to="/" className="btn-ghost">Home</Link>
          <Link to="/pdf-viewer" className="btn-ghost">Resume</Link>
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
