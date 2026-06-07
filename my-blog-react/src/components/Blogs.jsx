import { Link } from 'react-router-dom';

export default function Blogs({blogs, loading, error}) {
    // กันพัง: ถ้ายังไม่มีข้อมูล ใช้ลิสต์ว่างแทน (เลี่ยง blogs.data ที่อาจ undefined)
    const list = blogs && blogs.data ? blogs.data : [];

    return (
        <div className='w-full bg-[#f9f9f9] py-[50px] min-h-[300px]'>
        <div className='max-w-[1240px] mx-auto'>

            {/* แสดงสถานะเฉพาะพื้นที่รายการ blog — Navbar/Footer ยังอยู่ครบ */}
            {loading && (
                <p className='text-center text-gray-500 text-xl'>กำลังโหลดบทความ...</p>
            )}
            {!loading && error && (
                <p className='text-center text-gray-500 text-xl'>
                    ขณะนี้ยังโหลดบทความไม่ได้ (เซิร์ฟเวอร์ออฟไลน์)
                </p>
            )}
            {!loading && !error && list.length === 0 && (
                <p className='text-center text-gray-500 text-xl'>ยังไม่มีบทความ</p>
            )}

            <div className='grid lg:grid-cols-3 md:grid-cols-2 sm:grid-cols-2 ss:grid-cols-1 gap-8 px-4 text-black'>
            {list.map((blog)=>
                <Link key={blog.blogId} to={`/blog/${blog.blogId}`}>
                    <div  className='bg-white rounded-xl overflow-hidden drop-shadow-md'>
                        <img className='h-56 w-full object-cover' src={blog.coverImage} />
                        <div className='p-8'>
                            <h3 className='font-bold text-2xl my-1'>{blog.title}</h3>
                            <p className='text-gray-600 text-xl'>{blog.blogDesc}</p>
                        </div>
                    </div>
                </Link>

            )}

            </div>
        </div>
    </div>

    )
}
