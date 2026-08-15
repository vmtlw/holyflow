import React from 'react';
import { Link, useNavigate } from 'react-router-dom';

const Navbar = ({ user, onLogout }) => {
  const navigate = useNavigate();

  const handleLogout = () => {
    onLogout();
    navigate('/');
  };

  return (
    <nav className="bg-gray-800 py-4 mb-6">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center">
          <Link to="/" className="text-white text-2xl font-bold no-underline">HolyFlow</Link>
          <ul className="flex space-x-6">
            <li>
              <Link to="/" className="text-white no-underline px-3 py-2 rounded hover:bg-gray-700 transition duration-300">Главная</Link>
            </li>
            <li>
              <Link to="/songs" className="text-white no-underline px-3 py-2 rounded hover:bg-gray-700 transition duration-300">Песни</Link>
            </li>
            {user ? (
              <>
                <li>
                  <Link to="/add-song" className="text-white no-underline px-3 py-2 rounded hover:bg-gray-700 transition duration-300">Добавить песню</Link>
                </li>
                <li>
                  <Link to="/my-songs" className="text-white no-underline px-3 py-2 rounded hover:bg-gray-700 transition duration-300">Мои песни</Link>
                </li>
                <li>
                  <Link to="/profile" className="text-white no-underline px-3 py-2 rounded hover:bg-gray-700 transition duration-300">Профиль</Link>
                </li>
                <li>
                  <button 
                    onClick={handleLogout} 
                    className="text-white no-underline px-3 py-2 rounded hover:bg-gray-700 transition duration-300 cursor-pointer bg-transparent border-none"
                  >
                    Выйти
                  </button>
                </li>
              </>
            ) : (
              <>
                <li>
                  <Link to="/login" className="text-white no-underline px-3 py-2 rounded hover:bg-gray-700 transition duration-300">Вход</Link>
                </li>
                <li>
                  <Link to="/register" className="btn btn-primary px-3 py-2">Регистрация</Link>
                </li>
              </>
            )}
          </ul>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
