import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import './Chats.css';
import logoIcon from './assets/logo.svg';
import notificationIcon from './assets/notifications.svg';
import settingsIcon from './assets/settings.svg';
import BottomMenu from './components/BottomMenu.jsx';
import IconButton from './components/IconButton';
import allIcon from './assets/all.svg';
import unreadIcon from './assets/unread.svg';
import ChatPreview from './components/ChatPreview.jsx';

function Chats() {
    const [filter, setFilter] = useState('meeting');

    const [chats, setChats] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const API_URL = import.meta.env.VITE_API_URL;

    useEffect(() => {
        const fetchChats = async () => {
            setLoading(true);

            try {
                const token = localStorage.getItem('token');

                const response = await fetch(`${API_URL}/chats`, {
                headers: { Authorization: `Bearer ${token}` },
                });

                if (!response.ok) {
                    throw new Error("Ошибка загрузки");
                }

                const data = await response.json();
                setChats(data.chats ?? []);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchChats();


    }, []);

    const filteredChats = chats.filter(chat => {
        if (filter === 'meeting') {
            return chat.type === 'MEETING';
        }

        return chat.type === 'REQUEST';
    });

    return (
        <>
            <div className="chats-page">
                <header className="profile-header">
                          <div className="logo-container">
                            <img src={logoIcon} alt="Logo" className="logo-icon" />
                          </div>
                          <div className="header-icons">
                            <img src={notificationIcon} alt="Notifications" className="header-icon" />
                            <img src={settingsIcon} alt="Settings" className="header-icon" />
                          </div>
                </header>
                <div className="chats-page-content">

                    <h1 className="chats-header">Чаты</h1>

                    <div className="chats-top-bar">
                        <div className="chats-filter-container">
                            <div 
                                className={`chat-filter-item ${filter === 'meeting' ? 'chat-selected' : ''}`}
                                onClick={() => setFilter('meeting')}
                            >
                                {/* <img src={allIcon} alt="all" /> */}
                                <p>Мероприятия</p>
                            </div>

                            <div 
                                className={`chat-filter-item ${filter === 'request' ? 'chat-selected' : ''}`}
                                onClick={() => setFilter('request')}
                            >
                                {/* <img src={unreadIcon} alt="unread" /> */}
                                <p>Просьбы</p>
                            </div>
                        </div>

                        {/* <Link to="/create-request" className="add-request">
                            <span>+</span>
                        </Link> */}
                    </div>

                    <div className="list-of-chats">
                        {filteredChats.map(chat => (
                            <ChatPreview
                                key={chat.chat_id}
                                chat={chat}
                            />
                        ))}
                    </div>
                </div>
            </div>
            <BottomMenu initialSelected={'chats'} />
        </>
    );
}

export default Chats