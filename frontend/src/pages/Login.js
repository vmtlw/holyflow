import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import api from '../services/api';

const Login = ({ onLogin }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');

    try {
      const response = await api.post('/auth/login', {
        username,
        password
      });

      const { user, token } = response.data;
      onLogin(user, token);
      navigate('/');
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка входа');
    } finally {
      setLoading(false);
    }
  };

  const handleGitLabLogin = () => {
    // Redirect to GitLab OAuth
    window.location.href = `${api.defaults.baseURL}/auth/gitlab/login`;
  };

  return (
    <div className="card">
      <h2>Вход</h2>
      {error && <div className="alert alert-danger">{error}</div>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="username">Логин или Email:</label>
          <input
            type="text"
            id="username"
            className="form-control"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        </div>
        <div className="form-group">
          <label htmlFor="password">Пароль:</label>
          <input
            type="password"
            id="password"
            className="form-control"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        <button type="submit" className="btn" disabled={loading}>
          {loading ? 'Вход...' : 'Войти'}
        </button>
      </form>
      
      <div style={{ marginTop: '20px' }}>
        <button onClick={handleGitLabLogin} className="btn btn-secondary" style={{ marginRight: '10px' }}>
          Войти через GitLab
        </button>
      </div>
      
      <div style={{ marginTop: '20px' }}>
        <Link to="/register">Нет аккаунта? Зарегистрируйтесь</Link>
      </div>
      <div style={{ marginTop: '10px' }}>
        <Link to="/forgot-password">Забыли пароль?</Link>
      </div>
    </div>
  );
};

export default Login;