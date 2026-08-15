import React from 'react';
import { Link } from 'react-router-dom';

const Home = () => {
  return (
    <div className="space-y-8">
      <div className="bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl shadow-lg p-8 text-white">
        <h1 className="text-4xl font-bold mb-4">Добро пожаловать в HolyFlow</h1>
        <p className="text-xl mb-6">Сервис для хранения текстов христианских песен, их ритмичности и аудиозаписей.</p>
        <div className="bg-white/20 backdrop-blur-sm rounded-lg p-6 mb-6">
          <h2 className="text-2xl font-semibold mb-4">Здесь вы можете:</h2>
          <ul className="grid grid-cols-1 md:grid-cols-2 gap-3 mb-6">
            <li className="flex items-center">
              <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              Искать и просматривать тексты песен
            </li>
            <li className="flex items-center">
              <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              Слушать аудиозаписи песен
            </li>
            <li className="flex items-center">
              <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              Изучать ритмичность песен
            </li>
            <li className="flex items-center">
              <svg className="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              Создавать свой список любимых песен
            </li>
          </ul>
          <div className="mt-4">
            <Link to="/songs" className="btn btn-primary inline-flex items-center px-6 py-3 text-lg">
              Просмотреть песни
              <svg className="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path>
              </svg>
            </Link>
          </div>
        </div>
      </div>
      
      <div className="bg-white rounded-xl shadow-lg p-8">
        <h2 className="text-3xl font-bold text-gray-800 mb-4">О проекте</h2>
        <p className="text-gray-600 text-lg leading-relaxed">
          HolyFlow создан для того, чтобы помочь христианам находить и изучать песни для поклонения. 
          Наш сервис предоставляет удобный способ хранить и организовывать тексты песен, 
          а также доступ к аудиозаписям и ритмическим схемам.
        </p>
      </div>
    </div>
  );
};

export default Home;
