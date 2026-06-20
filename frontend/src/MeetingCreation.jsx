import './MeetingCreation.css'
import arrowIcon from './assets/arrowIcon.png'
import mapIcon from './assets/mapIcon.png'
import notificationIcon from './assets/notificationIcon.png'
import settingsIcon from './assets/settingsIcon.png'
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import gantelActiveIcon from './assets/gantelActiveIcon.png'
import gantelNotActiveIcon from './assets/gantelNotActiveIcon.png'
import socialActiveIcon from './assets/socialActiveIcon.png'
import socialNotActiveIcon from './assets/socialNotActiveIcon.png'
import learningActiveIcon from './assets/learningActiveIcon.png'
import learningNotActiveIcon from './assets/learningNotActiveIcon.png'
import MapBox from './components/MapBox'
import LocationPicker from './components/LocationPicker'

function MeetingCreation() {
  const navigate = useNavigate();
  const handleClose = () => {
    navigate('/');
  };
  const [selected, setSelected] = useState('Спорт');
  const [text, setText] = useState('');
  const maxLength = 500;
  const [title, setTitle] = useState('');
  const [date, setDate] = useState('');
  const [startTime, setStartTime] = useState('');
  const [endDate, setEndDate] = useState('');
  const [endTime, setEndTime] = useState('');
  const [address, setAddress] = useState('');
  const [maxPeople, setMaxPeople] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [, setError] = useState('');
  const [, setSuccess] = useState('');
  const [latitude, setLatitude] = useState(null);
  const [longitude, setLongitude] = useState(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const categoryIcons = {
    Спорт: {
      active: gantelActiveIcon,
      inactive: gantelNotActiveIcon
    },
    Соц: {
      active: socialActiveIcon,
      inactive: socialNotActiveIcon
    },
    Учеба: {
      active: learningActiveIcon,
      inactive: learningNotActiveIcon
    }
  };

  const combineDateTime = (dateStr, timeStr) => {
    if (!dateStr || !timeStr) return null;
    let formattedTime = timeStr;
    if (timeStr && timeStr.split(':').length === 2) {
      formattedTime = `${timeStr}:00`;
    }
    return `${dateStr}T${formattedTime}Z`;
  };

  const getMeetingTime = () => {
    return combineDateTime(date, startTime);
  };

  const handleLocationSelect = (lng, lat, locationAddress) => {
    if (lng === null && lat === null) {
      setLongitude(null);
      setLatitude(null);
      setAddress('');
    } else {
      setLongitude(lng);
      setLatitude(lat);
      setAddress(locationAddress);
    }
    setIsModalOpen(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setSuccess('');

    if (!title.trim()) {
      setError('Введите название мероприятия');
      setIsLoading(false);
      return;
    }

    if (!text.trim()) {
      setError('Введите описание мероприятия');
      setIsLoading(false);
      return;
    }

    if (!date || !startTime) {
      setError('Укажите дату и время начала мероприятия');
      setIsLoading(false);
      return;
    }

    const meetingTime = getMeetingTime();
    if (!meetingTime) {
      setError('Укажите корректную дату и время');
      setIsLoading(false);
      return;
    }

    const requestBody = {
      title: title.trim(),
      description: text.trim(),
      meeting_time: meetingTime,
      type: selected,
    };

    if (address.trim()) {
      requestBody.address = address.trim();
    }

    if (maxPeople && !isNaN(maxPeople) && parseInt(maxPeople) > 0) {
      requestBody.max_people = parseInt(maxPeople);
    }

    if (latitude != null && longitude != null) {
      requestBody.latitude = latitude;
      requestBody.longitude = longitude;
    }

    try {
      const token = 'temp';
      
      const response = await fetch('http://localhost:8080/meetings', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(requestBody),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.detail || 'Ошибка');
      }
      
      setTitle('');
      setText('');
      setDate('');
      setStartTime('');
      setEndDate('');
      setEndTime('');
      setAddress('');
      setMaxPeople('');
      setSelected('Спорт');
      setLatitude(null);
      setLongitude(null);
      navigate('/');
      
    } catch (err) {
      setError(err.message || 'Ошибка при создании мероприятия. Попробуйте еще раз.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <div className="meeting-creation-page">
        <div className="meeting-creation-page-content">
          <button className="close-button-creation" onClick={handleClose}>×</button>
          <div className="header-top">
            <h2>InnoConnect</h2>
            <button className="notification-button"><img src={notificationIcon} alt="" /></button>
            <button className="settings-button"><img src={settingsIcon} alt="" /></button>
          </div>
          <div className="text-header">
            <h1>Создание мероприятия</h1>
            <p>Расскажите о вашем мероприятии, чтобы другие пользователи могли присоединиться.</p>
          </div>
          <form onSubmit={handleSubmit}>
            <div className="creation-card">
              <div className="meeting-name">
                <h2>Название мероприятия</h2>
                <input 
                  type="text" 
                  placeholder="Введите название мероприятия"
                  value={title}
                  onChange={(e) => setTitle(e.target.value)}
                  required
                />
              </div>
              <div className="meeting-description">
                <h2>Описание мероприятия</h2>
                <textarea 
                  rows={5} 
                  placeholder="Расскажите подробно о вашем мероприятии" 
                  maxLength={maxLength}
                  value={text} 
                  onChange={(e) => setText(e.target.value)}
                  required
                />
                <div className="char-counter-inside">
                  {text.length}/{maxLength}
                </div>
              </div>
              <div className="meeting-time-date">
                <h2>Время мероприятия</h2>
                <p>Дата и время начала</p>
                <div className="choose-date-time-container">
                  <div className="choose-mini-left">
                    <input 
                      type="date" 
                      className="datetime-input"
                      value={date}
                      onChange={(e) => setDate(e.target.value)}
                      required
                    />
                    <input 
                      type="time" 
                      className="datetime-input"
                      value={startTime}
                      onChange={(e) => setStartTime(e.target.value)}
                      required
                    />
                  </div>
                  <img src={arrowIcon} alt=""></img>
                  <div className="choose-mini-right">
                    <input 
                      type="date" 
                      className="datetime-input"
                      value={endDate}
                      onChange={(e) => setEndDate(e.target.value)}
                    />
                    <input 
                      type="time" 
                      className="datetime-input"
                      value={endTime}
                      onChange={(e) => setEndTime(e.target.value)}
                    />
                  </div>
                </div>
              </div>
              <div className="meeting-place">
                <h2>Местоположение</h2>
                <div className="map-container">
                  <MapBox />
                  <div className="bottom-header-map">
                    <img src={mapIcon} alt=""></img>
                    <div className="bottom-header-h-p-combo">
                      <h2>Укажите адрес вашего мероприятия</h2>
                      <p>Нажмите на карту или введите адрес</p>
                    </div>
                    <button type="button" className="edit-location-button"
                        onClick={() => setIsModalOpen(true)}>Изменить</button>
                  </div>
                </div>
                
                <div className="seats-amount">
                  <h2>Количество мест <span className="optional-text">(Optional)</span></h2>
                  <input 
                    type="number" 
                    placeholder="Введите количество мест"
                    value={maxPeople}
                    onChange={(e) => setMaxPeople(e.target.value)}
                    min="1"
                  />
                </div>
                
                <div className="event-cataegories">
                  <h2>Категория мероприятия</h2>
                  <div className="category-switch">
                    <button 
                      type="button"
                      className={selected === 'Спорт' ? 'active' : ''}
                      onClick={() => setSelected('Спорт')}>
                      <img src={selected === 'Спорт' ? categoryIcons.Спорт.active : categoryIcons.Спорт.inactive} alt="" />Спорт
                    </button>
                    <button 
                      type="button"
                      className={selected === 'Соц' ? 'active' : ''}
                      onClick={() => setSelected('Соц')}>
                      <img src={selected === 'Соц' ? categoryIcons.Соц.active : categoryIcons.Соц.inactive} alt="" />Соц
                    </button>
                    <button 
                      type="button"
                      className={selected === 'Учеба' ? 'active' : ''}
                      onClick={() => setSelected('Учеба')}>
                      <img src={selected === 'Учеба' ? categoryIcons.Учеба.active : categoryIcons.Учеба.inactive} alt="" />Учеба
                    </button>
                  </div>
                </div>
                
                <button 
                  type="submit" 
                  className="create-event-button"
                  disabled={isLoading}
                  style={{opacity: isLoading ? 0.7 : 1, cursor: isLoading ? 'not-allowed' : 'pointer'}}>
                  {isLoading ? 'Создание...' : 'Создать мероприятие'}
                </button>
              </div>
            </div>
          </form>
        </div>
      </div>
        <LocationPicker
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onLocationSelect={handleLocationSelect}
        currentLatitude={latitude}
        currentLongitude={longitude}
        currentAddress={address}
        selectedCategory={selected}
      />
    </>
  )
}

export default MeetingCreation
