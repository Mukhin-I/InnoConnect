import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import MeetingCreation from './MeetingCreation.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <MeetingCreation />
  </StrictMode>,
)
