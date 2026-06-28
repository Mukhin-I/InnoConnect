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

function Map() {
  const [meetings, setMeetings] = useState([]);
  const [selectedMeetingId, setSelectedMeetingId] = useState(null);


  const [activeFilter, setActiveFilter] = useState('Спорт');
  const [highlightStyle, setHighlightStyle] = useState({ opacity: 0 }); // Start hidden until measured
  const containerRef = useRef(null);

  const handleFilterClick = (filter) => {
    setActiveFilter(filter);
  };

  useEffect(() => {
    async function fetchMeetings() {
      try {
        const response = await fetch("http://localhost:8080/meetings");

        if (!response.ok) {
          throw new Error("Failed to fetch meetings");
        }

        const data = await response.json();

        setMeetings(data.meetings);
      } catch (error) {
        console.error(error);
      }
    }

    fetchMeetings();
  }, []);


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
            <img src={logoIcon} alt="Logo" className="logo-icon" />
          </div>
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

        <div className="map-container-placeholder">
          <MapBox 
            meetings={meetings}
            selectedMeetingId={selectedMeetingId}
            setSelectedMeetingId={setSelectedMeetingId}
          />

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