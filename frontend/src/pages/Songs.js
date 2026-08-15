import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import api, { addFavorite } from '../services/api';

const Songs = () => {
  const [songs, setSongs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [search, setSearch] = useState('');

  useEffect(() => {
    fetchSongs();
  }, []);

  const fetchSongs = async () => {
    try {
      const response = await api.get('/songs');
      setSongs(response.data.songs);
    } catch (err) {
      console.error('Ошибка загрузки песен:', err);
      setError(`Ошибка загрузки песен: ${err.response?.data?.error || err.message || 'Неизвестная ошибка'}`);
    } finally {
      setLoading(false);
    }
  };

  const handleAddFavorite = async (songId) => {
    try {
      await addFavorite(songId);
      // Show success message or update UI as needed
      alert('Песня добавлена в избранное!');
    } catch (err) {
      console.error('Ошибка добавления в избранное:', err);
      alert(`Ошибка добавления в избранное: ${err.response?.data?.error || err.message || 'Неизвестная ошибка'}`);
    }
  };

  const filteredSongs = songs.filter(song =>
    song.title.toLowerCase().includes(search.toLowerCase()) ||
    song.artist.toLowerCase().includes(search.toLowerCase())
  );

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

  return (
    <div className="space-y-6">
      <div className="bg-white rounded-xl shadow-lg p-6">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between mb-6">
          <h2 className="text-3xl font-bold text-gray-800 mb-4 md:mb-0">Песни</h2>
          <div className="relative w-full md:w-96">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg className="w-5 h-5 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd" />
              </svg>
            </div>
            <input
              type="text"
              className="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              placeholder="Поиск песен..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>
        </div>
        
        {error && (
          <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-6 rounded">
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
        )}
        
        <div className="song-list">
          {filteredSongs.length > 0 ? (
            filteredSongs.map((song) => (
              <div key={song.id} className="song-list-card">
                <div className="song-card">
                  <h3 className="song-title">{song.title}</h3>
                  <p className="song-author">Автор: {song.artist}</p>
                  <p className="text-gray-600 mb-4 line-clamp-2">{song.text.substring(0, 100)}...</p>
                  <div className="flex flex-wrap gap-2">
                    <Link 
                      to={`/songs/${song.id}`} 
                      className="btn btn-primary px-4 py-2 text-sm"
                    >
                      Подробнее
                    </Link>
                    <button 
                      onClick={() => handleAddFavorite(song.id)} 
                      className="btn btn-secondary px-4 py-2 text-sm"
                    >
                      В избранное
                    </button>
                  </div>
                </div>
              </div>
            ))
          ) : (
            <div className="text-center py-12">
              <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <h3 className="mt-2 text-lg font-medium text-gray-900">Песни не найдены</h3>
              <p className="mt-1 text-gray-500">Попробуйте изменить параметры поиска</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default Songs;
