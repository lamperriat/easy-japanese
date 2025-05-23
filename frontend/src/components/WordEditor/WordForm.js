import { useEffect, useState } from 'react';
import { API_BASE_URL } from '../../services/api';
import Notification from '../Auth/Notification';
export default function WordForm({ initWordData, initBookId }) {
  const [formData, setFormData] = useState({
    id: 0,
    kanji: '',
    chinese: '',
    katakana: '',
    hiragana: '',
    type: '',
    example: []
  });
  const resetForm = () => {
    setFormData({
      id: 0,
      kanji: '',
      chinese: '',
      katakana: '',
      hiragana: '',
      type: '',
      example: []
    });
  }
  const resetFormUser = () => {
    resetForm();
    setApiMessage('表单清空完成');
    setEnableMsgButton(false);
  }
  const [selectedBook, setSelectedBook] = useState('1');
  const [apiMessage, setApiMessage] = useState('');
  const [enableMsgButton, setEnableMsgButton] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [notification, setNotification] = useState({ show: false, message: '', type: '' });
  const [searchResult, setSearchResult] = useState({
    id: 0,
    kanji: '',
    chinese: '',
    katakana: '',
    hiragana: '',
    type: '',
    example: []
  });
  
  
  const bookOptions = [
    { id: '1', name: '新标日初级上' },
    { id: '2', name: '新标日初级下' },
    { id: '3', name: '新标日中级上' },
    { id: '4', name: '新标日中级下' },
    { id: '5', name: '新标日高级上' },
    { id: '6', name: '新标日高级下' },
    { id: '-1', name: '个人'}, 
  ];

  useEffect(() => {
    if (initBookId) {
      setSelectedBook(initBookId);
    }
    if (initWordData) {
      setFormData({
        id: initWordData.id,
        kanji: initWordData.kanji,
        chinese: initWordData.chinese,
        katakana: initWordData.katakana,
        hiragana: initWordData.hiragana,
        type: initWordData.type,
        example: initWordData.example
      });
      setApiMessage('加载完成。清注意：此时您的提交会修改该词条，而非新增。');
      setEnableMsgButton(false);
    }
  }, [initWordData, initBookId]);
  const handleSubmit = async (actionType) => {
    setIsLoading(true);
    setEnableMsgButton(false);
    try {
      var endpoint = '';
      var method  = '';
      if (actionType === 'check') {
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/words/accurate-search`;
          method = 'POST';
        } else {
          endpoint = `${API_BASE_URL}/api/words/book_${selectedBook}/accurate-search`;
          method = 'POST';
        }
      } else if (actionType === 'submit') {
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/words`;
          method = 'POST';
        } else {
          endpoint = `${API_BASE_URL}/api/words/book_${selectedBook}`;
          method = 'POST';
        }
        if (formData.id === 0) {
          // 0: add new
          endpoint += '/add';
        } else {
          // 1: update
          endpoint += '/edit';
        }
      } else if (actionType === 'delete') {
        if (selectedBook === '-1') {
          endpoint = `${API_BASE_URL}/api/user/words/delete`;
          method = 'POST';
        } else {
          endpoint = `${API_BASE_URL}/api/words/book_${selectedBook}/delete`;
          method = 'POST';
        }
      } else {
        setApiMessage('未知操作');
        setIsLoading(false);
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
      var searchMsg = '';
      console.log(result);
      if (actionType === "check") {
        if (result.length > 0) {
          searchMsg = `找到${result.length}个相似词条: `;
          result.forEach((item) => {
            if (item.kanji) {
              searchMsg += `${item.kanji}(${item.hiragana})`;
            } else if (item.katakana) {
              searchMsg += `${item.katakana}`;
            }
            searchMsg += `; `;
          });
          searchMsg += "\n加载第一条到表单吗？";
          setEnableMsgButton(true);
          setSearchResult(result[0]);
        } else {
          searchMsg = '没有找到相似词条';
        }
        setApiMessage(searchMsg);
      } else {
        console.log(result);
        setApiMessage('操作成功');
      }
      
    } catch (error) {
      console.error('Error:', error);
      setApiMessage("网络请求失败");
    }
    setIsLoading(false);
  };

  const handleMsgButtonClick = () => {
    setFormData({
      ...formData,
      id: searchResult.id,
      kanji: searchResult.kanji,
      chinese: searchResult.chinese,
      katakana: searchResult.katakana,
      hiragana: searchResult.hiragana,
      type: searchResult.type,
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
          <label htmlFor="book-select">选择教材:</label>
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
        <div className='form-row'>
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
          <label>类型(名/动1/动2/动3/形1/形2/副/连)</label>
          <input
            value={formData.type}
            onChange={(e) => setFormData({...formData, type: e.target.value})}
          />
        </div>
        
        </div>

        <div className="form-row">
          <div className="form-group">
            <label>汉字</label>
            <input
              value={formData.kanji}
              onChange={(e) => setFormData({...formData, kanji: e.target.value})}
            />
          </div>
          <div className="form-group">
            <label>中文</label>
            <input
              value={formData.chinese}
              onChange={(e) => setFormData({...formData, chinese: e.target.value})}
            />
          </div>
        </div>
        <div className="form-row">
          <div className="form-group">
            <label>片假名</label>
            <input
              value={formData.katakana}
              onChange={(e) => setFormData({...formData, katakana: e.target.value})}
            />
          </div>          
          <div className="form-group">
            <label>平假名</label>
            <input
              value={formData.hiragana}
              onChange={(e) => setFormData({...formData, hiragana: e.target.value})}
            />
          </div>
        </div>

        {/* <div className="form-group">
          <label>例句</label>
          <textarea
            value={formData.example}
            onChange={(e) => setFormData({...formData, example: e.target.value})}
          />
        </div> */}
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
  );
}
