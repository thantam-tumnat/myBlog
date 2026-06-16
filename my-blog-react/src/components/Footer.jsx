import { FaGithub } from 'react-icons/fa';
import { HiOutlineMail } from 'react-icons/hi';

const GITHUB_URL = 'https://github.com/thantam-tumnat';
const EMAIL = 'ca.tumnat@gmail.com';

export default function Footer() {
  return (
    <footer className="border-t border-line bg-surface">
      <div className="wrap flex flex-col items-center justify-between gap-4 py-8 sm:flex-row">
        <p className="text-sm text-muted">
          © {new Date().getFullYear()} Thantam Tumnat · built with Go microservices
        </p>
        <div className="flex items-center gap-4 text-muted">
          <a
            href={GITHUB_URL}
            target="_blank"
            rel="noreferrer"
            aria-label="GitHub"
            className="transition-colors hover:text-ink"
          >
            <FaGithub size={20} />
          </a>
          <a
            href={`mailto:${EMAIL}`}
            aria-label="Email"
            className="transition-colors hover:text-ink"
          >
            <HiOutlineMail size={22} />
          </a>
        </div>
      </div>
    </footer>
  );
}
