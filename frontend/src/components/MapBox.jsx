import './MapBox.css'
import { useEffect, useRef } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';

function MapBox({
  center = [48.752, 55.752],
  zoom = 12,
  style = 'mapbox://styles/seanshushickkk/cmqazkvqy003d01qz36xq2coo',
}) {
  const mapRef = useRef(null);

  useEffect(() => {
    mapboxgl.accessToken = import.meta.env.VITE_MAPBOX_TOKEN;
    console.log(import.meta.env.VITE_MAPBOX_TOKEN)

    const map = new mapboxgl.Map({
      container: mapRef.current,
      style,
      center,
      zoom,
      attributionControl: false,
    });

    return () => map.remove();
  }, [center, zoom, style]);

  return (
    <div
      ref={mapRef}
      style={{
        width: '100%',
        height: '100%',
      }}
    />
  );
}

export default MapBox;
