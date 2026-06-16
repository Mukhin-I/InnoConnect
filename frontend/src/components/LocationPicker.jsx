import './LocationPicker.css'
import { useState } from 'react';
import MapBox from './MapBox';

function LocationPicker({ 
  isOpen, 
  onClose, 
  onLocationSelect,
  currentLatitude,
  currentLongitude,
  currentAddress 
}) {
  const [selectedPoint, setSelectedPoint] = useState(
    currentLatitude && currentLongitude 
      ? { lng: currentLongitude, lat: currentLatitude }
      : null
  );
  const [address, setAddress] = useState(currentAddress || '');
  const [isLoadingAddress, setIsLoadingAddress] = useState(false);
  
  const getAddressFromCoordinates = async (lng, lat) => {
    setIsLoadingAddress(true);
    try {
      const response = await fetch(
        `https://api.mapbox.com/geocoding/v5/mapbox.places/${lng},${lat}.json?access_token=${import.meta.env.VITE_MAPBOX_TOKEN}&language=ru`
      );
      const data = await response.json();
      if (data.features && data.features.length > 0) {
        setAddress(data.features[0].place_name);
      } else {
        setAddress('Адрес не найден');
      }
    } catch (error) {
      setAddress('Ошибка получения адреса', error);
    } finally {
      setIsLoadingAddress(false);
    }
  };

  const handleMapLocationSelect = async (lng, lat) => {
    setSelectedPoint({ lng, lat });
    await getAddressFromCoordinates(lng, lat);
  };

  const handleClearLocation = () => {
    setSelectedPoint(null);
    setAddress('');
  };

  const handleSave = () => {
    if (selectedPoint) {
      onLocationSelect(selectedPoint.lng, selectedPoint.lat, address);
    }
    onClose();
  };

  if (!isOpen) return null;

  return (
    <div className="modal-address" onClick={onClose}>
      <div className="modal-address-content" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>Выберите местоположение на карте</h2>
          <button className="close-button" onClick={onClose}>×</button>
        </div>
        <div className="modal-map-container">
          <MapBox 
            meetings={[]}
            selectedMeetingId={null}
            setSelectedMeetingId={() => {}}
            center={currentLatitude && currentLongitude 
              ? [currentLongitude, currentLatitude]
              : [48.752, 55.752]
            }
            zoom={13}
            isSelectable={true}
            onLocationSelect={handleMapLocationSelect}
            initialLatitude={selectedPoint?.lat || currentLatitude}
            initialLongitude={selectedPoint?.lng || currentLongitude}
            clearMarker={!selectedPoint}
          />
        </div>
        <div className="modal-address-container">
        <div className="address-row">
            <strong className="address-label">Адрес: </strong>
            <span className="address-value">
            {isLoadingAddress ? 'Загрузка...' : (address || 'Не выбран')}
            </span>
        </div>
        </div>
          
          <div className="modal-buttons">
            <button className="clear-button" onClick={handleClearLocation}>
              Очистить
            </button>
            <button className="cancel-button" onClick={onClose}>
              Отмена
            </button>
            <button 
              className="save-button" 
              onClick={handleSave}
              disabled={!selectedPoint}
            >
              Сохранить
            </button>
          </div>
        </div>
      </div>
  );
}

export default LocationPicker;