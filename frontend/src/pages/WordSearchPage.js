import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../services/api'; // 引入 API 地址
import './WordSearchPage.css';

// 词典选项（教材）
const bookOptions = [
  { id: '1', name: '新标日初级上' },
  { id: '2', name: '新标日初级下' },
  { id: '3', name: '新标日中级上' },
  { id: '4', name: '新标日中级下' },
  { id: '5', name: '新标日高级上' },
  { id: '6', name: '新标日高级下' },
];

const WordSearchPage = () => {
  const [words, setWords] = useState([]); // 词库数据
  const [filteredWords, setFilteredWords] = useState([]); // 过滤后的词
  const [searchQuery, setSearchQuery] = useState(''); // 搜索关键字
  const [selectedBook, setSelectedBook] = useState(bookOptions[0].id); // 默认选择第一本书 (使用 id)
  const [isLoading, setIsLoading] = useState(false); // 控制加载状态
  const [apiMessage, setApiMessage] = useState(''); // 显示 API 消息

  // 获取词库数据
  const fetchWordList = async () => {
    setIsLoading(true);
    setApiMessage('');
    try {
      const endpoint = `${API_BASE_URL}/api/dict/book_${selectedBook}/get`;

      const response = await fetch(endpoint, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': process.env.REACT_APP_API_KEY,
        },
      });

      const result = await response.json();
      if (response.ok) {
        setWords(result.words || []);
        setFilteredWords(result.words || []); // 初始状态下展示所有词汇
      } else {
        setApiMessage(result.error || '获取词表失败');
      }
    } catch (error) {
      setApiMessage('网络请求失败');
    }
    setIsLoading(false);
  };

  // 监听教材选择变化，获取词库数据
  useEffect(() => {
    if (selectedBook) {
      fetchWordList();
    }
  }, [selectedBook]);

  // 处理搜索功能
  const handleSearch = (query) => {
    setSearchQuery(query);
    if (!query.trim()) {
      setFilteredWords(words); // 为空时显示所有词
      return;
    }
    const filtered = words.filter((word) =>
      String(word).toLowerCase().includes(query.toLowerCase())
    );
    setFilteredWords(filtered);
  };

  return (
    <div className="word-search-page">
      <header className="main-header">
        <h1>词表搜索</h1>
        <h2>Word Search</h2>
      </header>

      <main className="search-content">
        {/* 选择词典 */}
        <section className="dictionary-section">
          <label htmlFor="dictionary">选择词典:</label>
          <select
            id="dictionary"
            value={selectedBook}
            onChange={(e) => setSelectedBook(e.target.value)}
          >
            {bookOptions.map((book) => (
              <option key={book.id} value={book.id}>
                {book.name}
              </option>
            ))}
          </select>
        </section>

        {/* 搜索框 */}
        <section className="search-section">
          <input
            type="text"
            placeholder="搜索单词..."
            value={searchQuery}
            onChange={(e) => handleSearch(e.target.value)}
          />
        </section>

        {/* 加载状态 & API 错误信息 */}
        {isLoading && <p>加载中...</p>}
        {apiMessage && <p className="error-message">{apiMessage}</p>}

        {/* 单词列表 */}
        <section className="word-list-section">
          <h3>单词列表</h3>
          <ul>
            {filteredWords.length > 0 ? (
              filteredWords.map((word, index) => <li key={index}>{word}</li>)
            ) : (
              <li>没有找到匹配的单词。</li>
            )}
          </ul>
        </section>
      </main>
    </div>
  );
};

export default WordSearchPage;
