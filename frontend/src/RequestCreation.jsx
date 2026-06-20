import './RequestCreation.css'
import mapIcon from './assets/mapIcon.png'
import notificationIcon from './assets/notificationIcon.png'
import settingsIcon from './assets/settingsIcon.png'
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import LocationPicker from './components/LocationPicker'
import helpActiveIcon from './assets/helpActiveIcon.png'
import helpNotActiveIcon from './assets/helpNotActiveIcon.png'
import thingsActiveIcon from './assets/thingsActiveIcon.png'
import thingsNotActiveIcon from './assets/thingsNotActiveIcon.png'
import carActiveIcon from './assets/carActiveIcon.png'
import carNotActiveIcon from './assets/carNotActiveIcon.png'
import otherActiveIcon from './assets/otherActiveIcon.png'
import otherNotActiveIcon from './assets/otherNotActiveIcon.png'
import IconButton from './components/IconButton';
import CreationInput from './components/CreationInput';
import pencilIcon from './assets/pencil.png';
import TextareaCounter from './components/TextareaCounter'
import LocationEditButton from './components/LocationEditButton';
import CategorySelector from './components/CategorySelector';
import CreateButton from './components/CreateButton';

function RequestCreation() {
  const navigate = useNavigate();
  const handleClose = () => {
    navigate('/');
  };
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [selected, setSelected] = useState('Помощь');
  const [address, setAddress] = useState('');
  const [endDate, setEndDate] = useState('');
  const [endTime, setEndTime] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [, setError] = useState('');
  const [, setSuccess] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);

  const categories = [
    {
      id: 'Помощь',
      label: 'Помощь',
      activeIcon: helpActiveIcon,
      inactiveIcon: helpNotActiveIcon,
    },
    {
      id: 'Вещи',
      label: 'Вещи',
      activeIcon: thingsActiveIcon,
      inactiveIcon: thingsNotActiveIcon,
    },
    {
      id: 'Транспорт',
      label: 'Транспорт',
      activeIcon: carActiveIcon,
      inactiveIcon: carNotActiveIcon,
    },
    {
      id: 'Прочее',
      label: 'Прочее',
      activeIcon: otherActiveIcon,
      inactiveIcon: otherNotActiveIcon,
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

  const getRequestTime = () => {
    return combineDateTime(endDate, endTime);
  };

  const formatAddress = (fullAddress) => {
    if (!fullAddress) return '';
    const parts = fullAddress.split(',').map(part => part.trim());
    let streetAndHouse = '';
    let city = '';
    for (let i = 0; i < parts.length; i++) {
      const part = parts[i];
      if (part.match(/ул\.|просп\.|пер\.|пл\.|бульв\.|шоссе|наб\.|аллея/i)) {
        streetAndHouse = part;
        if (i + 1 < parts.length && parts[i + 1].match(/\d/)) {
          streetAndHouse += ', ' + parts[i + 1];
        }
        break;
      }
    }
    if (!streetAndHouse) {
      for (let i = 0; i < parts.length; i++) {
        const part = parts[i];
        if (part.match(/\d/) && !part.match(/^\d+$/)) {
          streetAndHouse = part;
          break;
        }
      }
    }
    
    for (let i = parts.length - 1; i >= 0; i--) {
      const part = parts[i];
      if (part.match(/г\.|город/i)) {
        city = part;
        break;
      }
    }
    
    if (!city && streetAndHouse) {
      for (let i = parts.length - 1; i >= 0; i--) {
        const part = parts[i];
        if (part !== streetAndHouse && !part.match(/\d/) && part.length > 2) {
          city = part;
          break;
        }
      }
    }
    
    let formattedAddress = '';
    if (streetAndHouse) {
      formattedAddress = streetAndHouse;
    }
    if (city) {
      if (formattedAddress) {
        formattedAddress += ', ' + city;
      } else {
        formattedAddress = city;
      }
    }
    
    if (!formattedAddress) {
      const firstParts = parts.slice(0, 3).join(', ');
      return firstParts || fullAddress;
    }
    
    return formattedAddress;
  };

  const handleLocationSelect = (lng, lat, locationAddress) => {
    if (lng === null && lat === null) {
      setAddress('');
    } else {
      const formattedAddr = formatAddress(locationAddress);
      setAddress(formattedAddr);
    }
    setIsModalOpen(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setSuccess('');

    if (!title.trim()) {
      setError('Введите название просьбы');
      setIsLoading(false);
      return;
    }

    if (!description.trim()) {
      setError('Введите описание просьбы');
      setIsLoading(false);
      return;
    }

    if (!endDate || !endTime) {
      setError('Укажите крайний срок (дедлайн) просьбы');
      setIsLoading(false);
      return;
    }

    const requestTime = getRequestTime();
    if (!requestTime) {
      setError('Укажите корректную дату и время');
      setIsLoading(false);
      return;
    }

    const requestBody = {
      title: title.trim(),
      description: description.trim(),
      requester_address: address.trim(),
      type: selected,
      deadline: requestTime,
    };

    try {
      const token = 'temp';
      
      const response = await fetch('http://localhost:8080/requests', {
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
      setAddress('');
      setEndDate('');
      setEndTime('');
      setSelected('Помощь');
      navigate('/');
      
    } catch (err) {
      setError(err.message || 'Ошибка при создании просьбы. Попробуйте еще раз.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <div className="request-creation-page">
        <div className="request-creation-page-content">
            <button className="close-button-creation" onClick={handleClose}>×</button>
            <div className="header-top">
              <h2>InnoConnect</h2>
              <IconButton 
                icon={notificationIcon} 
                alt=""
              />
              <IconButton 
                icon={settingsIcon} 
                alt=""
              />
            </div>
            <div className="text-header">
              <h1>Создание просьбы</h1>
              <p>Расскажите соседям, какая помощь вам нужна.</p>
            </div>
            <form onSubmit={handleSubmit}>
            <div className="creation-card">
              <div className="request-name">
                <h2>Название просьбы</h2>
                <CreationInput
                    value={title}
                    onChange={(e) => setTitle(e.target.value)}
                    placeholder="Введите название просьбы"
                    icon={pencilIcon}
                    required
                />
              </div>
                <div className="request-description">
                <h2>Описание просьбы</h2>
                <TextareaCounter
                   value={description}
                   onChange={(e) => setDescription(e.target.value)}
                   maxLength={500}
                   rows={5}
                   placeholder="Расскажите подробнее, какая помощь нужна, что важно учесть"
                   required
                />
            </div>
            <div className="request-address-container">
                <h2>Куда принести</h2>
                <div className="address-card">
                    <div className="icon-address-container">
                    <img src={mapIcon} alt=""></img>
                    </div>
                    <div className="address">
                        <h2>{address}</h2>
                    </div>
                    <LocationEditButton onClick={() => setIsModalOpen(true)} />
                </div>
            </div>
            <div className="deadline-container">
                <h2>Крайний срок <span className="optional-text">(Дедлайн)</span></h2>
                <div className="deadline-choose">
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
            <div className="request-catergories">
                <h2>Категория просьбы</h2>
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
                  {isLoading ? 'Создание...' : 'Создать просьбу'}
              </CreateButton>
        </div>
        </form>
      </div>
      </div>
        <LocationPicker
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onLocationSelect={handleLocationSelect}
        currentAddress={address}
      />
    </>
  )
}

export default RequestCreation
