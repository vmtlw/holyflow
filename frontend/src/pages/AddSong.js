import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

const AddSong = () => {
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    title: '',
    artist: '',
    text: '',
    rhythm: '',
    mp3: null,
    cover: null
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleFileChange = (e) => {
    const { name, files } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: files[0]
    }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setSuccess(false);

    try {
      // Сначала создаем песню без файлов
      const songData = {
        title: formData.title,
        artist: formData.artist,
        text: formData.text,
        rhythm: formData.rhythm
      };

      const response = await api.post('/songs', songData);
      const songId = response.data.song.id;

      // Если есть MP3 файл, загружаем его
      if (formData.mp3) {
        const mp3Data = new FormData();
        mp3Data.append('file', formData.mp3);
        
        await api.post(`/songs/${songId}/upload-mp3`, mp3Data, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
      }

      // Если есть обложка, загружаем ее
      if (formData.cover) {
        const coverData = new FormData();
        coverData.append('file', formData.cover);
        
        await api.post(`/songs/${songId}/upload-cover`, coverData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
      }

      setSuccess(true);
      // Переходим к списку песен через 2 секунды
      setTimeout(() => {
        navigate('/songs');
      }, 2000);
    } catch (err) {
      setError(err.response?.data?.error || 'Ошибка добавления песни');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-xl shadow-lg p-6">
      <h2 className="text-3xl font-bold text-gray-800 mb-6">Добавить новую песню</h2>
      
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
      
      {success && (
        <div className="bg-green-50 border-l-4 border-green-500 p-4 mb-6 rounded">
          <div className="flex">
            <div className="flex-shrink-0">
              <svg className="h-5 w-5 text-green-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            </div>
            <div className="ml-3">
              <p className="text-sm text-green-700">Песня успешно добавлена!</p>
            </div>
          </div>
        </div>
      )}
      
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="form-group">
          <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">Название песни:</label>
          <input
            type="text"
            id="title"
            name="title"
            className="form-control"
            value={formData.title}
            onChange={handleChange}
            required
            placeholder="Введите название песни"
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="artist" className="block text-sm font-medium text-gray-700 mb-1">Автор:</label>
          <input
            type="text"
            id="artist"
            name="artist"
            className="form-control"
            value={formData.artist}
            onChange={handleChange}
            required
            placeholder="Введите имя автора"
          />
        </div>
        
        <div className="form-group">
          <label htmlFor="text" className="block text-sm font-medium text-gray-700 mb-1">Текст песни:</label>
          <textarea
            id="text"
            name="text"
            className="form-control"
            rows="10"
            value={formData.text}
            onChange={handleChange}
            required
            placeholder="Введите текст песни"
          ></textarea>
        </div>
        
        <div className="form-group">
          <label htmlFor="rhythm" className="block text-sm font-medium text-gray-700 mb-1">Ритмическая схема:</label>
          <textarea
            id="rhythm"
            name="rhythm"
            className="form-control"
            rows="3"
            value={formData.rhythm}
            onChange={handleChange}
            placeholder="Введите ритмическую схему (необязательно)"
          ></textarea>
        </div>
        
        <div className="form-group">
          <label htmlFor="mp3" className="block text-sm font-medium text-gray-700 mb-1">MP3 файл:</label>
          <input
            type="file"
            id="mp3"
            name="mp3"
            className="form-control block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
            accept="audio/mp3,audio/mpeg"
            onChange={handleFileChange}
          />
          {formData.mp3 && <div className="file-info mt-1 text-sm text-gray-500">Выбран файл: {formData.mp3.name}</div>}
        </div>
        
        <div className="form-group">
          <label htmlFor="cover" className="block text-sm font-medium text-gray-700 mb-1">Обложка:</label>
          <input
            type="file"
            id="cover"
            name="cover"
            className="form-control block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
            accept="image/*"
            onChange={handleFileChange}
          />
          {formData.cover && <div className="file-info mt-1 text-sm text-gray-500">Выбран файл: {formData.cover.name}</div>}
        </div>
        
        <button 
          type="submit" 
          className="btn btn-primary w-full py-3 px-4 rounded-md shadow-sm text-sm font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:w-auto sm:text-sm"
          disabled={loading}
        >
          {loading ? (
            <div className="flex items-center justify-center">
              <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Добавление...
            </div>
          ) : 'Добавить песню'}
        </button>
      </form>
    </div>
  );
};

export default AddSong;
