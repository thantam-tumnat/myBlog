import { Link } from 'react-router-dom';

// ── ข้อมูลโปรไฟล์ (แก้ค่าตรงนี้ที่เดียว) ──
const PROFILE_NAME = "Thantam Tumnat";
const PROFILE_IMG = "https://i.ibb.co/99wyY0vR/0-BADD6-F0-D2-C0-4151-A45-A-F3-F5-F7443-B2-D.jpg";

export default function Navbar() {
    return (
        <div className='w-full h-[80px] z-10 bg-white fixed drop-shadow-lg relative sticky top-0'>
            <div className='flex justify-between items-center w-full h-full md:max-w-[1240px] m-auto'>
                
                <div className='flex items-center'>
                    {/* <img src="/src/assets/Myblog.png" alt="logo" className='sm:ml-10 ss:ml-10 md:ml-3 w-full h-[40px]' /> */}
                    <Link to="/pdf-viewer">
                        <button className="flex items-center hover:opacity-[70%] rounded-md">
                            <img className='p-4 w-20 h-20 rounded-full' src={PROFILE_IMG} />
                            <h1 className='font-bold text-1xl text-center text-myblogbg pt-0'>{PROFILE_NAME}</h1>
                        </button>
                    </Link>
                </div>


                <div className='hidden md:flex sm:mr-10 md:mr-10'>
                    <Link to="/">
                        <button className="flex items-center hover:opacity-[70%] rounded-md">
                        <h1 className='pr-5 font-bold text-1xl text-center text-myblogbg pt-0'>{"Home"}</h1>
                        </button>
                    </Link>
                    <Link to="/pdf-viewer">
                        <button className="flex items-center hover:opacity-[70%] rounded-md">
                            <h1 className='pr-5 font-bold text-1xl text-center text-myblogbg pt-0'>{"About us"}</h1>
                        </button>
                    </Link>
                    <a href="https://github.com/thantam-tumnat" className="hover:opacity-[70%] rounded-md font-bold text-1xl text-center text-myblogbg pt-0">{"Github 📤"}</a>
                </div>

            </div>
        </div>
        
    )
}



