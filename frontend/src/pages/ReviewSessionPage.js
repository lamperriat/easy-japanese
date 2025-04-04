import React, { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { API_BASE_URL } from '../services/api';
import './ReviewSessionPage.css';

const ReviewSessionPage = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const queryParams = new URLSearchParams(location.search);
  const reviewType = queryParams.get('type');
  const batchSize = queryParams.get('batch');

  const [items, setItems] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [displayMode, setDisplayMode] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchReviewItems = async () => {
      try {
        setIsLoading(true);
        const response = await fetch(
          `${API_BASE_URL}/api/user/review/get/${reviewType}?batch=${batchSize}`
        );
        if (!response.ok) {
          throw new Error('Failed to fetch review items');
        }
        const data = await response.json();
        setItems(data);
        
        // Set initial display mode for words
        if (reviewType === 'word' && data.length > 0) {
          setRandomDisplayMode();
        }
        
        setIsLoading(false);
      } catch (err) {
        setError(err.message);
        setIsLoading(false);
      }
    };

    fetchReviewItems();
  }, [reviewType, batchSize]);

  const setRandomDisplayMode = () => {
    const random = Math.random();
    if (random < 0.3) {
      setDisplayMode('chinese');
    } else if (random < 0.6) {
      setDisplayMode('kanji');
    } else {
      setDisplayMode('kana');
    }
  };

  const handleNext = () => {
    if (currentIndex < items.length - 1) {
      setCurrentIndex(currentIndex + 1);
      if (reviewType === 'word') {
        setRandomDisplayMode();
      }
    } else {
      navigate('/review');
    }
  };

  const handlePrevious = () => {
    if (currentIndex > 0) {
      setCurrentIndex(currentIndex - 1);
    }
  };

  const handleAnswer = (isCorrect) => {
    // Here you would typically send the result to your API
    // For now, we'll just move to the next item
    handleNext();
  };

  const renderWord = () => {
    const word = items[currentIndex];
    
    switch (displayMode) {
      case 'chinese':
        return (
          <div className="word-display">
            <div className="main-display">{word.chinese}</div>
            <div className="hint">Chinese meaning</div>
          </div>
        );
      case 'kanji':
        return word.kanji ? (
          <div className="word-display">
            <div className="main-display">{word.kanji}</div>
            <div className="hint">Kanji</div>
          </div>
        ) : (
          <div className="word-display">
            <div className="main-display">{word.hiragana || word.katakana}</div>
            <div className="hint">Kana</div>
          </div>
        );
      case 'kana':
        return (
          <div className="word-display">
            <div className="main-display">{word.hiragana || word.katakana}</div>
            <div className="hint">Kana</div>
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
          <h3>Examples:</h3>
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

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  if (items.length === 0) {
    return (
      <div className="no-items">
        <p>No items to review.</p>
        <button onClick={() => navigate('/review')}>Back to Review Settings</button>
      </div>
    );
  }

  return (
    <div className="review-session">
      <div className="progress">
        {currentIndex + 1} / {items.length}
      </div>

      <div className="content">
        {reviewType === 'word' ? renderWord() : renderGrammar()}
      </div>

      <div className="navigation">
        <button className="nav-btn prev" onClick={handlePrevious}>
          Previous
        </button>
        <div className="answer-buttons">
          <button className="answer-btn incorrect" onClick={() => handleAnswer(false)}>
            Incorrect
          </button>
          <button className="answer-btn correct" onClick={() => handleAnswer(true)}>
            Correct
          </button>
        </div>
        <button className="nav-btn next" onClick={handleNext}>
          Next
        </button>
      </div>
    </div>
  );
};

export default ReviewSessionPage;