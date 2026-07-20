import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from './hooks/useAuth';
import './RespondToRequest.css';
import logoIcon from './assets/logo.svg';
import notificationIcon from './assets/notifications.svg';
import settingsIcon from './assets/settings.svg';
import heartIcon from './assets/hearts.svg';
import binIcon from './assets/bin.svg';
import avatarIcon from './assets/mock_ava.svg';
import messIcon from './assets/message.svg';
import locIcon from './assets/location.svg';
import callIcon from './assets/calling.svg';

const RespondToRequest = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const [requestData, setRequestData] = useState(null);
  const [loading, setLoading] = useState(true);
  const token = localStorage.getItem('token');
  const API_URL = import.meta.env.VITE_API_URL;
  const isAuthenticated = useAuth();

  useEffect(() => {
    const fetchRequest = async () => {
      try {
        const response = await fetch(`${API_URL}/requests/${id || 1}`);
        if (response.ok) {
          const data = await response.json();
          setRequestData(data);
        } else {
          useMockData();
        }
      } catch (error) {
        useMockData();
      } finally {
        setLoading(false);
      }
    };
    fetchRequest();
  }, [id, API_URL]);

  const useMockData = () => {
    setRequestData({
      request_id: id,
      creator: {
        id: 123,
        name: 'Иван Петров'
      },
      title: 'Помочь с выносом мусора',
      description: 'Здравствуйте! Нужно вынести мусор до помойки. Тяжелых пакетов 2-3 штуки. Бейбифокс в подарок.',
      requester_address: 'ул. Спортивная 128, Иннополис',
      type: 'Помощь',
      deadline: 'Сегодня, 15:45',
    });
  };

  const handleOpenChat = async () => {
    if (!isAuthenticated) {
      navigate('/welcome');
      return;
    }
    try {
      const response = await fetch(`${API_URL}/requests/${id}/chat`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (response.ok) {
        const data = await response.json();
        navigate(`/chat/${data.chat_id}`);
      } else {
        console.error('Ошибка при создании/получении чата');
      }
    } catch (error) {
      console.error('Ошибка:', error);
    }
  };

  const handleResponse = async () => {
    if (!isAuthenticated) {
      navigate('/welcome');
      return;
    }

    try {
      const response = await fetch(`${API_URL}/requests/${id}/apply`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.status === 200) {
        console.log("Successfully applied to request");

      } else if (response.status === 401) {
        localStorage.removeItem('token');
        navigate('/welcome');

      } else {
        const errorData = await response.json();
        console.error(
          errorData.message || "Failed to apply"
        );
      }

    } catch (error) {
      console.error("Error applying to request:", error);
    }
  };

  const typeOfReq = {
        "Помощь": "helpreq",
        "Вещи": "stuffreq",
        "Транспорт": "transportreq",
        "Прочее": "otherreq",
    }

  if (loading) return <div className="rtr-loading">Загрузка...</div>;
  if (!requestData) return <div className="rtr-error">Запрос не найден</div>;

  return (
    <div className="respond-page-wrapper">
       <header className="map-header">
          <div className="logo-container">
            <button className="rtr-back-btn" onClick={() => navigate('/chats')}>
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path d="M15 18L9 12L15 6" stroke="#1A1D1E" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    </svg>
  </button>
<img src={logoIcon} alt="Logo" style={{ width: 108, height: 25 }} />          </div>
          {/* <div className="header-icons">
            <img src={notificationIcon} alt="Notifications" className="header-icon" />
              <img src={settingsIcon} alt="Settings" className="header-icon" />
          </div> */}
        </header>

      <div className="rtr-card main-info-card">
        <div className="rtr-icon-box">
         <div className={`type-request-icon-card ${typeOfReq[requestData.type]}`}></div>
        </div>
        <div className="rtr-main-details">
          <span className="rtr-category">{requestData.type}</span>
          <h1 className="rtr-title">{requestData.title}</h1>
          <span className="rtr-date">{requestData.deadline}</span>
        </div>
        <div className="rtr-status-badge">Новое</div>
      </div>

      <div className="rtr-card author-card-full">
        <div className="rtr-author-left">
          <div className="rtr-avatar">
         <img src={avatarIcon} alt="Avatar" className="rtr-iconbox-icon" />
          </div>
          <div className="rtr-author-info">
            <div className="rtr-author-name-row">
              <p className="rtr-name">{requestData.creator.name}</p>
              <p className="rtr-verified">Верифицирован</p>
            </div>
            <span className="rtr-author-sub">Житель Иннополиса</span>
          </div>
        </div>
        <button className="rtr-call-btn">
          <img src={callIcon} alt="Call" className="rtr-iconbox-icon" />
        </button>
      </div>

      <div className="rtr-card desc-card">
        <h3 className="rtr-label">Описание</h3>
        <p className="rtr-desc-text">{requestData.description}</p>
      </div>

      <div className="rtr-card loc-card">
        <h3 className="rtr-label">Куда принести</h3>
        <div className="rtr-loc-row">
         <img src={locIcon} alt="Bin" className="rtr-iconbox-icon" />
          <div className="rtr-loc-details">
            <span className="rtr-address">{requestData.requester_address}</span>
            <span className="rtr-entrance">Подъезд 1</span>
          </div>
        </div>
      </div>

      <div className="rtr-bottom-actions">
        <button className="rtr-btn-primary">
         <img src={heartIcon} alt="Heart" className="rtr-iconbox-icon" onClick={handleResponse} />
          Откликнуться
        </button>
        <button className="rtr-btn-secondary"
        onClick={handleOpenChat}>
         <img src={messIcon} alt="Mess" className="rtr-iconbox-icon" />
          Написать автору
        </button>
      </div>
    </div>
  );
};

export default RespondToRequest;
