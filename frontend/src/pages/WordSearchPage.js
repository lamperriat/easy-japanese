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
      { id: '-1', name: '个人'},
    ]
  },
  {
    id: '2', name: '语法', options: [
      { id: '-1', name: '个人' }, 
      { id: '1', name: '全局' },
  ] },
  {
    id: '3', name: '阅读', options: [
      { id: '-1', name: '个人' }, 
      { id: '1', name: '全局' },
  ] },
]

const WordSearchPage = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const navigate = useNavigate();
  const [words, setWords] = useState([]);
  const [grammars, setGrammars] = useState([]);
  const [readings, setReadings] = useState([]);
  const [filteredWords, setFilteredWords] = useState([]);
  const [filteredGrammars, setFilteredGrammars] = useState([]);
  const [filteredReadings, setFilteredReadings] = useState([]);
  const [searchQuery, setSearchQuery] = useState(''); 
  const [selectedBook, setSelectedBook] = useState(searchOptions[0].options[0].id); 
  const [isLoading, setIsLoading] = useState(false); 
  const [apiMessage, setApiMessage] = useState('');
  const [searchType, setSearchType] = useState('1'); 
  const resultPerPage = 30;
  const handleWordClick = (word) => {
    navigate('/word-editor#word-form', {
      state: { word: word, selectedBook: selectedBook }
    });
  };
  const handleGrammarClick = (grammar) => {
    navigate('/word-editor#grammar-form', {
      state: { grammar: grammar, selectedBook: selectedBook }
    });
  };
  const handleReadingClick = (reading) => {
    navigate('/word-editor#reading-form', {
      state: { reading: reading, selectedBook: selectedBook }
    });
  };

  const fetchRemote = async (endpoint, method = 'GET', body = null) => {
    setIsLoading(true);
    setApiMessage('');
    try {
      var token = sessionStorage.getItem('token');
      // no token or token expired
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
      if (response.status === 401) {
        setApiMessage('登录已过期，请重新登录');
        sessionStorage.removeItem('token');
        setIsLoading(false);
        return [];
      }
      const responseText = await response.text(); 
      console.log("Raw API Response:", responseText); 
      
      if (response.status === 404) {
        return [];
      }
      try {
        const parsed = JSON.parse(responseText);
        if (parsed.count) {
          setTotalCount(parsed.count);
        }
        const result = Array.isArray(parsed) ? parsed : parsed.results; 
        console.log("Parsed API Response:", result); 
        
        if (response.ok) {
          return Array.isArray(result) ? result : [];
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
    const endpoint = selectedBook === '-1' ?
      `${API_BASE_URL}/api/user/words/get?page=${currentPage}`:
      `${API_BASE_URL}/api/words/book_${selectedBook}/get?page=${currentPage}`;
    const words = await fetchRemote(endpoint);
    setWords(words); 
  };
  const fetchGrammarList = async () => {
    const endpoint = selectedBook === '-1' ? 
      `${API_BASE_URL}/api/user/grammar/get?page=${currentPage}` :
      `${API_BASE_URL}/api/grammar/get?page=${currentPage}`;
    const grammars = await fetchRemote(endpoint);
    setGrammars(grammars);
  };
  const fetchReadingList = async () => {
    const endpoint = selectedBook === '-1' ?
      `${API_BASE_URL}/api/user/reading-material/get?page=${currentPage}` :
      `${API_BASE_URL}/api/reading-material/get?page=${currentPage}`;
    const readings = await fetchRemote(endpoint);
    setReadings(readings);
  };

  const fetchSimilar = async (query) => {
    var token = sessionStorage.getItem('token');
    if (!token) {
      setApiMessage('请先登录');
      setIsLoading(false);
      return [];
    }
    var endpoint = '';
    if (searchType === '1') {
      endpoint = selectedBook !== "-1" ? `${API_BASE_URL}/api/words/book_${selectedBook}/fuzzy-search?query=${encodeURIComponent(query)}&page=${currentPage}`
          : `${API_BASE_URL}/api/user/words/fuzzy-search?query=${encodeURIComponent(query)}&page=${currentPage}`;
    } else if (searchType === '2') {
      if (selectedBook === '-1') {
        endpoint = `${API_BASE_URL}/api/user/grammar/search?query=${encodeURIComponent(query)}&page=${currentPage}`;
      } else {
        endpoint = `${API_BASE_URL}/api/grammar/search?query=${encodeURIComponent(query)}&page=${currentPage}`;
      }
    } else if (searchType === '3') {
      if (selectedBook === '-1') {
        endpoint = `${API_BASE_URL}/api/user/reading-material/search?query=${encodeURIComponent(query)}&page=${currentPage}`;
      } else {
        endpoint = `${API_BASE_URL}/api/reading-material/search?query=${encodeURIComponent(query)}&page=${currentPage}`;
      }
    } else {
      return [];
    }
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
      if (searchType === '1') {
        setFilteredWords(responseData.results);
      } else if (searchType === '2') {
        setFilteredGrammars(responseData.results);
      } else if (searchType === '3') {
        setFilteredReadings(responseData.results);
      }
    } else {
      if (searchType === '1') {
        setFilteredWords([]);
      } else if (searchType === '2') {
        setFilteredGrammars([]);
      } else if (searchType === '3') {
        setFilteredReadings([]);
      }
      setApiMessage('没有找到匹配的结果');
    }
  };
  

  useEffect(() => {
    if (searchType === '1') {
      fetchWordList(); 
    } else if (searchType === '2') {
      fetchGrammarList(); 
    } else if (searchType === '3') {
      fetchReadingList(); 
    }
  }, [searchType, selectedBook, currentPage]);

  useEffect(() => {
    setFilteredWords(words); 
  }, [words]);
  useEffect(() => {
    setFilteredGrammars(grammars);
  }, [grammars]);
  useEffect(() => {
    setFilteredReadings(readings);
  }, [readings]);

  const handleSearch = (query) => {
    setSearchQuery(query);
    if (!query.trim()) {
      if (searchType === '1') {
        setFilteredWords(words); 
      } else if (searchType === '2') {
        setFilteredGrammars(grammars); 
      } else if (searchType === '3') {
        setFilteredReadings(readings); 
      }
      setApiMessage('');  
    } else {
      fetchSimilar(query); 
    }
  };

  
  
  const currentOptions = searchOptions[searchType - 1].options || [];

  useEffect(() => {
    if (currentOptions.length > 0) {
      setSelectedBook(currentOptions[0].id);
    }
  }, [searchType]);
  const wordResult = 
    filteredWords.length > 0 ? (
      filteredWords.map((word, index) => (
        <li key={index} onClick={() => handleWordClick(word)}>
          <div>汉字: {word.kanji || '无'}</div>
          <div>中文: {word.chinese || '无'}</div>
          <div>平假名: {word.hiragana || '无'}</div>
          <div>片假名: {word.katakana || '无'}</div>
        </li>
      ))
    ) : (
      <li>没有找到匹配的单词。</li>
    );
  const grammarResult =
    filteredGrammars.length > 0 ? (
      filteredGrammars.map((grammar, index) => (
        <li key={index} onClick={() => handleGrammarClick(grammar)}>
          <div>语法: {grammar.description || '无'}</div>
        </li>
      ))
    ) : (
      <li>没有找到匹配的语法。</li>
    );
  const readingResult =
    filteredReadings.length > 0 ? (
      filteredReadings.map((reading, index) => (
        <li key={index} onClick={() => handleReadingClick(reading)}>
          <div>标题: {reading.title || '无'}</div>
          <div>内容: {reading.content || '无'}</div>
        </li>
      ))
    ) : (
      <li>没有找到匹配的阅读材料。</li>
    );
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
            value={searchType}
            onChange={(e) => setSearchType(e.target.value)}
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
          <h3>结果列表</h3>
          <div className="results-count">
            共计{totalCount}条结果
          </div>
          <div className="word-grid">
            {searchType === '1' && wordResult}
            {searchType === '2' && grammarResult}
            {searchType === '3' && readingResult}
          </div>
        </section>
        {totalCount > resultPerPage && (
          <div className='pagination'>
            <button
              disabled={currentPage === 1}
              onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
            >
              上一页
            </button>
            <span>第 {currentPage} 页</span>
            <button
              disabled={currentPage * resultPerPage >= totalCount}
              onClick={() => setCurrentPage(Math.min(Math.ceil(totalCount / resultPerPage), currentPage + 1))}
            >
              下一页
            </button>
          </div>
        )}
      </main>
    </div>
  );
};

export default WordSearchPage;
