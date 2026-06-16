import Navbar from '../components/Navbar'
import Hero from '../components/Hero'
import Blogs from '../components/Blogs'
import Footer from '../components/Footer'
import React from 'react'

function Homepage({ blogs, loading, error }) {
  return (
    <div className="flex min-h-screen flex-col">
      <Navbar />
      <main className="flex-1">
        <Hero />
        <Blogs blogs={blogs} loading={loading} error={error} />
      </main>
      <Footer />
    </div>
  )
}

export default Homepage
