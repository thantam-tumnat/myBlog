import {FaFacebook, FaGithub, FaInstagram, FaTwitter, FaTwitch} from 'react-icons/fa';

export default function Footer() {
    return (
        <div className='w-full bg-myblogbg text-gray-300 py-8 px-2'>
        {/* <div className='max-w-[1240px] mx-auto grid grid-cols-2 md:grid-cols-3 border-b-2 border-gray-600 py-8'>
            <div>
                <h6 className='font-bold uppercase py-2'>MyBlogs</h6>
                <ol>
                    <li className='py-1'>Your blogs</li>
                    <li className='py-1'>Let's blogs</li>
                </ol>
            </div>
            <div>
                <h6 className='font-bold uppercase py-2'>About us</h6>
                <ol>
                    <li className='py-1'>Contact us</li>
                    <li className='py-1'>087-233-1295</li>
                </ol>
            </div>
            <div>
                <h6 className='font-bold uppercase py-2'>Terms and Conditions</h6>
                <ol>
                    <li className='py-1'>Privacy Policy</li>
                    <li className='py-1'>User Agreement</li>
                </ol>
            </div>

        </div> */}

        <div className='flex flex-col max-w-[1240px] px-2 py-4 m-auto justify-between sm:flex-row text-center text-gray-500 items-center'>
            <p>2024 MYBLOGS, LLC. All rights reserved.</p>
            <div className='flex justify-between sm:w-[300px] pt-4 text-2xl gap-2'>
                <FaFacebook />
                <FaGithub />
                <FaInstagram />
                <FaTwitch />
                <FaTwitter />
            </div>
        </div>
    </div>
  
    )
}
