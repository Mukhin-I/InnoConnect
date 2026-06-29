import './MeetingCreation.css'
import arrowIcon from './assets/arrowIcon.png'
import mapIcon from './assets/mapIcon.png'
import notificationIcon from './assets/notificationIcon.png'
import calendarIcon from './assets/calendarIcon.png'
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
import IconButton from './components/IconButton'
import CreationInput from './components/CreationInput'
import TextareaCounter from './components/TextareaCounter'
import LocationEditButton from './components/LocationEditButton'
import seatsIcon from './assets/seatsIcon.png'
import CategorySelector from './components/CategorySelector'
import CreateButton from './components/CreateButton'
import logoIcon from './assets/logo.svg';


function MeetingCreation() {
  const navigate = useNavigate();
  const handleClose = () => {
    navigate('/');
  };
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [selected, setSelected] = useState('Спорт');
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

  const categories = [
    {
      id: 'Спорт',
      label: 'Спорт',
      activeIcon: gantelActiveIcon,
      inactiveIcon: gantelNotActiveIcon
    },
    {
      id: 'Соц',
      label: 'Соц',
      activeIcon: socialActiveIcon,
      inactiveIcon: socialNotActiveIcon
    },
    {
      id: 'Учеба',
      label: 'Учеба',
      activeIcon: learningActiveIcon,
      inactiveIcon: learningNotActiveIcon
    },
  ];

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

    if (!description.trim()) {
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
      description: description.trim(),
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
      setDescription('');
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
                <header className="chat-header">
                        <div className="header-left">
                        <button className="rtr-back-btn" onClick={() => navigate('/chats')}>
                            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                              <path d="M15 18L9 12L15 6" stroke="#1A1D1E" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                            </svg>
                          </button>
                          <div className="logo-container">
                            <img src={logoIcon} alt="Logo" className="logo-icon" />
                          </div>
                        </div>
                        <div className="header-icons">
                          <img src={notificationIcon} alt="Notifications" className="header-icon" />
                          <img src={settingsIcon} alt="Settings" className="header-icon" />
                        </div>
                      </header>
        <div className="meeting-creation-page-content">
          <div className="text-header">
            <h1>Создание мероприятия</h1>
            <p className="text-header-p">Расскажите о вашем мероприятии, чтобы другие пользователи могли присоединиться.</p>
          </div>
          <form onSubmit={handleSubmit}>
            <div className="creation-card">
              <div className="meeting-name">
                <h2>Название мероприятия</h2>
                <CreationInput
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Введите название мероприятия"
                icon={calendarIcon}
                required
                />
              </div>
              <div className="meeting-description">
                <h2>Описание мероприятия</h2>
                <TextareaCounter
                   value={description}
                   onChange={(e) => setDescription(e.target.value)}
                   maxLength={500}
                   rows={5}
                   placeholder="Расскажите подробно о вашем мероприятии"
                   required
                />
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
                    <LocationEditButton onClick={() => setIsModalOpen(true)} />
                  </div>
                </div>
                
                <div className="seats-amount">
                  <h2>Количество мест <span className="optional-text">(Optional)</span></h2>
                  <CreationInput
                      type="number"
                      placeholder="Введите количество мест"
                      value={maxPeople}
                      onChange={(e) => setMaxPeople(e.target.value)}
                      min="1"
                      icon={seatsIcon}
                  />
                </div>
                
                <div className="event-cataegories">
                  <h2>Категория мероприятия</h2>
                  <CategorySelector
                    categories={categories}
                    selectedCategory={selected}
                    onSelectCategory={setSelected}
                  />
                </div>
                <CreateButton
                  type="submit"
                  disabled={isLoading}
                  style={{
                      opacity: isLoading ? 0.7 : 1,
                      cursor: isLoading ? 'not-allowed' : 'pointer'
                  }}
              >
                  {isLoading ? 'Создание...' : 'Создать мероприятие'}
              </CreateButton>
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
