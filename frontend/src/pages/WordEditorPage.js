import React, { useEffect, useState } from 'react';
import WordForm from '../components/WordEditor/WordForm';
import GrammarForm from '../components/WordEditor/GrammarForm';
import ReadingForm from '../components/WordEditor/ReadingForm';
import '../components/WordEditor/WordForm.css';
import { useLocation } from 'react-router-dom';

export default function WordEditorPage() {
  const location = useLocation();
  const [wordData, setWordData] = useState(null);
  const [bookId, setBookId] = useState(null);
  useEffect(() => {
    if (location.hash) {
      const id = location.hash.replace('#', '');
      const element = document.getElementById(id);
      if (element) {
        element.scrollIntoView({ behavior: 'smooth' });
      }
    }
    if (location.state && location.state.word) {
      const word = location.state.word;
      setWordData(word);
      setBookId(location.state.selectedBook);
    }
  })
  return (
    <div className="word-editor-page">
      <header className="main-header">
        <h1>词库管理</h1>
        <h2>Word Database Management</h2>
      </header>


      <main className="editor-content">
        <section id='word-form' className="word-form-section">
          <h3>添加/编辑单词</h3>
          <WordForm initWordData={wordData} initBookId={bookId} />
        </section>
        <section id='grammar-form' className="word-form-section">
          <h3>添加/编辑语法</h3>
          <GrammarForm />
        </section>
        <section id='reading-form' className="word-form-section">
          <h3>添加/编辑阅读</h3>
          <ReadingForm />
        </section>
      </main>
    </div>
  );
}