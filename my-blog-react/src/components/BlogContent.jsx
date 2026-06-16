import { useParams, Link } from 'react-router-dom';

const FALLBACK_IMG =
  'data:image/svg+xml;utf8,' +
  encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="600"><rect width="100%" height="100%" fill="#f4f4f5"/><text x="50%" y="50%" fill="#a1a1aa" font-family="sans-serif" font-size="28" text-anchor="middle" dominant-baseline="middle">no image</text></svg>'
  );

export default function BlogsContent({ blogs }) {
  const { id } = useParams();

  // กันพัง: ถ้าไม่มีข้อมูล ใช้ลิสต์ว่าง แล้วหาไม่เจอก็เป็น object ว่าง
  const list = blogs && blogs.data ? blogs.data : [];
  const blog = list.find((item) => item.blogId == id) || {};

  return (
    <article className="wrap max-w-prose py-12">
      <Link to="/" className="mb-8 inline-flex items-center gap-1 text-sm text-muted hover:text-ink">
        <span>←</span> Back to posts
      </Link>

      <h1 className="text-3xl leading-tight sm:text-4xl">{blog.title}</h1>

      {/* ผู้เขียน */}
      {blog.userName && (
        <div className="mt-5 flex items-center gap-3 border-b border-line pb-8">
          <img
            src={blog.userImage || FALLBACK_IMG}
            onError={(e) => { e.currentTarget.src = FALLBACK_IMG; }}
            alt={blog.userName}
            className="h-11 w-11 rounded-full object-cover ring-1 ring-line"
          />
          <div>
            <p className="font-medium text-ink">{blog.userName}</p>
            <p className="text-sm text-muted">{blog.userDesc}</p>
          </div>
        </div>
      )}

      {/* รูปปก */}
      <img
        src={blog.coverImage || FALLBACK_IMG}
        onError={(e) => { e.currentTarget.src = FALLBACK_IMG; }}
        alt={blog.title}
        className="mt-8 aspect-[16/9] w-full rounded-2xl object-cover"
      />

      {/* เนื้อหา */}
      <div className="prose-content line-break mt-8">{blog.content}</div>
    </article>
  );
}
