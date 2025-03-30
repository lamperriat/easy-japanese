import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../services/api'; 
import { useNavigate } from 'react-router-dom';
import './WordSearchPage.css';
const searchOptions = [
  {
    id: '1', name: '单词', options: [
      { id: '1', name: '新标日初级上' },
      { id: '2', name: '新标日初级下' },
      { id: '3', name: '新标日中级上' },
      { id: '4', name: '新标日中级下' },
      { id: '5', name: '新标日高级上' },
      { id: '6', name: '新标日高级下' },
      { id: '-1', name: 'global'},
    ]
  },
  {
    id: '2', name: '语法', options: [
      { id: '1', name: 'user' }, 
      { id: '2', name: 'global' },
  ] },
  {
    id: '3', name: '阅读', options: [
      { id: '1', name: 'user' }, 
      { id: '2', name: 'global' },
  ] },
]
const bookOptions = [
  { id: '1', name: '新标日初级上' },
  { id: '2', name: '新标日初级下' },
  { id: '3', name: '新标日中级上' },
  { id: '4', name: '新标日中级下' },
  { id: '5', name: '新标日高级上' },
  { id: '6', name: '新标日高级下' },
  { id: '-1', name: 'global'}
];

const WordSearchPage = () => {
  const navigate = useNavigate();
  const [words, setWords] = useState([]);
  const [filteredWords, setFilteredWords] = useState([]);
  const [searchQuery, setSearchQuery] = useState(''); 
  const [searchType, setSearchType] = useState('1'); // word=1, grammar=2, reading=3
  const [selectedBook, setSelectedBook] = useState(bookOptions[0].id); 
  const [isLoading, setIsLoading] = useState(false); 
  const [apiMessage, setApiMessage] = useState('');

  const handleWordClick = (word) => {
    navigate('/word-editor#word-form', {
      state: { word: word, selectedBook: selectedBook }
    });
  };

  const fetchWords = async (endpoint, method = 'GET', body = null) => {
    setIsLoading(true);
    setApiMessage('');
    try {
      var token = sessionStorage.getItem('token');
      if (!token) {
        setApiMessage('请先登录');
        setIsLoading(false);
        return [];
      }
      const response = await fetch(endpoint, {
        method,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token,
        },
        body: body ? JSON.stringify(body) : null,
      });
  
      const responseText = await response.text(); 
      console.log("Raw API Response:", responseText); 
  
      try {
        const result = JSON.parse(responseText).results; 
        console.log("Parsed API Response:", result); 
  
        if (response.ok) {
          return Array.isArray(result) ? result : result.words || result.data || [];
        } else {
          setApiMessage(result.error || '获取词表失败');
          return [];
        }
      } catch (error) {
        console.error("JSON Parsing Error:", error);
         setApiMessage('解析API响应失败: 非法的 JSON 格式');
        return [];
      }
  
    } catch (error) {
      console.error("Fetch Error:", error);
      setApiMessage('网络请求失败: ' + error.message);
      return [];
    } finally {
      setIsLoading(false);
    }
  };
  

  const fetchWordList = async () => {
    const endpoint = `${API_BASE_URL}/api/words/book_${selectedBook}/get`;
    const words = await fetchWords(endpoint);
    console.log("Fetched Word List:", words); 
    setWords(words); 
  };

  const fetchSimilarWords = async (query) => {
    var token = sessionStorage.getItem('token');
    if (!token) {
      setApiMessage('请先登录');
      setIsLoading(false);
      return [];
    }
    const endpoint = `${API_BASE_URL}/api/words/book_${selectedBook}/fuzzy-search?query=${encodeURIComponent(query)}`;
    const response = await fetch(endpoint, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': token,
      },
    });
  
    const responseData = await response.json();
    console.log("Fetched Response:", responseData);  
  
    if (responseData.results && responseData.results.length > 0) {
      setFilteredWords(responseData.results);
    } else {
      setFilteredWords([]); 
    }
  };
  

  useEffect(() => {
    if (selectedBook) {
      fetchWordList(); 
    }
  }, [selectedBook]);

  useEffect(() => {
    setFilteredWords(words); 
    console.log("Filtered Words Updated:", words); 
  }, [words]);

  useEffect(() => {
    console.log("Filtered Words Updated:", filteredWords); 
  }, [filteredWords]);
  

  const handleSearch = (query) => {
    setSearchQuery(query);
    if (!query.trim()) {
      setFilteredWords(words); 
      setApiMessage('');  
    } else {
      fetchSimilarWords(query); 
    }
  };

  const [searchCategory, setSearchCategory] = useState('1'); 
  
  const currentOptions = searchOptions[searchCategory - 1].options || [];

  useEffect(() => {
    if (currentOptions.length > 0) {
      setSelectedBook(currentOptions[0].id);
    }
  }, [searchCategory]);

  return (
    <div className="word-search-page">
      <header className="main-header">
        <h1>词库搜索</h1>
        <h2>Word Search</h2>
      </header>

      <main className="search-content">
        <section className="dictionary-section">
          <label htmlFor="category">选择类别:</label>
          <select
            id="category"
            value={searchCategory}
            onChange={(e) => setSearchCategory(e.target.value)}
          >
            {searchOptions.map((option) => (
              <option key={option.id} value={option.id}>
                {option.name}
              </option>
            ))}
          </select>
        </section>

        <section className="dictionary-section">
          <label htmlFor="dictionary">选择词典:</label>
          <select
            id="dictionary"
            value={selectedBook}
            onChange={(e) => setSelectedBook(e.target.value)}
          >
            {currentOptions.map((option) => (
              <option key={option.id} value={option.id}>
                {option.name}
              </option>
            ))}
          </select>
        </section>

        <section className="search-section">
          <input
            type="text"
            placeholder="搜索单词..."
            value={searchQuery}
            onChange={(e) => handleSearch(e.target.value)}
          />
        </section>

        {isLoading && <p>加载中...</p>}
        {apiMessage && <p className="error-message">{apiMessage}</p>}

        <section className="word-list-section">
          <h3>单词列表</h3>
          <ul>
            {filteredWords.length > 0 ? (
              filteredWords.map((word, index) => (
                <li key={index} onClick={() => handleWordClick(word)}>
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
