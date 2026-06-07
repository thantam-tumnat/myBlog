import Navbar  from '../components/Navbar'
import Blogs from '../components/Blogs'
import Footer from '../components/Footer'
import React from 'react'

function Homepage({blogs, loading, error}) {
  return (
      <div>
        <Navbar />
        <Blogs blogs={blogs} loading={loading} error={error} />
        <Footer />
      </div>
  )
}

export default Homepage
