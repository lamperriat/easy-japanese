import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { API_BASE_URL } from '../services/api';
import Notification from '../components/Auth/Notification';
import './ReviewSessionPage.css';

const ReviewSessionPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const queryParams = new URLSearchParams(location.search);
  const reviewType = queryParams.get('type');
  const batchSize = queryParams.get('batch');

  const [items, setItems] = useState([]);
  const [correctness, setCorrectness] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [displayMode, setDisplayMode] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [notification, setNotification] = useState({ show: false, message: '', type: '' });
  // const [error, setError] = useState(null);
  const [userAnswer, setUserAnswer] = useState('');
  const [showFullInfo, setShowFullInfo] = useState(false);
  const [answered, setAnswered] = useState(false);
  useEffect(() => {
    const fetchReviewItems = async () => {
      setIsLoading(true);
      var token = sessionStorage.getItem('token');
      if (!token) {
        setNotification({
          show: true,
          message: '请先登录',
          type: 'error'
        });
        setTimeout(() => {
          setNotification({ show: false, message: '', type: '' });
        }, 3000);
        setIsLoading(false);
        return;
      }
      const response = await fetch(
        `${API_BASE_URL}/api/user/review/${reviewType}/get?batch=${batchSize}`, 
        {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': token
          }
        }
      );
      if (response.status === 401) {
        setNotification({
          show: true,
          message: '登录已过期，请重新登录',
          type: 'error'
        });
        setTimeout(() => {
          setNotification({ show: false, message: '', type: '' });
        }, 3000);
        setIsLoading(false);
        return;
      }
      const data = await response.json();
      setItems(data);
      setCorrectness(new Array(data.length).fill(false));
      setIsLoading(false);
    };

    fetchReviewItems();
  }, [reviewType, batchSize]);

  const clear = () => {
    setUserAnswer('');
    setAnswered(false);
    setShowFullInfo(false);
  }

  const setRandomDisplayMode = () => {
    if (currentIndex < 0 || currentIndex >= items.length) return;
    const random = Math.random();
    if (random < 0.3) {
      setDisplayMode('chinese');
    } else {
      if (items[currentIndex].kanji) {
        if (random < 0.65) {
          setDisplayMode('kanji');
        } else {
          setDisplayMode('kana');
        }
      } else {
        setDisplayMode('kana');
      }
    }
  };

  useEffect(() => {
    if (items.length > 0 && reviewType === 'word' && currentIndex < items.length) {
      setRandomDisplayMode();
    }
  }, [items, reviewType]);

  const handlePrevious = () => {
    if (currentIndex > 0) {
      setCurrentIndex(currentIndex - 1);
      if (reviewType === 'word') {
        setRandomDisplayMode();
      }
      clear();
    }
  };

  const handleAnswer = (isCorrect) => {
    setShowFullInfo(true);
    setCorrectness((prev) => {
      const newCorrectness = [...prev];
      newCorrectness[currentIndex] = isCorrect;
      return newCorrectness;
    });
    setItems((prev) => {
      const newItems = [...prev];
      newItems[currentIndex].userAnswer = userAnswer;
      return newItems;
    });
    setAnswered(true);
    // for grammars, we directly go next
    if (reviewType === 'grammar') {
      handleNext();
    }
  };

  const handleNext = () => {
    if (currentIndex < items.length - 1) {
      setCurrentIndex(currentIndex + 1);
      if (reviewType === 'word') {
        setRandomDisplayMode();
      }
      clear();
    } else {
      // TODO: Handle end of review session
      navigate('/summary', {
        state: {
          reviewedWords: items,
          correctness: correctness,
          reviewType: reviewType,
        }
      });
    }
  };

  const renderWord = () => {
    const word = items[currentIndex];
    if (showFullInfo) {
      return (
        <div className="word-full-info">
          <div className="info-row">
            <span className="label">中文:</span>
            <span className="value">{word.chinese}</span>
          </div>
          {word.kanji && (
            <div className="info-row">
              <span className="label">汉字:</span>
              <span className="value">{word.kanji}</span>
            </div>
          )}
          <div className="info-row">
            <span className="label">平假名:</span>
            <span className="value">{word.hiragana}</span>
          </div>
          <div className="info-row">
            <span className="label">片假名:</span>
            <span className="value">{word.katakana}</span>
          </div>
          <div className="info-row">
            <span className="label">词性:</span>
            <span className="value">{word.type}</span>
          </div>
          <div className="info-row">
            <span className="label">你的回答:</span>
            <span className="value">{userAnswer}</span>
          </div>
          {word.example && word.example.length > 0 && (
            <div className="examples">
              <h3>例句:</h3>
              {word.example.map((ex, idx) => (
                <div key={idx} className="example">
                  <div className="japanese">{ex.example}</div>
                  <div className="translation">{ex.chinese}</div>
                </div>
              ))}
            </div>
          )}
        </div>
      );
    }
    switch (displayMode) {
      case 'chinese':
        return (
          <div className="word-display">
            <div className="main-display">{word.chinese}</div>
            <div className="hint">对应的汉字/假名是？</div>
          </div>
        );
      case 'kanji':
        return word.kanji ? (
          <div className="word-display">
            <div className="main-display">{word.kanji}</div>
            <div className="hint">对应的中文/假名是？</div>
          </div>
        ) : (
          <div className="word-display">
            <div className="main-display">{word.hiragana || word.katakana}</div>
            <div className="hint">对应的中文/汉字是？</div>
          </div>
        );
      case 'kana':
        return (
          <div className="word-display">
            <div className="main-display">{word.hiragana || word.katakana}</div>
            <div className="hint">对应的中文/汉字是？</div>
          </div>
        );
      default:
        return null;
    }
  };

  const renderGrammar = () => {
    const grammar = items[currentIndex];
    return (
      <div className="grammar-display">
        <div className="description">{grammar.description}</div>
        <div className="examples">
          <h3>例句:</h3>
          {grammar.example.map((ex, idx) => (
            <div key={idx} className="example">
              <div className="japanese">{ex.example}</div>
              <div className="translation">{ex.chinese}</div>
            </div>
          ))}
        </div>
      </div>
    );
  };

  if (isLoading) {
    return <div className="loading">Loading...</div>;
  }

  // if (error) {
  //   return <div className="error">Error: {error}</div>;
  // }

  if (items.length === 0) {
    return (
      <div className="no-items">
        <p>No items to review.</p>
        <button onClick={() => navigate('/review')}>Back to Review Settings</button>
        {notification.show && (
          <Notification 
          message={notification.message} 
          type={notification.type} 
          />
        )}
      </div>
    );
  }
  return (
    <div className="review-session">
      <div className="header">
        <button className="exit-btn" onClick={() => navigate('/review')}>
          退出复习 
        </button>
        <div className="progress">
          {currentIndex + 1} / {items.length}
        </div>
      </div>

      <div className="content">
        {reviewType === 'word' ? renderWord() : renderGrammar()}
      </div>

      {!showFullInfo && (
        <div className="user-answer">
          <textarea
            value={userAnswer}
            onChange={(e) => setUserAnswer(e.target.value)}
            onKeyDown={(e) => {
              if (reviewType === 'word' && e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault(); 
                setShowFullInfo(true);
              }
            }}
            placeholder={reviewType == 'word' ? "写下你的回答(非必须)\n按下Enter显示正确答案" :
              "将你的联想记录在这(不会保存至后台)"}
            rows="2"
          />
        </div>
      )}

      <div className="navigation">
        <button className="nav-btn prev" onClick={handlePrevious} disabled={currentIndex === 0}>
          上一个
        </button>
        
        {!answered ? (
          <div className="answer-buttons">
            <button className="answer-btn incorrect" onClick={() => handleAnswer(false)}>
              我答错了
            </button>
            <button className="answer-btn correct" onClick={() => handleAnswer(true)}>
              我答对了
            </button>
          </div>
        ) : (
          <button className="nav-btn next" onClick={handleNext}>
            {currentIndex < items.length - 1 ? '下一个' : '完成'}
          </button>
        )}
      </div>
      {notification.show && (
        <Notification 
        message={notification.message} 
        type={notification.type} 
        />
      )}
    </div>
  );
};


export default ReviewSessionPage;