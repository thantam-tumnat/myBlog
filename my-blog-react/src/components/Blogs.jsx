import { Link } from 'react-router-dom';

// รูปสำรองเมื่อ coverImage ว่างหรือโหลดไม่ขึ้น
const FALLBACK_IMG =
  'data:image/svg+xml;utf8,' +
  encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" width="600" height="400"><rect width="100%" height="100%" fill="#f4f4f5"/><text x="50%" y="50%" fill="#a1a1aa" font-family="sans-serif" font-size="20" text-anchor="middle" dominant-baseline="middle">no image</text></svg>'
  );

export default function Blogs({ blogs, loading, error }) {
  // กันพัง: ถ้ายังไม่มีข้อมูล ใช้ลิสต์ว่างแทน
  const list = blogs && blogs.data ? blogs.data : [];

  return (
    <section className="wrap border-t border-line py-14">
      <div className="mb-8 flex items-baseline justify-between">
        <h2 className="text-2xl sm:text-3xl">Latest posts</h2>
        {!loading && !error && list.length > 0 && (
          <span className="text-sm text-muted">{list.length} บทความ</span>
        )}
      </div>

      {/* สถานะ */}
      {loading && <p className="py-10 text-center text-muted">กำลังโหลดบทความ…</p>}
      {!loading && error && (
        <p className="py-10 text-center text-muted">
          ขณะนี้ยังโหลดบทความไม่ได้ (เซิร์ฟเวอร์ออฟไลน์)
        </p>
      )}
      {!loading && !error && list.length === 0 && (
        <p className="py-10 text-center text-muted">ยังไม่มีบทความ</p>
      )}

      {/* การ์ดบล็อก — ซ้ายชิด ขนาดสม่ำเสมอ hover ยกขึ้น */}
      <div className="grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
        {list.map((blog) => (
          <Link
            key={blog.blogId}
            to={`/blog/${blog.blogId}`}
            className="group flex flex-col overflow-hidden rounded-2xl border border-line bg-white transition-all duration-200 hover:-translate-y-1 hover:shadow-lg hover:shadow-zinc-200"
          >
            <div className="aspect-[16/10] overflow-hidden bg-surface">
              <img
                src={blog.coverImage || FALLBACK_IMG}
                onError={(e) => {
                  e.currentTarget.src = FALLBACK_IMG;
                }}
                alt={blog.title}
                className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
              />
            </div>
            <div className="flex flex-1 flex-col p-6">
              <h3 className="text-xl font-bold leading-snug text-ink transition-colors group-hover:text-accent">
                {blog.title}
              </h3>
              <p className="mt-2 line-clamp-3 text-body">{blog.blogDesc}</p>
              <span className="mt-4 inline-flex items-center gap-1 text-sm font-medium text-accent">
                Read more
                <span className="transition-transform group-hover:translate-x-1">→</span>
              </span>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}
