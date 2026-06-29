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
    const [filter, setFilter] = useState('all');

    const [chats, setChats] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchChats = async () => {
            setLoading(true);

            try {
                const response = await fetch("http://localhost:8080/chats");

                if (!response.ok) {
                    throw new Error("Ошибка загрузки");
                }

                const data = await response.json();
                setChats(data.chats);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchChats();


    }, []);

    const filteredChats = chats.filter(chat => {
        if (filter === 'all') {
            return true;
        }

        return chat.unread_count > 0;
    });

    return (
        <>
            <div className="chats-page">
                <div className="chats-page-content">
                    <div className="header-top">
                    <h2>InnoConnect</h2>
                    <IconButton 
                        icon={notificationIcon} 
                        alt=""
                    />
                    <IconButton 
                        icon={settingsIcon} 
                        alt=""
                    />
                    </div>

                    <h1 className="chats-header">Чаты</h1>

                    <div className="chats-top-bar">
                        <div className="chats-filter-container">
                            <div 
                                className={`chat-filter-item ${filter === 'all' ? 'chat-selected' : ''}`}
                                onClick={() => setFilter('all')}
                            >
                                <img src={allIcon} alt="all" />
                                <p>Все</p>
                            </div>

                            <div 
                                className={`chat-filter-item ${filter === 'unread' ? 'chat-selected chat-selected-unread' : ''}`}
                                onClick={() => setFilter('unread')}
                            >
                                <img src={unreadIcon} alt="unread" />
                                <p>Непрочитанные</p>
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