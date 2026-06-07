import Navbar  from '../components/Navbar'
import BlogContent from '../components/BlogContent'
import Footer from '../components/Footer'
import React from 'react'

function BlogContentPage({blogs}) {
  console.log(blogs)
  return (
      <div>
        <Navbar />
        <BlogContent blogs={blogs}/>   
        <Footer />
      </div> 
  )
}

export default BlogContentPage