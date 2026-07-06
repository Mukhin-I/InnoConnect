import './Chat.css';
import logoIcon from './assets/logo.svg';
import notificationIcon from './assets/notifications.svg';
import settingsIcon from './assets/settings.svg';
import taskIcon from './assets/bin.svg';
import avatarImg from './assets/mock_ava.svg';
import contactIcon from './assets/calling.svg';
import attachIcon from './assets/paperclip.svg';
import sendIcon from './assets/send.svg';
import receivedIcon from './assets/received_mess.svg';
import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

function Chat() {
  const API_URL = import.meta.env.VITE_API_URL;
  const navigate = useNavigate();
  const { id: chatId } = useParams();

  const [contact, setContact] = useState({ name: 'Иван Петров', role: 'Житель Иннополиса' });
  const [messages, setMessages] = useState([]);
  const [myId, setMyId] = useState(null);
  const [draft, setDraft] = useState('');
  const [loading, setLoading] = useState(true);

  const task = {
    category: 'Помощь',
    title: 'Помочь с выносом мусора',
    time: 'Сегодня, 15:45',
  };

  const token = localStorage.getItem('token');

  const formatTime = (iso) =>
    new Date(iso).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });

  const toView = (list, meId) =>
    list.map((m) => ({
      id: m.id,
      type: m.sender.id === meId ? 'out' : 'in',
      text: m.text,
      time: formatTime(m.sent_at),
    }));

  useEffect(() => {
    const load = async () => {
      try {
        const meRes = await fetch(`${API_URL}/me`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        const me = meRes.ok ? await meRes.json() : null;
        const meId = me ? me.id : null;
        setMyId(meId);

        const infoRes = await fetch(`${API_URL}/chats/${chatId}`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (infoRes.ok) {
          const info = await infoRes.json();
          const other = info.participants.find((p) => p.id !== meId) || info.participants[0];
          if (other) setContact({ name: other.name, role: 'Житель Иннополиса' });
        }

        const msgRes = await fetch(`${API_URL}/chats/${chatId}/messages`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        if (msgRes.ok) {
          const data = await msgRes.json();
          setMessages(toView(data.messages, meId));
        } else {
          //useMockData();
          setMessages([]);
        }
      } catch (error) {
        //useMockData();
        setMessages([]);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [chatId]);

  const useMockData = () => {
    setMessages([
      { type: 'out', text: 'Здравствуйте! Могу помочь с выносом мусора', time: '15:45' },
      { type: 'in', text: 'Здравствуйте! Отлично, буду очень благодарен', time: '15:45' },
      { type: 'in', text: 'Буду дома в 20:00, подойдет?', time: '16:45' },
      { type: 'out', text: 'Ну если будешь дома значит сам и выкинешь', time: '16:47' },
      { type: 'in', text: 'Да пошел ты', time: '16:45' },
    ]);
  };

  const handleSend = async () => {
    const text = draft.trim();
    if (!text) return;
    setDraft('');
    try {
      const res = await fetch(`${API_URL}/chats/${chatId}/messages`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ text }),
      });
      if (res.ok) {
        const m = await res.json();
        setMessages((prev) => [...prev, ...toView([m], myId)]);
      }
    } catch (error) {
      console.error(error);
    }
  };

  if (loading) return <div className="rtr-loading">Загрузка...</div>;

  return (
    <div className="chat-page">
      <header className="chat-header">
        <div className="header-left">
        <button className="rtr-back-btn" onClick={() => navigate('/chats')}>
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M15 18L9 12L15 6" stroke="#1A1D1E" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
          </button>
          <div className="logo-container">
            <img src={logoIcon} alt="Logo" className="logo-icon" />
          </div>
        </div>
        <div className="header-icons">
          <img src={notificationIcon} alt="Notifications" className="header-icon" />
          <img src={settingsIcon} alt="Settings" className="header-icon" />
        </div>
      </header>

      <div className="rtr-card main-info-card rtr-card-1">
              <div className="rtr-icon-box">
               <img src={taskIcon} alt="Bin" className="rtr-iconbox-icon" />
              </div>
              <div className="rtr-main-details">
                <span className="rtr-category">{task.category}</span>
                <h1 className="rtr-title">{task.title}</h1>
                <span className="rtr-date">{task.time}</span>
              </div>
            </div>

 <div className="rtr-card author-card-full">
        <div className="rtr-author-left">
          <div className="rtr-avatar">
         <img src={avatarImg} alt="Avatar" className="rtr-iconbox-icon" />
          </div>
          <div className="rtr-author-info">
            <div className="rtr-author-name-row">
              <p className="rtr-name">{contact.name}</p>
              <p className="rtr-verified">Верифицирован</p>
            </div>
            <span className="rtr-author-sub">Житель Иннополиса</span>
          </div>
        </div>
        <button className="rtr-call-btn">
          <img src={contactIcon} alt="Call" className="rtr-iconbox-icon" />
        </button>
      </div>


      <div className="chat-messages">
  <div className="date-divider"><span>Сегодня</span></div>
  {messages.map((m, i) => {
    return (
      <div key={m.id ?? i} className={`msg-row ${m.type}`}>
        <div className={`bubble ${m.type}`}>
          <p className="bubble-text">{m.text}</p>
          <span className="bubble-meta">
            <span className="bubble-time">{m.time}</span>
            {m.type === 'out' && <img src={receivedIcon} alt="" className="bubble-check" />}
          </span>
        </div>
      </div>
    );
  })}
</div>

      <div className="chat-input">
        <span className="input-attach">
          <img src={attachIcon} alt="" className="input-attach-icon" />
        </span>
        <input
  type="text"
  className="input-field"
  placeholder="Напишите сообщение..."
  value={draft}
  onChange={(e) => setDraft(e.target.value)}
  onKeyDown={(e) => e.key === 'Enter' && handleSend()}
/>
<button className="input-send" onClick={handleSend}>
  <img src={sendIcon} alt="Отправить" className="input-send-icon" />
</button>
      </div>
    </div>
  );
}

export default Chat;
