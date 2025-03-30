import { useEffect, useState } from 'react';
import { API_BASE_URL } from '../../services/api';
import Notification from '../Auth/Notification';

export default function ReadingForm({ initReadingData, initBookId }) {
  const [formData, setFormData] = useState({
    id: 0,
    title: '', 
    content: '', 
    chinese: '',
  });
  const resetForm = () => {
    setFormData({
        id: 0,
        title: '',
        content: '',
        chinese: ''
    });
  }
  const resetFormUser = () => {
    resetForm();
    setApiMessage('表单清空完成');
  }
  const [selectedBook, setSelectedBook] = useState('1');
  const [apiMessage, setApiMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [notification, setNotification] = useState({ show: false, message: '', type: '' });
  const bookOptions = [
    { id: '0', name: 'global' }, 
    { id: '-1', name: 'user' },
  ]
  useEffect(() => {
    if (initBookId) {
      setSelectedBook(initBookId);
    }
    if (initReadingData) {
      setFormData({
        id: initReadingData.id,
        title: initReadingData.title,
        content: initReadingData.content,
        chinese: initReadingData.chinese
      });
      setApiMessage('加载完成。清注意：此时您的提交会修改该词条，而非新增。');
    }
  }, [initReadingData, initBookId]);
  const handleSubmit = async (actionType) => {
    setIsLoading(true);
    try {
      var endpoint = '';
      var method  = '';
      if (actionType === 'check') {
        setApiMessage('由于阅读材料只能全文检索，暂时不提供重复检查');
        return;
      } else if (actionType === 'submit') {
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/reading-material/add`;
          method = 'POST';
        } else {
          endpoint = `${API_BASE_URL}/api/reading-material/add`;
          method = 'POST';
        }
      } else if (actionType === 'delete') {
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/reading-material/delete`;
          method = 'POST';
        } else {
          endpoint = `${API_BASE_URL}/api/reading-material/delete`;
          method = 'POST';
        }
      } else {
        setApiMessage('未知操作');
        return;
      }
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
      const response = await fetch(endpoint, {
        method: method, 
        headers: {
          'Content-Type': 'application/json',
          'Authorization': token,
        },
        body: JSON.stringify(formData)
      });

      const result = await response.json();
      if (response.ok && actionType !== "check") {
        resetForm();
      }
      console.log(result);
      setApiMessage("操作成功");
    } catch (error) {
      setApiMessage("网络请求失败");
    }
    setIsLoading(false);
  };

    return (
      <div className="word-editor">
        <form>
        <div className="form-row">
          <label htmlFor="book-select">选择数据库:</label>
          <select
            id="book-select"
            value={selectedBook}
            onChange={(e) => setSelectedBook(e.target.value)}
            className="form-control"
          >
            {bookOptions.map(book => (
              <option key={book.id} value={book.id}>{book.name}</option>
            ))}
          </select>
          </div>
        <div className="form-group">
          <label>ID</label>
          <input
            value={formData.id}
            readOnly
            disabled
            style={{ backgroundColor: '#f0f0f0' }}
          />
        </div>
        <div className="form-group">
          <label>标题</label>
          <input
            value={formData.title}
            onChange={(e) => setFormData({...formData, type: e.target.value})}
          />
        </div>
        

          <div className="form-group">
              <label htmlFor="">内容:</label>
              <textarea
                style={{ height: '100px', width: '100%' }}
                placeholder='内容'
                value={formData.content}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              ></textarea>
          </div>
          
          <div className="form-group">
              <label htmlFor="">翻译:</label>
              <textarea
                style={{ height: '100px', width: '100%' }}
                placeholder='翻译'
                value={formData.chinese}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
              ></textarea>
          </div>
        <div className="button-group">
          <button 
            type="button"
            onClick={() => handleSubmit('check')}
            disabled={isLoading}
          >
            {isLoading ? '检查中...' : '检查类似词条'}
          </button>
          
          <button 
            type="button" 
            onClick={() => handleSubmit('submit')}
            disabled={isLoading}
          >
            {isLoading ? '提交中...' : '提交词条'}
          </button>
        </div>

        <div className="button-group">
          <button 
            type="button" 
            onClick={resetFormUser}
            disabled={isLoading}
            style={{ backgroundColor: '#6c757d', color: 'white' }}
          >
            重置表单
          </button>
          
          <button 
            type="button" 
            onClick={() => handleSubmit('delete')}
            disabled={isLoading || formData.id === 0}
            style={{ backgroundColor: '#dc3545', color: 'white' }}
          >
            删除词条
          </button>
        </div>
          {apiMessage && <div className="api-message">{apiMessage}</div>}
        </form>
        {notification.show && (
          <Notification 
            message={notification.message} 
            type={notification.type} 
          />
        )}
      </div>
    )
}