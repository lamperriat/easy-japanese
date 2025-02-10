import { useState } from 'react';
export default function WordForm() {
  const [formData, setFormData] = useState({
    id: '',
    kanji: '',
    hiragana: '',
    english: '',
    example: ''
  });
  
  const [apiMessage, setApiMessage] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (actionType) => {
    setIsLoading(true);
    try {
      const endpoint = actionType === 'check' 
        ? '/api/words/check' 
        : '/api/words/submit';

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': process.env.REACT_APP_API_KEY
        },
        body: JSON.stringify(formData)
      });

      const result = await response.json();
      setApiMessage(result.message || "操作成功");
    } catch (error) {
      setApiMessage("网络请求失败");
    }
    setIsLoading(false);
  };

  return (
    <div className="word-editor">
      <form>
        <div className="form-group">
          <label>唯一标识符</label>
          <input
            type="text"
            value={formData.id}
            onChange={(e) => setFormData({...formData, id: e.target.value})}
          />
        </div>

        <div className="form-row">
          <div className="form-group">
            <label>汉字表记</label>
            <input
              value={formData.kanji}
              onChange={(e) => setFormData({...formData, kanji: e.target.value})}
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

        <div className="form-group">
          <label>英文翻译</label>
          <input
            value={formData.english}
            onChange={(e) => setFormData({...formData, english: e.target.value})}
          />
        </div>

        <div className="form-group">
          <label>例句</label>
          <textarea
            value={formData.example}
            onChange={(e) => setFormData({...formData, example: e.target.value})}
          />
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