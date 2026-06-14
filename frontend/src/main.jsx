import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import EventCard from './EventCard.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <EventCard />
  </StrictMode>,
)
