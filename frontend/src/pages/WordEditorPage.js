import React from 'react';
import WordForm from '../components/WordEditor/WordForm';

export default function WordEditorPage() {
  return (
    <div className="word-editor-page">
      <header className="main-header">
        <h1>词库管理</h1>
        <h2>Word Database Management</h2>
      </header>


      <main className="editor-content">
        <section className="word-form-section">
          <h3>添加/编辑单词</h3>
          <WordForm />
        </section>
      </main>
    </div>
  );
}