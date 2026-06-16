import Navbar from '../components/Navbar'
import BlogContent from '../components/BlogContent'
import Footer from '../components/Footer'
import React from 'react'

function BlogContentPage({ blogs }) {
  return (
    <div className="flex min-h-screen flex-col">
      <Navbar />
      <main className="flex-1">
        <BlogContent blogs={blogs} />
      </main>
      <Footer />
    </div>
  )
}

export default BlogContentPage
