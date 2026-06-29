import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import './index.css'
import Map from './Map.jsx'
import MeetingCreation from './MeetingCreation.jsx'
import RespondToRequest from './RespondToRequest.jsx'
import RequestCreation from './RequestCreation.jsx'
import Requests from './Requests.jsx'
import Chat from './Chat.jsx'
import Chats from './Chats.jsx'
import Profile from './Profile.jsx'
import Welcome from './Welcome.jsx'
import Register from './Registration.jsx'
import Login from './Login.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Map />} />
        <Route path="/create-meeting" element={<MeetingCreation />} />
        <Route path="/request/:id" element={<RespondToRequest />} />
        <Route path="/create-request" element={<RequestCreation />} />
        <Route path="/requests" element={<Requests />} />
        <Route path="/chats" element={<Chats />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/chat/:id" element={<Chat />} />
        <Route path="/welcome" element={<Welcome />} />
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>
)
