import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import './Navigation.css';
import AuthPopup from './Auth/AuthPopup';
import Notification from './Auth/Notification';
import { API_BASE_URL } from '../services/api';

export default function Navigation() {
  const [showAuthPopup, setShowAuthPopup] = useState(false);
  const [notification, setNotification] = useState({ show: false, message: '', type: '' });

  const handleAuthSubmit = async (username, apiKey) => {
    try {
      const response = await fetch(`${API_BASE_URL}/api/user/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': apiKey
        },
        body: JSON.stringify({ username: username })
      });
      
      if (response.ok) {
        setNotification({ 
          show: true, 
          message: '成功提交', 
          type: 'success' 
        });
        // Save to localStorage for persistence
        localStorage.setItem('username', username);
        localStorage.setItem('apiKey', apiKey);
        setShowAuthPopup(false);
      } else {
        const errorData = await response.json();
        setNotification({ 
          show: true, 
          message: `提交失败: ${errorData.error || '未知错误'}`, 
          type: 'error' 
        });
        console.log('Error:', errorData.error || '未知错误');
      }
    } catch (error) {
      setNotification({ 
        show: true, 
        message: `提交失败: ${error.message}`, 
        type: 'error' 
      });
    }
    
    // Auto-hide notification after 3 seconds
    setTimeout(() => {
      setNotification({ show: false, message: '', type: '' });
    }, 3000);
  };

  // Get username from localStorage (if available)
  const savedUsername = localStorage.getItem('username') || '用户名';

  return (
    <>
      <nav className="main-nav">
        <ul>
          <li><Link to="/word-editor"> 修改词库</Link></li>
          <li><Link to="/word-search"> 词库搜索</Link></li>
          <li><Link to="/"> 返回主页</Link></li>
          <li className="right-item">
            <button 
              className="username-btn" 
              onClick={() => setShowAuthPopup(true)}
            >
              {savedUsername}
            </button>
          </li>
        </ul>
      </nav>

      {showAuthPopup && (
        <AuthPopup 
          onClose={() => setShowAuthPopup(false)}
          onSubmit={handleAuthSubmit}
        />
      )}

      {notification.show && (
        <Notification 
          message={notification.message} 
          type={notification.type} 
        />
      )}
    </>
  );
}