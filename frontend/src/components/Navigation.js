import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import './Navigation.css';
import AuthPopup from './Auth/AuthPopup';
import Notification from './Auth/Notification';
import { API_BASE_URL } from '../services/api';

export default function Navigation() {
  const [showAuthPopup, setShowAuthPopup] = useState(false);
  const [notification, setNotification] = useState({ show: false, message: '', type: '' });

  const handleAuthSubmit = async (username, apikey) => {
    try {
      var token = sessionStorage.getItem('token');
      if (!token) {
        const response = await fetch(`${API_BASE_URL}/api/auth/token`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-API-Key': apikey
          },
          body: JSON.stringify({ username: username })
        });
        if (!response.ok) {
          const errorData = await response.json();
          setNotification({
            show: true,
            message: `提交失败: ${errorData.error || '未知错误'}`,
            type: 'error'
          });
        }
        const data = await response.json();
        token = data.token;
        sessionStorage.setItem('token', token);
      }
      const response = await fetch(`${API_BASE_URL}/api/user/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token
        },
        body: JSON.stringify({ username: username })
      });
      
      if (response.ok) {
        setNotification({
          show: true,
          message: '成功提交',
          type: 'success'
        });
        localStorage.setItem('username', username);
        setShowAuthPopup(false);
      } else if (response.status === 409) {
        setNotification({ 
          show: true, 
          message: '该API已经注册', 
          type: 'success' 
        });
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