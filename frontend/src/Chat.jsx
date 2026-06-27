import './Chat.css';
import { useNavigate } from 'react-router-dom';
import logoIcon from './assets/logo.svg';
import notificationIcon from './assets/notifications.svg';
import settingsIcon from './assets/settings.svg';
import taskIcon from './assets/bin.svg';
import avatarImg from './assets/mock_ava.svg';
import contactIcon from './assets/calling.svg';
import attachIcon from './assets/paperclip.svg';
import sendIcon from './assets/send.svg';
import receivedIcon from './assets/received_mess.svg';

function Chat() {


  const contact = {
    name: 'Иван Петров',
    role: 'Житель Иннополиса',
  };
const navigate = useNavigate();
  const task = {
    category: 'Помощь',
    title: 'Помочь с выносом мусора',
    time: 'Сегодня, 15:45',
  };

  const messages = [
    { type: 'date', label: 'Сегодня' },
    { type: 'out', text: 'Здравствуйте! Могу помочь с выносом мусора', time: '15:45' },
    { type: 'in', text: 'Здравствуйте! Отлично, буду очень благодарен', time: '15:45' },
    { type: 'new', label: 'Новые сообщения' },
    { type: 'in', text: 'Буду дома в 20:00, подойдет?', time: '16:45' },
    { type: 'out', text: 'Ну если будешь дома значит сам и выкинешь', time: '16:47' },
    { type: 'in', text: 'Да пошел ты', time: '16:45' },
  ];

  return (
    <div className="chat-page">
      <header className="chat-header">
        <div className="header-left">
        <button className="rtr-back-btn" onClick={() => navigate('/chats')}>
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
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
        {messages.map((m, i) => {
          if (m.type === 'date') {
            return <div key={i} className="date-divider"><span>{m.label}</span></div>;
          }
          if (m.type === 'new') {
            return <div key={i} className="new-divider"><span>{m.label}</span></div>;
          }
          return (
            <div key={i} className={`msg-row ${m.type}`}>
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
        <input type="text" className="input-field" placeholder="Напишите сообщение..." />
        <button className="input-send">
          <img src={sendIcon} alt="Отправить" className="input-send-icon" />
        </button>
      </div>
    </div>
  );
}

export default Chat;
