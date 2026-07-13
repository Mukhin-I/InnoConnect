import { Link } from 'react-router-dom';
import './Map.css';
import { useState, useEffect, useRef } from 'react';
import logoIcon from './assets/logo.svg';
import filterIcon from './assets/filter.svg';
import sportIcon from './assets/football.svg';
import coffeeIcon from './assets/coffee-cup.svg';
import studyIcon from './assets/books.svg';
import notificationIcon from './assets/notifications.svg';
import settingsIcon from './assets/settings.svg';
import MapBox from './components/MapBox';
import EventCard from './EventCard.jsx';
import BottomMenu from './components/BottomMenu.jsx';
import CardOfMeeting from './components/CardOfMeeting.jsx';
import mapIcon from './assets/map-toggle.svg';
import listIcon from './assets/list-toggle.svg';

const tomorrowAt = (hours, minutes = 0) => {
  const d = new Date();
  d.setDate(d.getDate() + 1);
  d.setHours(hours, minutes, 0, 0);
  return d.toISOString();
};

const MOCK_MEETINGS = [
  {
    id: 'mock-1',
    title: 'Утренняя йога в парке',
    type: 'Спорт',
    meeting_time: tomorrowAt(8, 45),
    address: 'Парк у Артспейса',
    latitude: 55.7522,
    longitude: 48.7444,
    current_people: 4,
    max_people: 10,
  },
  {
    id: 'mock-2',
    title: 'Пробежка вокруг озера',
    type: 'Спорт',
    meeting_time: tomorrowAt(19, 0),
    address: 'Озеро у Технопарка',
    latitude: 55.7530,
    longitude: 48.7450,
    current_people: 3,
    max_people: 8,
  },
  {
    id: 'mock-3',
    title: 'Кофе и настолки',
    type: 'Соц',
    meeting_time: tomorrowAt(18, 30),
    address: 'Кофейня «Пар»',
    latitude: 55.7515,
    longitude: 48.7430,
    current_people: 5,
    max_people: 12,
  },
  {
    id: 'mock-4',
    title: 'Языковой обмен EN/RU',
    type: 'Соц',
    meeting_time: tomorrowAt(20, 0),
    address: 'Лобби Технопарка',
    latitude: 55.7519,
    longitude: 48.7460,
    current_people: 6,
    max_people: 15,
  },
  {
    id: 'mock-5',
    title: 'Learning club: System Design',
    type: 'Учеба',
    meeting_time: tomorrowAt(17, 0),
    address: 'Университет Иннополис, ауд. 108',
    latitude: 55.7540,
    longitude: 48.7455,
    current_people: 8,
    max_people: 20,
  },
  {
    id: 'mock-6',
    title: 'Разбор задач по алгоритмам',
    type: 'Учеба',
    meeting_time: tomorrowAt(21, 0),
    address: 'Библиотека кампуса',
    latitude: 55.7548,
    longitude: 48.7440,
    current_people: 2,
    max_people: 6,
  },
];

function Map() {
  const [meetings, setMeetings] = useState(MOCK_MEETINGS);
  const [selectedMeetingId, setSelectedMeetingId] = useState(null);


  const [activeFilter, setActiveFilter] = useState('Спорт');
  const [viewMode, setViewMode] = useState('map');
  const filteredMeetings = (meetings || []).filter(
    meeting => meeting.type === activeFilter
  );
  const [highlightStyle, setHighlightStyle] = useState({ opacity: 0 }); // Start hidden until measured
  const containerRef = useRef(null);

  const API_URL = import.meta.env.VITE_API_URL;

  const handleFilterClick = (filter) => {
    setActiveFilter(filter);
    setSelectedMeetingId(null);
  };

  useEffect(() => {
    async function fetchMeetings() {
      try {
        const response = await fetch(`${API_URL}/meetings`);

        if (!response.ok) {
          throw new Error("Failed to fetch meetings");
        }

        const data = await response.json();

        if (data.meetings && data.meetings.length > 0) {
          setMeetings(data.meetings);
        }
      } catch (error) {
        console.error(error);
      }
    }

    fetchMeetings();
  }, [API_URL]);


  useEffect(() => {
    if (containerRef.current) {
      const activeBtn = containerRef.current.querySelector('.active-filter');
      if (activeBtn) {
        setHighlightStyle({
          left: activeBtn.offsetLeft,
          width: activeBtn.offsetWidth,
          height: activeBtn.offsetHeight,
          top: activeBtn.offsetTop,
          opacity: 1,
        });
      }
    }
  }, [activeFilter]);

  return (
    <>
      <div className="map-page">
        <header className="map-header">
          <div className="logo-container">
<img src={logoIcon} alt="Logo" style={{ width: 108, height: 25 }} />          </div>
          <div className="header-icons">
            <img src={notificationIcon} alt="Notifications" className="header-icon" />
              <img src={settingsIcon} alt="Settings" className="header-icon" />
          </div>
        </header>

        <div className="filters-bar">
          <div className="filters-group" ref={containerRef}>
            <div className="highlight-pill" style={highlightStyle}></div>
            <button className="filter-btn filter-settings">
              <img src={filterIcon} alt="Filter" className="filter-icon" />
              Фильтры
            </button>
          
            <button className={`filter-btn category-btn ${activeFilter === 'Спорт' ? 'active-filter' : ''}`} onClick={() => handleFilterClick('Спорт')}>
              <img src={sportIcon} alt="Sport" className="filter-icon" />
              Спорт
            </button>
            <button className={`filter-btn category-btn ${activeFilter === 'Соц' ? 'active-filter' : ''}`} onClick={() => handleFilterClick('Соц')}>
              <img src={coffeeIcon} alt="Soc" className="filter-icon" />
              Соц
            </button>
            <button className={`filter-btn category-btn ${activeFilter === 'Учеба' ? 'active-filter' : ''}`} onClick={() => handleFilterClick('Учеба')}>
              <img src={studyIcon} alt="Study" className="filter-icon" />
              Учеба
            </button>
          </div>
      
          <Link to="/create-meeting" className="add-btn">
            <span>+</span>
          </Link>
        </div>

        <div className="view-toggle">
          <div
            className="view-toggle-pill"
            style={{ transform: viewMode === 'list' ? 'translateX(100%)' : 'translateX(0)' }}
          />
          <button
            type="button"
            className={`view-toggle-btn ${viewMode === 'map' ? 'active' : ''}`}
            onClick={() => setViewMode('map')}
          >
            <img src={mapIcon} alt="" />
            Карта
          </button>
          <button
            type="button"
            className={`view-toggle-btn ${viewMode === 'list' ? 'active' : ''}`}
            onClick={() => setViewMode('list')}
          >
            <img src={listIcon} alt="" />
            Список
          </button>
        </div>

        <div className="map-container-placeholder">
          {viewMode === 'map' ? (
            <MapBox
              meetings={filteredMeetings}
              selectedMeetingId={selectedMeetingId}
              setSelectedMeetingId={setSelectedMeetingId}
            />
          ) : (
            <div className="list-of-meetings-container">
              {filteredMeetings.map((meeting) => (
                <CardOfMeeting
                  key={meeting.id}
                  meeting={meeting}
                  onClick={setSelectedMeetingId}
                />
              ))}
            </div>
          )}

          {selectedMeetingId !== null && (
            <EventCard
              eventId={selectedMeetingId}
              onClose={() => setSelectedMeetingId(null)}
            />
          )}
        </div>
      </div>
      <BottomMenu initialSelected={'map'} />
    </>
  );
}

export default Map;