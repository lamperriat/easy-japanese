import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './ReviewPage.css';

const ReviewPage = () => {
  const [reviewType, setReviewType] = useState('word');
  const [batchSize, setBatchSize] = useState(10);
  const navigate = useNavigate();

  const handleStartReview = () => {
    navigate(`/review/session?type=${reviewType}&batch=${batchSize}`);
  };

  return (
    <div className="review-page">
      <h1>Review Settings</h1>
      <div className="review-options">
        <div className="option-group">
          <h2>Review Type</h2>
          <div className="radio-group">
            <label>
              <input
                type="radio"
                value="word"
                checked={reviewType === 'word'}
                onChange={() => setReviewType('word')}
              />
              Words
            </label>
            <label>
              <input
                type="radio"
                value="grammar"
                checked={reviewType === 'grammar'}
                onChange={() => setReviewType('grammar')}
              />
              Grammar
            </label>
          </div>
        </div>

        <div className="option-group">
          <h2>Batch Size</h2>
          <select
            value={batchSize}
            onChange={(e) => setBatchSize(parseInt(e.target.value))}
          >
            <option value="5">5</option>
            <option value="10">10</option>
            <option value="20">20</option>
            <option value="30">30</option>
            <option value="50">50</option>
          </select>
        </div>

        <button className="start-review-btn" onClick={handleStartReview}>
          Start Review
        </button>
      </div>
    </div>
  );
};

export default ReviewPage;