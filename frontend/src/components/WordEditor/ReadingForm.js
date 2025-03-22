import { useState } from 'react';
import { API_BASE_URL } from '../../services/api';


export default function ReadingForm() {
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
  const [selectedBook, setSelectedBook] = useState('1');
  const [apiMessage, setApiMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const bookOptions = [
    { id: '0', name: 'global' }, 
    { id: '-1', name: 'user' },
  ]

  const handleSubmit = async (actionType) => {
    setIsLoading(true);
    try {
      var endpoint = '';
      var method  = '';
      if (actionType === 'check') {
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/reading-material/search`;
          method = 'GET';
        } else {
          endpoint = `${API_BASE_URL}/api/reading-material/search`;
          method = 'GET';
        }
      } else {
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/reading-material/add`;
          method = 'POST';
        } else {
          endpoint = `${API_BASE_URL}/api/reading-material/add`;
          method = 'POST';
        }
      }
      const response = await fetch(endpoint, {
        method: method, 
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': process.env.REACT_APP_API_KEY
        },
        body: JSON.stringify(formData)
      });

      const result = await response.json();
      if (response.ok && actionType === "submit") {
        resetForm();
      }
      setApiMessage(result.error || result.message || "操作成功");
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

        {apiMessage && <div className="api-message">{apiMessage}</div>}
        </form>
      </div>
    )
}