import './EventCard.css'
import calendarIcon from './assets/calendar.png';
import peopleIcon from './assets/Profile.png';
import addressIcon from './assets/location.png';
import chatIcon from './assets/message.png';
import closeIcon from './assets/closeIcon.png';
import arrowIcon from './assets/right-arrow.svg';
import avatarIcon from './assets/mock_ava.svg';
import participantsIcon from './assets/people.png';
import sportIcon from './assets/gantelActiveIcon.png';
import socialIcon from './assets/socialActiveIcon.png';
import studyIcon from './assets/learningActiveIcon.png';
import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from './hooks/useAuth';

const CATEGORY_ITEMS = {
  "Спорт": sportIcon,
  "Соц": socialIcon,
  "Учеба": studyIcon,
};

// Вверху файла, после импортов
const mockEvent = {
  title: "Футбол на свежем воздухе",
  type: "Спорт",
  meeting_time: "2026-07-15T15:30:00",
  current_people: 7,
  max_people: 12,
  description: "Приходите поиграть в футбол! Берите с собой воду и хорошее настроение. Уровень подготовки не важен.",
  address: "Спортивный комплекс ",
};

const mockParticipants = [
  {
    id: 1,
    name: 'Иван Петров',
    subtitle: 'Организатор',
    isOrganizer: true,
  },
  {
    id: 2,
    name: 'Марина Смирнова',
    subtitle: 'Житель Иннополиса',
    isOrganizer: false,
  },
  {
    id: 3,
    name: 'Дмитрий Васильев',
    subtitle: 'Житель Иннополиса',
    isOrganizer: false,
  },
];

function EventCard({ eventId, onClose }) {
  // const [event, setEvent] = useState(null);
  // const [loading, setLoading] = useState(true);
  // const [error, setError] = useState(null);
  const [event, setEvent] = useState(mockEvent); // <- сразу мок
const [loading, setLoading] = useState(false);  // <- false чтобы не крутился лоадер
const [error, setError] = useState(null);

  const navigate = useNavigate();
  const token = localStorage.getItem('token');
  const API_URL = import.meta.env.VITE_API_URL;
  const isAuthenticated = useAuth();

  const formatTime = (date) => {
    const hours = date.getHours();
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `${hours}:${minutes}`;
  };

  const formatDate = (dateString) => {
    const start = new Date(dateString);
    const end = new Date(start.getTime() + 90 * 60 * 1000);
    const now = new Date();
    const isToday = start.toDateString() === now.toDateString();
    const prefix = isToday
      ? 'Сегодня'
      : `${start.getDate().toString().padStart(2, '0')}.${(start.getMonth() + 1)
          .toString()
          .padStart(2, '0')}`;

    return `${prefix}, ${formatTime(start)} - ${formatTime(end)}`;
  };

  const formatPeople = (currentPeople, maxPeople) => {
    if (maxPeople) {
      return `${currentPeople}/${maxPeople} мест`;
    } else {
      return `${currentPeople} зарегистрировано`;
    }
  };

  const handleOpenChat = async () => {
    if (!isAuthenticated) {
      navigate('/welcome');
      return;
    }
    try {
      const response = await fetch(`${API_URL}/meetings/${eventId}/chat`, {
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

  const handleSign = async () => {
    if (!isAuthenticated) {
      navigate('/welcome');
      return;
    }

    try {
      const response = await fetch(`${API_URL}/meetings/${eventId}`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (response.status === 200) {
        const updatedEvent = await fetch(`${API_URL}/meetings/${eventId}`);
        if (updatedEvent.ok) {
          const data = await updatedEvent.json();
          setEvent(data);
        }
      } else if (response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('tokenType');
        localStorage.removeItem('expiresIn');
        navigate('/welcome');
      } else {
        const errorData = await response.json();
        console.error(`Ошибка: ${errorData.message || 'Не удалось записаться'}`);
      }
    } catch (error) {
      console.error('Ошибка при записи:', error);
    }
  }

  const participants = event.participants?.length
    ? event.participants.map((participant, index) => ({
        id: participant.id ?? index,
        name:
          participant.name ||
          `${participant.first_name || ''} ${participant.last_name || ''}`.trim() ||
          'Участник',
        subtitle: participant.subtitle || participant.role || 'Житель Иннополиса',
        isOrganizer: Boolean(participant.is_organizer) || index === 0,
      }))
    : mockParticipants;

  const participantsCount = event.max_people
    ? `${event.current_people}/${event.max_people}`
    : `${event.current_people || participants.length}`;

  const distanceLabel = event.distance ? `${event.distance}` : '150м от вас';

  useEffect(() => {
    if (!eventId) return;

    const fetchEvent = async () => {
      setLoading(true);
      try {
        const response = await fetch(`${API_URL}/meetings/${eventId}`);
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
  }, [eventId, API_URL]);

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
          <button className="close-button" onClick={onClose}><img src={closeIcon} alt="x" /></button>

          <div className="card-top">
            <div className="icon-container">
              <img src={CATEGORY_ITEMS[event.type]} alt="" />
            </div>
            <div className="card-right-content">
              {event.type && (
                <div className="categories-container">
                  <span className="event-type">{event.type}</span>
                </div>
              )}
              <div className="card-header">
                <h2 className="card-header-title">{event.title}</h2>
                <div className="card-info-row">
                  <div className="card-date">
                    <img src={calendarIcon} alt=""></img>
                    <p className='card-formatted-date'>{formatDate(event.meeting_time)}</p>
                  </div>
                  <div className="card-people-amount">
                    <img src={peopleIcon} alt=""></img>
                    <p className='card-formatted-date'>{formatPeople(event.current_people, event.max_people)}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="description">
            <p className='event-description'>{event.description}</p>
          </div>

          <div className="address-card-event">
            <div className="address-left">
              <img src={addressIcon} alt=""></img>
              <p className="address-description">{event.address || 'Адрес не указан'}</p>
            </div>
            <p className="address-distance">{distanceLabel}</p>
          </div>

          <div className="participants-section">
            <div className="participants-header">
              <div className="participants-title-wrap">
                <img src={participantsIcon} alt="Участники" />
                <p className="address-description">Участники</p>
                <span>{participantsCount}</span>
              </div>
              <img src={arrowIcon} alt="Свернуть" className="participants-arrow" />
            </div>

            <div className="participants-list">
              {participants.map((participant) => (
                <div className="participant-row" key={participant.id}>
                  <div className="participant-left">
                    <img src={avatarIcon} alt={participant.name} className="participant-avatar" />
                    <div className="participant-meta">
                      <p className="participant-name-event">{participant.name}</p>
                      <p className="participant-subtitle">{participant.subtitle}</p>
                    </div>
                  </div>
                      {participant.isOrganizer && <span className="organizer-chip">Организатор</span>}
                </div>
              ))}
            </div>
          </div>

          <div className="button-container">
            <button className="sign-button" onClick={handleSign}>Записаться</button>
            <button className="chat-button" onClick={handleOpenChat}>
              <img src={chatIcon} alt="chat" />
            </button>
          </div>
        </div>
      </div>
    </>
  )
}

export default EventCard
