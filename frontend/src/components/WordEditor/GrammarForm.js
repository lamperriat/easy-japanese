import { useEffect, useState } from 'react';
import { API_BASE_URL } from '../../services/api';
import Notification from '../Auth/Notification';

export default function GrammarForm({ initGrammarData, initBookId }) {
  const [formData, setFormData] = useState({
    id: 0,
    description: '', 
    example: []
  });
  const resetForm = () => {
    setFormData({
      id: 0,
      description: '',
      example: []
    });
  }
  const resetFormUser = () => {
    resetForm();
    setApiMessage('表单清空完成');
  }
  const [selectedBook, setSelectedBook] = useState('1');
  const [apiMessage, setApiMessage] = useState('');
  const [enableMsgButton, setEnableMsgButton] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [notification, setNotification] = useState({ show: false, message: '', type: '' });
  const [searchResult, setSearchResult] = useState({
    id: 0,
    description: '',
    example: []
  });
  
  const bookOptions = [
    { id: '1', name: '全局' }, 
    { id: '-1', name: '个人' },
  ]
  useEffect(() => {
    if (initBookId) {
      setSelectedBook(initBookId);
    }
    if (initGrammarData) {
      setFormData({
        id: initGrammarData.id,
        description: initGrammarData.description,
        example: initGrammarData.example || []
      });
      setApiMessage('加载完成。清注意：此时您的提交会修改该词条，而非新增。');
    }
  }, [initGrammarData, initBookId]);

  const handleSubmit = async (actionType) => {
    setIsLoading(true);
    try {
      var endpoint = '';
      var method  = '';
      if (actionType === 'check') {
        // TODO: use search-title api to improve performance
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/grammar/search?query=${formData.description}`;
          method = 'GET';
        } else {
          endpoint = `${API_BASE_URL}/api/grammar/search?query=${formData.description}`;
          method = 'GET';
        }
      } else if (actionType === 'submit') {
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/grammar/add`;
          method = 'POST';
        } else {
          endpoint = `${API_BASE_URL}/api/grammar/add`;
          method = 'POST';
        }
      } else if (actionType === 'delete') {
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/grammar/delete`;
          method = 'POST';
        } else {
          endpoint = `${API_BASE_URL}/api/grammar/delete`;
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
      var response;
      if (method === 'GET') {
        response = await fetch(endpoint, {
          method: method,
          headers: {
            'Content-Type': 'application/json',
            'Authorization': token
          }
        });
      } else {
        response = await fetch(endpoint, {
          method: method,
          headers: {
            'Content-Type': 'application/json',
            'Authorization': token
          },
          body: JSON.stringify(formData)
        });
      }

      const result = await response.json();
      if (response.ok && actionType !== "check") {
        resetForm();
      }
      var searchMsg = "";
      console.log(result);
      if (actionType === 'check') {
        const searchResults = result.results;
        if (searchResults.length > 0) {
          searchMsg = `找到 ${result.length} 个相似的词条：`;
          if (searchResults.length > 3) {
            searchMsg += `\n- ${searchResults[0].description}`;
            searchMsg += `\n- ${searchResults[1].description}`;
            searchMsg += `\n- ${searchResults[2].description}`;
            searchMsg += `\n...`;
          } else {
            searchResults.forEach(item => {
              searchMsg += `\n- ${item.description}`;
            });
          }
          searchMsg += "\n加载第一条到表单吗？";
          setEnableMsgButton(true);
          if (searchResults[0].example === null) {
            searchResults[0].example = [];
          }
          setSearchResult(searchResults[0]);
        } else {
          searchMsg = "没有找到相似的词条";
        }
        setApiMessage(searchMsg);
      } else {
        setApiMessage("操作成功");
      }
    } catch (error) {
      console.error("Error:", error);
      setApiMessage("网络请求失败");
    }
    setIsLoading(false);
  };
  const handleMsgButtonClick = () => {
    setFormData({
      ...formData,
      id: searchResult.id,
      description: searchResult.description,
      example: searchResult.example
    });
    setApiMessage('加载完成。清注意：此时您的提交会修改该词条，而非新增。');
    setEnableMsgButton(false);
  }

    return (
      <div className="word-editor">
      <form>
        <div className="form-row">
        
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
        <label>选择数据库：</label>
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
        </div>
        <div className="form-group">
          <label htmlFor="description">语法:</label>
          <textarea
          style={{ height: '100px', width: '100%' }}
          placeholder='语法描述'
          value={formData.description}
          onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          ></textarea>
        </div>
        <div className="form-group">
        <label>例句</label>
        {formData.example.map((ex, index) => (
        <div key={index}>
          <textarea
        style={{height: '40px'}}
        placeholder="例句"
        value={ex.example}
        onChange={(e) => {
        const newExamples = [...formData.example];
        newExamples[index].example = e.target.value;
        setFormData({ ...formData, example: newExamples });
        }}
          />
          <textarea
        style={{height: '40px'}}
        placeholder="中文翻译"
        value={ex.chinese}
        onChange={(e) => {
        const newExamples = [...formData.example];
        newExamples[index].chinese = e.target.value;
        setFormData({ ...formData, example: newExamples });
        }}
          />
        </div>
        ))}
        <div style={{ marginTop: '10px' }}>
        <button
          type="button"
          onClick={() => {
          setFormData({
            ...formData,
            example: [...formData.example, { example: '', chinese: '' }]
          });
          }}
        >
          +
        </button>
        <button
          type="button"
          onClick={() => {
          if (formData.example.length > 0) {
            const newExamples = [...formData.example];
            newExamples.pop();
            setFormData({
            ...formData,
            example: newExamples
            });
          }
          }}
          disabled={formData.example.length === 0}
          style={{ marginLeft: '2px' }}
        >
          -
        </button>
        </div>
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
      {apiMessage && (
        <div className="api-message">
        {apiMessage}
        {enableMsgButton && (
          <button 
          type="button" 
          onClick={handleMsgButtonClick}
          style={{ marginLeft: '10px' }}
          >
          点击加载
          </button>
        )}
        </div>
      )}
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