import React, { useState } from 'react';
import './App.css';

function App() {
  const [randomNumber, setRandomNumber] = useState(null);
  const [loading, setLoading] = useState(false);

  const fetchRandomNumber = async () => {
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/api/random', {
        headers: {
          'X-Api-Key': process.env.REACT_APP_API_KEY 
        }
      });
      const data = await response.json();
      setRandomNumber(data.random);
    } catch (error) {
      console.error('Error fetching random number:', error);
    }
    setLoading(false);
  };

  return (
    <div className="App">
      <h1>Random Number Generator</h1>
      <button 
        onClick={fetchRandomNumber}
        disabled={loading}
      >
        {loading ? 'Loading...' : 'Get Random Number'}
      </button>
      {randomNumber !== null && (
        <div className="result">
          <h2>Your Random Number:</h2>
          <p>{randomNumber}</p>
        </div>
      )}
    </div>
  );
}

export default App;