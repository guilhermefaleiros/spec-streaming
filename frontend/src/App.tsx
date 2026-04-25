import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { HomePage } from './pages/home-page'
import { VideoPage } from './pages/video-page'

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/videos/:id" element={<VideoPage />} />
      </Routes>
    </BrowserRouter>
  )
}
