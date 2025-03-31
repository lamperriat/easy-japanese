import React, { useState, useEffect } from 'react';
import './AuthPopup.css';
import Notification from './Notification';
const AuthPopup = ({ onClose, onSubmit }) => {
  const [username, setUsername] = useState('');
  const [apiKey, setApiKey] = useState('');
  const [notification, setNotification] = useState({ show: false, message: '', type: '' });
  
  // Load from localStorage if available
  useEffect(() => {
    const savedUsername = localStorage.getItem('username');
    
    if (savedUsername) setUsername(savedUsername);
  }, []);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (username && apiKey) {
      onSubmit(username, apiKey);
    }
  };

  return (
    <div className="auth-popup-overlay">
      <div className="auth-popup">
        <button className="close-button" onClick={onClose}>×</button>
        <h2>用户登录</h2>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="username">用户名</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="请输入用户名"
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="apiKey">API Key</label>
            <input
              type="password"
              id="apiKey"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder="请输入API Key"
              required
            />
          </div>
          <div className="button-group">
            <button 
              type="button" 
              className="logout-button" 
              onClick={() => {
                sessionStorage.removeItem('token');
                setNotification({
                  show: true,
                  message: '已登出',
                  type: 'success'
                });
                setTimeout(() => {
                  setNotification({ show: false, message: '', type: '' });
                }, 3000);
              }}
            >
              登出
            </button>
            <button type="submit" className="submit-button">提交</button>
          </div>
        </form>
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

export default AuthPopup;