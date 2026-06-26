import React from 'react'
import Homepage from './pages/HomePage'
import BlogContentPage from './pages/BlogContentPage'
// import { Routes, Route } from 'react-router-dom'
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import useFetch from './hooks/useFetch'
import PdfPage from './pages/PdfPage';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import { AuthProvider } from './context/AuthContext';

function App() {

  const BLOG_API = import.meta.env.VITE_BLOG_API || 'http://localhost:8002/v1/myblogs'
  let {loading, data, error} = useFetch(`${BLOG_API}/getBlogs`)

  // ส่งสถานะ loading/error เข้าไปในหน้า แทนการบล็อกทั้งแอป
  // เพื่อให้ Navbar/Footer (ปุ่มโปรไฟล์, About us) แสดงเสมอ แม้ backend ofline
  return (
    <AuthProvider>
    <Router>
    <div>
      <Routes>
        <Route path='/' element={<Homepage blogs={data} loading={loading} error={error} />}></Route>
        <Route path='/blog/:id' element={<BlogContentPage blogs={data} loading={loading} error={error} />}></Route>
        <Route path='/pdf-viewer' element={<PdfPage></PdfPage>}></Route>
        <Route path='/login' element={<LoginPage />}></Route>
        <Route path='/register' element={<RegisterPage />}></Route>
      </Routes>
    </div>
    </Router>
    </AuthProvider>
  )
}

export default App
