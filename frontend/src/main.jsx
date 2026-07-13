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
import GroupChats from './GroupChats.jsx'
import MeetingChatRedirect from './MeetingChatRedirect.jsx'
import ProtectedRoute from './components/ProtectedRoute.jsx'
import EventCard from './EventCard.jsx'

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Map />} />
        <Route path="/create-meeting" element={<ProtectedRoute><MeetingCreation /></ProtectedRoute>} />
        <Route path="/request/:id" element={<RespondToRequest />} />
        <Route path="/create-request" element={<ProtectedRoute><RequestCreation /></ProtectedRoute>} />
        <Route path="/requests" element={<Requests />} />
        <Route path="/chats" element={<ProtectedRoute><Chats /></ProtectedRoute>} />
        <Route path="/profile" element={<ProtectedRoute><Profile /></ProtectedRoute>} />
        <Route path="/chats/:id" element={<ProtectedRoute><Chat /></ProtectedRoute>} />
        <Route path="/group-chats/:id" element={<ProtectedRoute><GroupChats /></ProtectedRoute>} />
        <Route path="/meetings/:id/chat" element={<ProtectedRoute><MeetingChatRedirect /></ProtectedRoute>} />
        <Route path="/welcome" element={<Welcome />} />
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>
)


