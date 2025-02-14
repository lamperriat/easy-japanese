import { useState } from 'react';
import { API_BASE_URL } from '../../services/api';
export default function WordForm() {
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
  const [selectedBook, setSelectedBook] = useState('1');
  const [apiMessage, setApiMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const bookOptions = [
    { id: '1', name: '新标日初级上' },
    { id: '2', name: '新标日初级下' },
    { id: '3', name: '新标日中级上' },
    { id: '4', name: '新标日中级下' },
    { id: '5', name: '新标日高级上' },
    { id: '6', name: '新标日高级下' },
  ];
  const handleSubmit = async (actionType) => {
    setIsLoading(true);
    try {
      const endpoint = actionType === 'check' 
        ? `${API_BASE_URL}/api/words/book_${selectedBook}/check`
        : `${API_BASE_URL}/api/words/book_${selectedBook}/submit`;

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': process.env.REACT_APP_API_KEY
        },
        body: JSON.stringify(formData)
      });

      const result = await response.json();
      if (response.ok && actionType == "submit") {
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
        <div className='form-row'>
        <div className="form-group">
          <label>ID(留空会填补默认值)</label>
          <input
            type="text"
            value={formData.id}
            onChange={(e) => setFormData({...formData, id: e.target.value})}
          />
        </div>
        <div className="form-group">
          <label>类型</label>
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
            <div key={index} className="example-pair">
              <textarea
          placeholder="例句"
          value={ex.example}
          onChange={(e) => {
            const newExamples = [...formData.example];
            newExamples[index].example = e.target.value;
            setFormData({ ...formData, example: newExamples });
          }}
              />
              <textarea
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
  );
}