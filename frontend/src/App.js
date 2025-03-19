import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage.js';
import WordEditorPage from './pages/WordEditorPage.js';
import Navigation from './components/Navigation.js';
import './App.css';
export default function App() {
  return (
    <Router>
      <Navigation />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/word-editor" element={<WordEditorPage />} />
      </Routes>
    </Router>
  );
}