import './GroupChats.css';
import logoIcon from './assets/logo.svg';
import notificationIcon from './assets/notifications.svg';
import settingsIcon from './assets/settings.svg';
import taskIcon from './assets/guntell.svg';
import avatarImg from './assets/mock_ava.svg';
import attachIcon from './assets/paperclip.svg';
import sendIcon from './assets/send.svg';
import receivedIcon from './assets/received_mess.svg';
import placeIcon from './assets/location.svg';
import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

function GroupChat() {
  const navigate = useNavigate();
  const { id: meetingId } = useParams();
  const [chatId, setChatId] = useState(null);
  const [participants, setParticipants] = useState([]);
  const [messages, setMessages] = useState([]);
  const [myId, setMyId] = useState(null);
  const [draft, setDraft] = useState('');
  const [loading, setLoading] = useState(true);

  const task = {
    category: 'Спорт',
    title: 'Утренняя пробежка',
    time: 'Сегодня, 7:30 - 9:00',
    place: 'Спорткомплекс',
  };

  const token = localStorage.getItem('token');

  const formatTime = (iso) =>
    new Date(iso).toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });

  const toView = (list, meId) =>
    list.map((m) => ({
      id: m.id,
      type: m.sender.id === meId ? 'out' : 'in',
      name: m.sender.name,
      text: m.text,
      time: formatTime(m.sent_at),
    }));

  useEffect(() => {
  const load = async () => {
    try {
      // 1. встреча -> chat_id
      const chatRes = await fetch(`http://localhost:8080/meetings/${meetingId}/chat`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!chatRes.ok) { useMockData(); return; }
      const { chat_id } = await chatRes.json();
      setChatId(chat_id);

      // 2. кто я (нужно для in/out)
      const meRes = await fetch('http://localhost:8080/me', {
        headers: { Authorization: `Bearer ${token}` },
      });
      const me = meRes.ok ? await meRes.json() : null;
      const meId = me ? me.id : null;
      setMyId(meId);

      // 3. инфа о чате (участники)
      const infoRes = await fetch(`http://localhost:8080/chats/${chat_id}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (infoRes.ok) {
        const info = await infoRes.json();
        setParticipants(
          info.participants.map((p) => ({
            id: p.id,
            name: p.id === meId ? 'Вы' : p.name,
          }))
        );
      }

      // 4. сообщения
      const msgRes = await fetch(`http://localhost:8080/chats/${chat_id}/messages`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (msgRes.ok) {
        const data = await msgRes.json();
        setMessages(toView(data.messages, meId));
      } else {
        useMockData();
      }
    } catch (error) {
      useMockData();
    } finally {
      setLoading(false);
    }
  };
  load();
}, [meetingId]);

  const useMockData = () => {
    setParticipants([
      { id: 0, name: 'Вы' },
      { id: 1, name: 'Анна' },
      { id: 2, name: 'Павел' },
      { id: 3, name: 'Сергей' },
      { id: 4, name: 'Саша' },
      { id: 5, name: 'Дмитрий' },
      { id: 6, name: 'Марина' },
      { id: 7, name: 'Олег' },
      { id: 8, name: 'Катя' },
      { id: 9, name: 'Игорь' },
      { id: 10, name: 'Лена' },
      { id: 11, name: 'Артем' },
    ]);
    setMessages([
      { type: 'out', text: 'Всем привет!\nНапоминаю, всчтречаемся у входа в 7:25', time: '6:45' },
      { type: 'in', name: 'Анна', text: 'Привет! Я буду чуть раньше', time: '7:00'},
      { type: 'in', name: 'Дмитрий', text: 'Привет! Буду впервые.\nСильно высокий темп?', time: '7:01' },
      { type: 'out', text: 'Да не, чисто так по кайфу', time: '6:45' },
      { type: 'in', name: 'Марина', text: 'Возьму с собой воду. Кому еще надо - пишите', time: '7:01'},
    ]);
  };

  const handleSend = async () => {
    const text = draft.trim();
    if (!text) return;
    setDraft('');
    try {
      const res = await fetch(`http://localhost:8080/chats/${chatId}/messages`, {
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

  const shown = participants.slice(0, 5);
  const extra = participants.length - shown.length;

  return (
    <div className="chat-page group-chat">
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
          <span className="rtr-place">
            <img src={placeIcon} alt="" className="rtr-place-icon" />
            {task.place}
          </span>
        </div>
        <svg className="rtr-task-chevron" width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M15 18L9 12L15 6" stroke="#C0C1C4" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
        </svg>
      </div>

      <div className="rtr-card participants-card">
        <div className="participants-head">
          <span className="participants-title">Участники</span>
          <span className="participants-count">· {participants.length}</span>
        </div>
        <div className="participants-list">
          {shown.map((p) => (
            <div className="participant" key={p.id}>
              <div className="participant-avatar">
                <img src={avatarImg} alt={p.name} className="participant-img" />
                <span className="participant-dot"></span>
              </div>
              <span className="participant-name">{p.name}</span>
            </div>
          ))}
          {extra > 0 && (
            <div className="participant">
              <div className="participant-avatar participant-more">+{extra}</div>
              <span className="participant-name">Еще</span>
            </div>
          )}
        </div>
      </div>

      <div className="chat-messages">
        <div className="date-divider"><span>Сегодня</span></div>
        {messages.map((m, i) => (
          <div key={m.id ?? i} className={`msg-row ${m.type}`}>
            {m.type === 'in' && <img src={avatarImg} alt="" className="msg-avatar" />}
            <div className={`bubble ${m.type}`}>
              {m.type === 'in' && <span className="bubble-name">{m.name}</span>}
              <p className="bubble-text">{m.text}</p>
              <span className="bubble-meta">
                <span className="bubble-time">{m.time}</span>
                {m.type === 'out' && <img src={receivedIcon} alt="" className="bubble-check" />}
              </span>
            </div>
          </div>
        ))}
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

export default GroupChat;
