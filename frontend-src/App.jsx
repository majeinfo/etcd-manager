import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [endpoints, setEndpoints] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const fetchStatus = async () => {
    try {
      setLoading(true);
      const response = await fetch('http://localhost:8080/api/status');
      const data = await response.json();
      setEndpoints(data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch endpoints status');
    } finally {
      setLoading(false);
    }
  };

  const compactDatabase = async () => {
    try {
      setLoading(true);
      const response = await fetch('http://localhost:8080/api/compact', {
        method: 'POST',
      });
      const data = await response.json();
      if (response.ok) {
        alert('Database compacted successfully');
        fetchStatus(); // Refresh status
      } else {
        throw new Error(data.error);
      }
    } catch (err) {
      setError('Failed to compact database');
    } finally {
      setLoading(false);
    }
  };

  const defragDatabase = async () => {
    try {
      setLoading(true);
      const response = await fetch('http://localhost:8080/api/defrag', {
        method: 'POST',
      });
      const data = await response.json();
      if (response.ok) {
        alert('Database defragmented successfully');
        fetchStatus(); // Refresh status
      } else {
        throw new Error(data.error);
      }
    } catch (err) {
      setError('Failed to defragment database');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStatus();
    const interval = setInterval(fetchStatus, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <h1>etcd Cluster Manager</h1>
      </header>
      
      <main className="App-main">
        {error && (
          <div className="error-message">
            {error}
          </div>
        )}
        
        <div className="actions">
          <button 
            onClick={compactDatabase} 
            disabled={loading}
            className="action-button"
          >
            Compact Database
          </button>
          <button 
            onClick={defragDatabase} 
            disabled={loading}
            className="action-button"
          >
            Defragment Database
          </button>
          <button 
            onClick={fetchStatus} 
            disabled={loading}
            className="action-button refresh"
          >
            Refresh Status
          </button>
        </div>

        <div className="endpoints-grid">
          {endpoints.map((endpoint, index) => (
            <div key={index} className="endpoint-card">
              <h3>{endpoint.endpoint}</h3>
              <div className="endpoint-details">
                <p>
                  <strong>Version:</strong> {endpoint.version}
                </p>
                <p>
                  <strong>DB Size:</strong> {(endpoint.dbSize / (1024 * 1024)).toFixed(2)} MB
                </p>
                <p>
                  <strong>DB Size In Use:</strong> {(endpoint.dbSizeInUse / (1024 * 1024)).toFixed(2)} MB
                </p>
                <p>
                  <strong>Leader:</strong> {endpoint.leader ? 'Yes' : 'No'}
                </p>
              </div>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}

export default App

