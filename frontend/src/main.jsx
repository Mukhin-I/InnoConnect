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
       <Route path="/chats/:id" element={<Chat />} />
        <Route path="/group-chats/:id" element={<GroupChats />} />
        <Route path="/meetings/:id/chat" element={<MeetingChatRedirect />} />
        <Route path="/welcome" element={<Welcome />} />
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>
)

