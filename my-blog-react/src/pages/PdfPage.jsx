import Navbar  from '../components/Navbar'
// import Blogs from '../components/Blogs'
import Footer from '../components/Footer'
import React from 'react'



function PdfPage() {
  return (
      <div>
        <Navbar />
        <div className='pl-96 pr-96 flex items-center'>
            <embed src="/src/files/thantam_resume.pdf" type="application/pdf" width="100%" height="1200px" />
        </div>
        <Footer />
      </div> 
  )
}

export default PdfPage
