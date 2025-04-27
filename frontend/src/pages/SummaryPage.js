import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import './SummaryPage.css';


const SummaryPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const [progress, setProgress] = useState(0);
  const [selectedWord, setSelectedWord] = useState(null);

  // Assume these are passed via state from the previous page
  const reviewedWords = location.state?.reviewedWords || [];
  console.log('Reviewed Words:', reviewedWords);
  const correctness = location.state?.correctness || [];
  const reviewType = location.state?.reviewType || 'word';
  
  // Calculate accuracy
  const correctCount = correctness.filter(Boolean).length;
  const accuracy = reviewedWords.length > 0 ? (correctCount / reviewedWords.length) * 100 : 0;
  const accuracyColor = accuracy >= 80 ? 'green' : accuracy >= 60 ? 'yellow' : 'red';

  useEffect(() => {
    // Animate progress bar
    const duration = 1500; // 1.5 seconds
    const startTime = Date.now();
    
    const animate = () => {
      const elapsed = Date.now() - startTime;
      const progress = Math.min(elapsed / duration, 1);
      setProgress(progress * accuracy);
      
      if (progress < 1) {
        requestAnimationFrame(animate);
      }
    };
    
    animate();
  }, [accuracy]);
  if (!reviewedWords || reviewedWords.length === 0 || !correctness) {
    navigate('/review');
    return null; 
  }
  const handleWordClick = (index) => {
    setSelectedWord(reviewedWords[index]);
  };

  const handleBackToReview = () => {
    navigate('/review');
  };
  const strClamp = (str, maxLength) => {
    if (str.length > maxLength) {
      return str.slice(0, maxLength) + '...';
    }
    return str;
  }
  return (
    <div className="summary-page">
      <h1>复习总结</h1>
      
      <div className="accuracy-section">
        <h2>正确率</h2>
        <div className="progress-bar-container">
          <div 
            className={`progress-bar ${accuracyColor}`}
            style={{ width: `${progress}%` }}
          ></div>
          <div className="progress-text">{Math.round(progress)}%</div>
        </div>
        <p className="accuracy-message">
          在{reviewedWords.length}个{reviewType === 'word' ? '单词' : '语法'}中你答对了{correctCount}个!
        </p>
      </div>
      
      <div className="words-grid">
        {reviewedWords.map((word, index) => (
          <div
            key={word.id}
            className={`word-card ${correctness[index] ? 'correct' : 'incorrect'}`}
            onClick={() => handleWordClick(index)}
          >
            {reviewType === 'word' ? word.kanji || word.hiragana || word.katakana
            : strClamp(word.description, 8)}
          </div>
        ))}
      </div>
      
      {reviewType === 'word' ? (selectedWord && (
        <div className="word-details">
          <h3>单词详情</h3>
          <div className="details-table">
            <div className="detail-row">
              <span className="detail-label">中文:</span>
              <span className="detail-value">{selectedWord.chinese}</span>
            </div>
            {selectedWord.kanji && (
              <div className="detail-row">
                <span className="detail-label">汉字:</span>
                <span className="detail-value">{selectedWord.kanji}</span>
              </div>
            )}
            <div className="detail-row">
              <span className="detail-label">平假名:</span>
              <span className="detail-value">{selectedWord.hiragana}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">片假名:</span>
              <span className="detail-value">{selectedWord.katakana}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">词性:</span>
              <span className="detail-value">{selectedWord.type}</span>
            </div>
            <div className="detail-row">
              <span className="detail-label">你的回答:</span>
              <span className="detail-value">{selectedWord.userAnswer}</span>
            </div>
            {selectedWord.example && selectedWord.example.length > 0 && (
              <div className="examples-section">
                <h4>例句:</h4>
                {selectedWord.example.map((ex, idx) => (
                  <div key={idx} className="example">
                    <div className="japanese">{ex.example}</div>
                    <div className="translation">{ex.chinese}</div>
                  </div>
                ))}
              </div>
            )}
          </div>
          <button 
            className="close-details"
            onClick={() => setSelectedWord(null)}
          >
            关闭
          </button>
        </div>
      )) : (selectedWord && (
        <div className="word-details">
          <h3>语法详情</h3>
          <div className="details-table">
            <div className="detail-row">
              <span className="detail-label">描述:</span>
              <span className="detail-value" style={{ whiteSpace: 'pre-wrap', textAlign: 'left' }}>
                {selectedWord.description.replace(/\\n/g, '\n')}
              </span>
            </div>
            {selectedWord.example && selectedWord.example.length > 0 && (
              <div className="examples-section">
                <h4>例句:</h4>
                {selectedWord.example.map((ex, idx) => (
                  <div key={idx} className="example">
                    <div className="japanese">{ex.example}</div>
                    <div className="translation">{ex.chinese}</div>
                  </div>
                ))}
              </div>
            )}
          </div>
          <button 
            className="close-details"
            onClick={() => setSelectedWord(null)}
          >
            关闭
          </button>
        </div>
      ))}
      
      <button className="back-button" onClick={handleBackToReview}>
        返回复习页面
      </button>
    </div>
  );
};

export default SummaryPage;