import './EventCard.css'
import calendarIcon from './assets/calendar.png';
import peopleIcon from './assets/Profile.png';
import addressIcon from './assets/location.png';
import chatIcon from './assets/message.png';
import closeIcon from './assets/closeIcon.png';
import sportIcon from './assets/gantelActiveIcon.png';
import socialIcon from './assets/socialActiveIcon.png';
import studyIcon from './assets/learningActiveIcon.png';
import { useState, useEffect } from 'react';
import Category from './components/Category';
import { useNavigate } from 'react-router-dom';

const CATEGORY_ITEMS = {
  "Спорт": sportIcon,
  "Соц": socialIcon,
  "Учеба": studyIcon,
};

function EventCard({ eventId, onClose }) {
  const [event, setEvent] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const token = localStorage.getItem('token');

  // Date in such format: day.month, hour:minutes like 02.02, 15:05
  const formatDate = (dateString) => {
    const date = new Date(dateString);  
    const day = date.getDate().toString().padStart(2, '0');
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `${day}.${month}, ${hours}:${minutes}`;
  };

  const formatPeople = (currentPeople, maxPeople) => {
    if (maxPeople) {
      return `${currentPeople}/${maxPeople} мест`;
    } else {
      return `${currentPeople} зарегистрировано`;
    }
  };

  const handleOpenChat = async () => {
    try {
      const response = await fetch(`http://localhost:8080/meetings/${eventId}/chat`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (response.ok) {
        const data = await response.json();
        navigate(`/group-chat/${data.chat_id}`);
      } else {
        console.error('Ошибка при получении чата');
      }
    } catch (error) {
      console.error('Ошибка:', error);
    }
  };

  useEffect(() => {
    if (!eventId) return;

    const fetchEvent = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/meetings/${eventId}`);
        if (!response.ok) throw new Error('Ошибка загрузки');
        const data = await response.json();
        setEvent(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchEvent();
  }, [eventId]);

  if (loading) return (
    <div className="card-page">
      <div className="card-container-loading"></div>
    </div>
  );

  if (error) return (
    <div className="card-page">
      <div className="card-container-error">Ошибка: {error}</div>
    </div>
  );

  if (!event) return null;

  return (
    <>
    <div className="card-page">
        <div className="card-container">
          <div className="card-top">
            <div className="icon-container">
              <img src={CATEGORY_ITEMS[event.type]} alt="" />
            </div>
            <div className="card-right-content">
                {event.type && (
                  <div className="categories-container">
                    <Category label={event.type} />
                  </div>
                )}
              <div className="card-header">
                <h2>{event.title}</h2>
                <div className="card-info-row">
                <div className="card-date">
                  <img src={calendarIcon} alt=""></img>
                  <p>{formatDate(event.meeting_time)}</p>
                </div>
                <div className="card-people-amount">
                  <img src={peopleIcon} alt=""></img>
                  <p>{formatPeople(event.current_people, event.max_people)}</p>
                </div>
                </div>
              </div>
            </div>
            </div>
            <div className="description">
              <p>{event.description}</p>
            </div>
            <div className="address">
              <img src={addressIcon} alt=""></img>
              <h2>{event.address || 'Адрес не указан'}</h2>
            </div>
            <div className="button-container">
              <button className="sign-button">Записаться</button>
              <button className="chat-button"
              onClick={handleOpenChat}>
                <img src={chatIcon} alt="chat" />
              </button>
            </div>
            <button className="close-button" onClick={onClose}><img src={closeIcon} alt="x" /></button>
          </div>
        </div>
    </>
  )
}

export default EventCard
