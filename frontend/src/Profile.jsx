import './Profile.css';
import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import BottomMenu from './components/BottomMenu.jsx';
import logoIcon from './assets/logo.svg';
import notificationIcon from './assets/notifications.svg';
import settingsIcon from './assets/settings.svg';
import avatarImg from './assets/mock_ava.svg';    
import pinIcon from './assets/location.svg';
import verifiedIcon from './assets/verified_prof.svg';
import helpedIcon from './assets/helped.svg';
import requestsIcon from './assets/requests.svg';
import rightArrowIcon from './assets/right-arrow.svg';
import personalIcon from './assets/personal.svg';
import supportIcon from './assets/support.svg';
import logoutIcon from './assets/logout.svg';
import chevronIcon from './assets/chevron.svg';
import thanksIcon from './assets/thanks.svg';
import cityImg from './assets/city.svg';         
import EventCard from './EventCard.jsx';

 function Profile() {
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const API_URL = import.meta.env.VITE_API_URL;

   const useMockData = () => {
    setUser({
      name: 'Иван Петров',
      location: 'Иннополис, Россия',
      verified: true,
      helped: 24,
      requests: 16,
    });
  };

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const token = localStorage.getItem('token');
        const response = await fetch(`${API_URL}/me`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (response.ok) {
          const data = await response.json();
          setUser({
            name: data.name,
            email: data.email,
            location: 'Иннополис, Россия',
            verified: true,
            helped: 24,
            requests: 16,
          });
        } else {
          useMockData();
        }
      } catch (error) {
        useMockData();
      } finally {
        setLoading(false);
      }
    };
    fetchUser();
  }, [API_URL]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    navigate('/login');
  };

  if (loading) 
    return (
      <>
        <div className="profile-page">
         <header className="map-header">
          <div className="logo-container">
        <img src={logoIcon} alt="Logo" style={{ width: 108, height: 25 }} />          </div>
                  <div className="header-icons">
                    <img src={notificationIcon} alt="Notifications" className="header-icon" />
                      <img src={settingsIcon} alt="Settings" className="header-icon" />
               </div>
          </header>
        <h1 className="page-title">Профиль</h1>
        <p className="rtr-loading">Загрузка...</p>
        
      </div>
      <BottomMenu initialSelected={'profile'} />
      </>
    );
  if (!user) return <div className="rtr-error">Профиль не найден</div>;

  return (
    <>
      <div className="profile-page">
         <header className="map-header">
          <div className="logo-container">
        <img src={logoIcon} alt="Logo" style={{ width: 108, height: 25 }} />          </div>
                  <div className="header-icons">
                    <img src={notificationIcon} alt="Notifications" className="header-icon" />
                      <img src={settingsIcon} alt="Settings" className="header-icon" />
               </div>
          </header>

        <h1 className="page-title">Профиль</h1>

        <section className="profile-card">
          <div className="profile-top">
            <img src={avatarImg} alt="Аватар" className="profile-avatar" />
            <div className="profile-info">
              <div className="profile-name">{user.name}</div>
              <div className="profile-location">
                <img src={pinIcon} alt="" className="location-icon" />
                <span>{user.location}</span>
              </div>
              <div className="profile-bottom">
            {user.verified && (
              <span className="verified-badge">
                <img src={verifiedIcon} alt="" className="badge-icon" />
                Верифицирован
              </span>
            )}

            <div className="profile-stats">
              <div className="stat">
                <span className="stat-icon-box">
                  <img src={helpedIcon} alt="" className="stat-icon" />
                </span>
                <span className="stat-text">
                  <span className="stat-num">{user.helped}</span>
                  <div className="stat-label">Помог</div>
                </span>
              </div>
              <div className="stat stat2">
                <span className="stat-icon-box">
                  <img src={requestsIcon} alt="" className="stat-icon" />
                </span>
              <span className="stat-text">
                  <span className="stat-num">{user.requests}</span>
                  <div className="stat-label">Просьб</div>
                </span>
              </div>
            </div>
          </div>
            </div>
          </div>
        </section>

        <section className="thanks-card">
          <div className='thanks-text'>
          <div className="thanks-icon-box">
            <img src={thanksIcon} alt="" className="thanks-icon" />
          </div>
          <h2 className="thanks-title">Спасибо,</h2>
          <p className="thanks-text">что делаете Иннополис более дружным местом.</p>
          </div>
          <div className="thanks-art-box">
            <img src={cityImg} alt="" className="thanks-art" />
          </div>
        </section>

                <nav className="menu-card">
          {/*
          <Link to="/profile/personal" className="menu-item">
            <span className="menu-icon-box">
              <img src={personalIcon} alt="" className="menu-icon" />
            </span>
            <span className="menu-label">Личная информация</span>
            <img src={rightArrowIcon} alt="" className="menu-chevron" />
          </Link>

          <Link to="/notifications" className="menu-item">
            <span className="menu-icon-box">
              <img src={supportIcon} alt="" className="menu-icon" />
            </span>
            <span className="menu-label">Уведомления</span>
            <img src={rightArrowIcon} alt="" className="menu-chevron" />
          </Link>

          <Link to="/support" className="menu-item">
            <span className="menu-icon-box">
              <img src={chevronIcon} alt="" className="menu-icon" />
            </span>
            <span className="menu-label">Помощь и поддержка</span>
            <img src={rightArrowIcon} alt="" className="menu-chevron" />
          </Link>
          */}
          <button className="menu-item" onClick={handleLogout}>
            <span className="menu-icon-box">
              <img src={logoutIcon} alt="" className="menu-icon" />
            </span>
            <span className="menu-label">Выйти из аккаунта</span>
            <img src={rightArrowIcon} alt="" className="menu-chevron" />
          </button>
        </nav>
      </div>
      <BottomMenu initialSelected={'profile'} />
    </>
  );
}

export default Profile;
