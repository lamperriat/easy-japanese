import { API_BASE_URL } from '../services/api'; 

import "./AdminPage.css";
import React, { useState } from 'react';
import axios from 'axios';

const QueryTypes = {
  APIKEY_ADD: 'apikeyAdd',
  APIKEY_DEL: 'apikeyDel',
  ACCOUNT_ADD: 'accountAdd',
  INVALID: 'invalid',
  [Symbol.for('isEnum')]: true
}

Object.freeze(QueryTypes);

const AdminPage = () => {
  const [credentials, setCredentials] = useState({
    username: '',
    password: ''
  });
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [apiResponse, setApiResponse] = useState(null);
  const [loading, setLoading] = useState(false);
  const [apiKeyInput, setApiKeyInput] = useState('');
  const [queryType, setQueryType] = useState(QueryTypes.INVALID); // apikey or account
  const [newAccount, setNewAccount] = useState({
    username: '',
    password: ''
  });

  const handleLogin = (e) => {
    e.preventDefault();
    setIsAuthenticated(true);
  };

  const handleLogout = () => {
    setIsAuthenticated(false);
    setCredentials({ username: '', password: '' });
    setApiResponse(null);
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setCredentials(prev => ({ ...prev, [name]: value }));
  };

  const handleAccountInputChange = (e) => {
    const { name, value } = e.target;
    setNewAccount(prev => ({ ...prev, [name]: value }));
  };

  const createApiKey = async () => {
    setLoading(true);
    try {
      const response = await axios.get(`${API_BASE_URL}/api/admin/apikey/create`, {
        headers: {
          'Admin-Name': credentials.username,
          'Admin-Password': credentials.password
        }
      });
      setQueryType(QueryTypes.APIKEY_ADD);
      setApiResponse(response.data);
    } catch (error) {
      setApiResponse(error.response?.data || { error: 'An error occurred' });
    } finally {
      setLoading(false);
    }
  };

  const deleteApiKey = async () => {
    if (!apiKeyInput) {
      setApiResponse({ error: 'Please enter an API key to delete' });
      return;
    }

    setLoading(true);
    try {
      const response = await axios.post(
        `${API_BASE_URL}/api/admin/apikey/delete`,
        { key: apiKeyInput },
        {
          headers: {
            'Admin-Name': credentials.username,
            'Admin-Password': credentials.password
          }
        }
      );
      setQueryType(QueryTypes.APIKEY_DEL);
      setApiResponse(response.data);
      setApiKeyInput('');
    } catch (error) {
      setApiResponse(error.response?.data || { error: 'An error occurred' });
    } finally {
      setLoading(false);
    }
  };

  const createAccount = async () => {
    if (!newAccount.username || !newAccount.password) {
      setApiResponse({ error: 'Username and password are required' });
      return;
    }

    setLoading(true);
    try {
      const response = await axios.post(
        `${API_BASE_URL}/api/admin/account/create`,
        newAccount,
        {
          headers: {
            'Admin-Name': credentials.username,
            'Admin-Password': credentials.password
          }
        }
      );
      setQueryType(QueryTypes.ACCOUNT_ADD);
      setApiResponse(response.data);
      setNewAccount({ username: '', password: '' });
    } catch (error) {
      setApiResponse(error.response?.data || { error: 'An error occurred' });
    } finally {
      setLoading(false);
    }
  };

  if (!isAuthenticated) {
    return (
      <div className="admin-container">
        <div className="login-card">
          <h1 className="login-title">Admin Login</h1>
          <form onSubmit={handleLogin} className="login-form">
            <div className="form-group">
              <label htmlFor="username" className="form-label">Username:</label>
              <input
                type="text"
                id="username"
                name="username"
                value={credentials.username}
                onChange={handleInputChange}
                className="form-input"
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="password" className="form-label">Password:</label>
              <input
                type="password"
                id="password"
                name="password"
                value={credentials.password}
                onChange={handleInputChange}
                className="form-input"
                required
              />
            </div>
            <button type="submit" className="login-button">Login</button>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div className="admin-container">
      <div className="dashboard-card">
        <div className="admin-header">
          <div>
            <h1 className="dashboard-title">Admin Dashboard</h1>
            <p className="user-info">Logged in as: <span className="username">{credentials.username}</span></p>
          </div>
          <button onClick={handleLogout} className="logout-button">Logout</button>
        </div>

        <div className="section-divider"></div>

        <div className="admin-section">
          <h2 className="section-title">API Key Management</h2>
          
          <div className="action-card">
            <h3 className="action-title">Create API Key</h3>
            <p className="action-description">Generate a new API key for authentication</p>
            <button 
              onClick={createApiKey} 
              disabled={loading}
              className={`action-button ${loading ? 'loading' : ''}`}
            >
              {loading ? 'Generating...' : 'Create New API Key'}
            </button>
          </div>
          {apiResponse && queryType === QueryTypes.APIKEY_ADD &&  (
            <div className={`api-response ${apiResponse.error ? 'error' : 'success'}`}>
              <h3 className="response-title">New API KEY:</h3>
              <pre className="response-content">{apiResponse.key + "\nPlease copy and save it properly. You will not have access to it after refreshing the page." || apiResponse.error}</pre>
            </div>
          )}

          <div className="action-card">
            <h3 className="action-title">Delete API Key</h3>
            <p className="action-description">Remove an existing API key</p>
            <input
              type="text"
              placeholder="Enter API key to delete"
              value={apiKeyInput}
              onChange={(e) => setApiKeyInput(e.target.value)}
              className="form-input"
            />
            <button 
              onClick={deleteApiKey} 
              disabled={loading || !apiKeyInput}
              className={`action-button ${loading ? 'loading' : ''}`}
            >
              {loading ? 'Deleting...' : 'Delete API Key'}
            </button>
          </div>
        </div>
        {apiResponse && queryType === QueryTypes.APIKEY_DEL &&  (
            <div className={`api-response ${apiResponse.error ? 'error' : 'success'}`}>
              <h3 className="response-title">Result:</h3>
              <pre className="response-content">{apiResponse.message || apiResponse.error}</pre>
            </div>
          )}

        <div className="section-divider"></div>

        <div className="admin-section">
          <h2 className="section-title">Account Management</h2>
          
          <div className="action-card">
            <h3 className="action-title">Create New Account</h3>
            <p className="action-description">Create a new user account</p>
            <div className="form-group">
              <label className="form-label">Username:</label>
              <input
                type="text"
                name="username"
                value={newAccount.username}
                onChange={handleAccountInputChange}
                className="form-input"
              />
            </div>
            <div className="form-group">
              <label className="form-label">Password:</label>
              <input
                type="password"
                name="password"
                value={newAccount.password}
                onChange={handleAccountInputChange}
                className="form-input"
              />
            </div>
            <button 
              onClick={createAccount} 
              disabled={loading || !newAccount.username || !newAccount.password}
              className={`action-button ${loading ? 'loading' : ''}`}
            >
              {loading ? 'Creating...' : 'Create Account'}
            </button>
          </div>
          {apiResponse && queryType === QueryTypes.ACCOUNT_ADD &&  (
            <div className={`api-response ${apiResponse.error ? 'error' : 'success'}`}>
              <h3 className="response-title">Result:</h3>
              <pre className="response-content">{apiResponse.message || apiResponse.error}</pre>
            </div>
          )}
        </div>


      </div>
    </div>
  );
};

export default AdminPage;