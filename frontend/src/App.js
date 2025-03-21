import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage.js';
import WordEditorPage from './pages/WordEditorPage.js';
import WordSearchPage from './pages/WordSearchPage.js';
import Navigation from './components/Navigation.js';
import './App.css';
export default function App() {
  return (
    <Router>
      <Navigation />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/word-editor" element={<WordEditorPage />} />
        <Route path="/word-search" element={<WordSearchPage />} /> {/* 词表搜索页面 */}
      </Routes>
    </Router>
  );
}