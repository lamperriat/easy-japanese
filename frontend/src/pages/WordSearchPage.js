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


  const fetchWords = async (endpoint, method = 'GET', body = null) => {
    setIsLoading(true);
    setApiMessage('');
    try {
      const response = await fetch(endpoint, {
        method,
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': process.env.REACT_APP_API_KEY,
        },
        body: body ? JSON.stringify(body) : null,
      });
  
      const responseText = await response.text(); // 获取原始的响应文本
      console.log("Raw API Response:", responseText); // 打印原始响应
  
      try {
        const result = JSON.parse(responseText); // 尝试将响应文本解析为 JSON
        console.log("Parsed API Response:", result); // 打印解析后的 JSON
  
        if (response.ok) {
          return Array.isArray(result) ? result : result.words || result.data || [];
        } else {
          setApiMessage(result.error || '获取词表失败');
          return [];
        }
      } catch (error) {
        console.error("JSON Parsing Error:", error); // 打印 JSON 解析错误
         setApiMessage('解析API响应失败: 非法的 JSON 格式');
        return [];
      }
  
    } catch (error) {
      console.error("Fetch Error:", error); // 打印 fetch 请求错误
      setApiMessage('网络请求失败: ' + error.message);
      return [];
    } finally {
      setIsLoading(false);
    }
  };
  

  // 获取词库数据
  const fetchWordList = async () => {
    const endpoint = `${API_BASE_URL}/api/words/book_${selectedBook}/get`;
    const words = await fetchWords(endpoint);
    console.log("Fetched Word List:", words); // 打印获取的词库
    setWords(words); // 更新词库数据
  };

  // 获取相似单词
  const fetchSimilarWords = async (query) => {
    const endpoint = `${API_BASE_URL}/api/words/book_${selectedBook}/fuzzy-search?query=${encodeURIComponent(query)}`;
    const response = await fetch(endpoint, {
      method: 'GET',  // 保持 GET 请求
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': process.env.REACT_APP_API_KEY,
      },
    });
  
    const responseData = await response.json();
    console.log("Fetched Response:", responseData);  // 打印实际返回的数据
  
    // 检查返回的数据是否有 similar 字段
    if (responseData.results && responseData.results.length > 0) {
      setFilteredWords(responseData.results);
    } else {
      setFilteredWords([]); // 如果没有数据，清空列表
    }
  };
  

  // 监听教材选择变化，获取词库数据
  useEffect(() => {
    if (selectedBook) {
      fetchWordList(); // 选择教材时，加载所有单词
    }
  }, [selectedBook]);

  // 监听 words 数据变化，更新 filteredWords
  useEffect(() => {
    setFilteredWords(words); // 每当 words 更新时，更新 filteredWords
    console.log("Filtered Words Updated:", words); // 打印更新后的filteredWords
  }, [words]);

  useEffect(() => {
    console.log("Filtered Words Updated:", filteredWords); // 打印 updated filteredWords
  }, [filteredWords]);
  

  // 处理搜索功能
  const handleSearch = (query) => {
    setSearchQuery(query);
    if (!query.trim()) {
      // 如果没有输入，显示所有单词
      setFilteredWords(words); // 恢复到所有的单词列表
      setApiMessage('');  // 清空消息
    } else {
      fetchSimilarWords(query); // 有输入时，发起请求获取相似单词
    }
  };

  return (
    <div className="word-search-page">
      <header className="main-header">
        <h1>词库搜索</h1>
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
              filteredWords.map((word, index) => (
                <li key={index}>
                  <div>Kanji: {word.kanji || '无'}</div>
                  <div>Chinese: {word.chinese || '无'}</div>
                  <div>Hiragana: {word.hiragana || '无'}</div>
                  <div>Katakana: {word.katakana || '无'}</div>
                </li>
              ))
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
