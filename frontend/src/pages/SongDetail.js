import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api, { deleteSong } from '../services/api';

const SongDetail = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [song, setSong] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [userId, setUserId] = useState(null);

  useEffect(() => {
    // Получаем ID пользователя из токена или localStorage
    const token = localStorage.getItem('token');
    if (token) {
      try {
        // Декодируем JWT токен для получения ID пользователя
        const payload = JSON.parse(atob(token.split('.')[1]));
        setUserId(payload.userID || payload.userId || payload.id);
      } catch (e) {
        console.error('Ошибка декодирования токена:', e);
      }
    }
    
    fetchSong();
  }, [id]);

  const fetchSong = async () => {
    try {
      const response = await api.get(`/songs/${id}`);
      setSong(response.data.song);
    } catch (err) {
      setError('Ошибка загрузки песни');
    } finally {
      setLoading(false);
    }
  };

  const handleAddToFavorites = async () => {
    try {
      await api.post(`/users/me/favorites/${id}`);
      // Обновляем состояние песни, чтобы показать, что она в избранном
      setSong({ ...song, is_favorite: true });
    } catch (err) {
      console.error('Ошибка добавления в избранное:', err);
    }
  };

  const handleRemoveFromFavorites = async () => {
    try {
      await api.delete(`/users/me/favorites/${id}`);
      // Обновляем состояние песни, чтобы показать, что она удалена из избранного
      setSong({ ...song, is_favorite: false });
    } catch (err) {
      console.error('Ошибка удаления из избранного:', err);
    }
  };

  const handleDeleteSong = async () => {
    if (!window.confirm('Вы уверены, что хотите удалить эту песню?')) {
      return;
    }

    try {
      await deleteSong(id);
      // Перенаправляем пользователя на страницу со списком песен
      navigate('/songs');
    } catch (err) {
      console.error('Ошибка удаления песни:', err);
      alert(`Ошибка удаления песни: ${err.response?.data?.error || err.message || 'Неизвестная ошибка'}`);
    }
  };

  if (loading) {
    return (
      <div className="container">
        <div className="flex justify-center items-center h-64">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
          <span className="ml-3 text-xl">Загрузка...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container">
        <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-sm text-red-700">{error}</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (!song) {
    return (
      <div className="container">
        <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 rounded">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-sm text-yellow-700">Песня не найдена</p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Проверяем, является ли текущий пользователь владельцем песни
  // Также проверяем, есть ли у песни createdById (в некоторых случаях может отсутствовать)
  const isOwner = userId && (song.createdById === userId || song.user_id === userId);

  return (
    <div className="bg-white rounded-xl shadow-lg p-6">
      <div className="flex flex-col md:flex-row md:items-start md:justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-800 mb-2">{song.title}</h1>
          <p className="text-xl text-gray-600">Автор: {song.artist}</p>
        </div>
        
        <div className="flex flex-wrap gap-2 mt-4 md:mt-0">
          {song.is_favorite ? (
            <button 
              onClick={handleRemoveFromFavorites} 
              className="btn favorite-btn px-4 py-2 flex items-center"
            >
              <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z" clipRule="evenodd" />
              </svg>
              В избранном
            </button>
          ) : (
            <button 
              onClick={handleAddToFavorites} 
              className="btn btn-secondary px-4 py-2 flex items-center"
            >
              <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path>
              </svg>
              В избранное
            </button>
          )}
          
          {isOwner && (
            <button 
              onClick={handleDeleteSong} 
              className="btn btn-danger px-4 py-2 flex items-center"
            >
              <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
              </svg>
              Удалить
            </button>
          )}
        </div>
      </div>
      
          {song.mp3_url && (
            <div className="bg-gray-50 rounded-lg p-6 mb-6">
              <h3 className="text-xl font-semibold text-gray-800 mb-4 flex items-center">
                <svg className="w-6 h-6 mr-2 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3"></path>
                </svg>
                Аудиозапись
              </h3>
              <div className="audio-player">
                <audio controls>
                  <source src={song.mp3_url} type="audio/mpeg" />
                  Ваш браузер не поддерживает аудио элемент.
                </audio>
              </div>
            </div>
          )}
      
      {song.rhythm && (
        <div className="bg-gray-50 rounded-lg p-6 mb-6">
          <h3 className="text-xl font-semibold text-gray-800 mb-4 flex items-center">
            <svg className="w-6 h-6 mr-2 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z"></path>
            </svg>
            Ритмическая схема
          </h3>
          <div className="rhythm-pattern font-mono bg-gray-100 p-4 rounded-md overflow-x-auto">
            {song.rhythm}
          </div>
        </div>
      )}
      
      <div className="bg-gray-50 rounded-lg p-6">
        <h3 className="text-xl font-semibold text-gray-800 mb-4 flex items-center">
          <svg className="w-6 h-6 mr-2 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
          </svg>
          Текст песни
        </h3>
        <div className="song-text">{song.text}</div>
      </div>
    </div>
  );
};

export default SongDetail;