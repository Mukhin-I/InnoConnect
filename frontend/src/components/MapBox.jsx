import './MapBox.css'
import { useEffect, useRef } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import pinIconInactive from '../assets/pin-sport-inactive.svg';
import pinIconActive from '../assets/pin-sport-active.svg';

function MapBox({
  meetings = [],
  selectedMeetingId,
  setSelectedMeetingId,
  center = [48.752, 55.752],
  zoom = 12,
  style = 'mapbox://styles/seanshushickkk/cmqazkvqy003d01qz36xq2coo',
  isSelectable = false,
  onLocationSelect = null,
  initialLatitude = null,
  initialLongitude = null,
  clearMarker = false
}) {
  const mapContainerRef = useRef(null);
  const mapRef = useRef(null);
  const selectMarkerRef = useRef(null);
  const isMapLoadedRef = useRef(false);

  const addSelectMarker = (lng, lat) => {
    if (!mapRef.current) return;
    
    if (selectMarkerRef.current) {
      selectMarkerRef.current.remove();
    }
    
    const el = document.createElement("img");
    el.src = pinIconActive;
    el.style.width = "34px";
    el.style.height = "34px";
    el.style.cursor = "pointer";
    
    selectMarkerRef.current = new mapboxgl.Marker(el)
      .setLngLat([lng, lat])
      .addTo(mapRef.current);
  };

  const removeSelectMarker = () => {
    if (selectMarkerRef.current) {
      selectMarkerRef.current.remove();
      selectMarkerRef.current = null;
    }
  };

  useEffect(() => {
    if (!mapContainerRef.current) return;

    mapboxgl.accessToken = import.meta.env.VITE_MAPBOX_TOKEN;

    mapRef.current = new mapboxgl.Map({
      container: mapContainerRef.current,
      style,
      center,
      zoom,
      attributionControl: false,
    });

    mapRef.current.on('load', () => {
      isMapLoadedRef.current = true;
      
      if (initialLatitude && initialLongitude && !clearMarker) {
        addSelectMarker(initialLongitude, initialLatitude);
      }
    });

    mapRef.current.on('click', (e) => {
      const { lng, lat } = e.lngLat;
      
      if (isSelectable && onLocationSelect) {
        addSelectMarker(lng, lat);
        onLocationSelect(lng, lat);
      }
    });

    return () => {
      if (selectMarkerRef.current) {
        selectMarkerRef.current.remove();
      }
      if (mapRef.current) {
        mapRef.current.remove();
      }
    };
  }, []);

  useEffect(() => {
    if (!mapRef.current || !isMapLoadedRef.current) return;
    
    if (clearMarker) {
      removeSelectMarker();
    } else if (initialLatitude && initialLongitude) {
      addSelectMarker(initialLongitude, initialLatitude);
    }
  }, [initialLatitude, initialLongitude, clearMarker]);

  const markersRef = useRef([]);

  useEffect(() => {
    if (!mapRef.current) return;

    // удалить старые пины
    markersRef.current.forEach(marker => marker.remove());
    markersRef.current = [];

    meetings.forEach(meeting => {
      if (!meeting.latitude || !meeting.longitude) {
        return;
      }

      const el = document.createElement("img");

      el.src =
        meeting.id === selectedMeetingId
          ? pinIconActive
          : pinIconInactive;
      el.style.width = meeting.id === selectedMeetingId ? "34px" : "24px";
      el.style.height = meeting.id === selectedMeetingId ? "34px" : "24px";

      el.addEventListener('click', (e) => {
        e.stopPropagation();
        setSelectedMeetingId(meeting.id);
      });

      const marker = new mapboxgl.Marker(el)
        .setLngLat([
          meeting.longitude,
          meeting.latitude,
        ])
        .addTo(mapRef.current);

      markersRef.current.push(marker);
    });

  }, [meetings, selectedMeetingId, setSelectedMeetingId]);

  return (
    <div
      ref={mapContainerRef}
      style={{
        width: '100%',
        height: '100%',
      }}
    />
  );
}

export default MapBox;