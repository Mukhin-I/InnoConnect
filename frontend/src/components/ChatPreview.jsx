import React from 'react';
import { useNavigate } from 'react-router-dom';
import './ChatPreview.css';
import chatAvatar from './chat-avatar.svg'

function ChatPreview( {chat} ) {
    const navigate = useNavigate();

    const formatTime = (dateString) => {
        const date = new Date(dateString);
        return date.toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
        });
    };

    const path = chat.type === 'MEETING'
        ? `/group-chats/${chat.chat_id}`
        : `/chats/${chat.chat_id}`;

    const title = chat.type === 'MEETING'
        ? (chat.request_title ?? 'Групповой чат')
        : (chat.participants?.[0]?.name ?? 'Чат');

    return(
        <>
            <div className="chat-card" onClick={() => navigate(path)}>
                <div className="chat-left">
                    <img src={chatAvatar} alt="avatar" className="avatar" />
                    <div className="chat-info">
                        <div className="chat-info-top">
                            <h2>{title}</h2>
                            {chat.request_title && (
                                <div className="request-associated"><p>{chat.request_title}</p></div>
                            )}
                        </div>

                        <p className="last-message">{chat.last_message?.text ?? 'Нет сообщений'}</p>
                    </div>
                </div>
                <div className="chat-right">
                    {chat.last_message && (
                        <p className="sent-at">{formatTime(chat.last_message.sent_at)}</p>
                    )}
                    {chat.unread_count > 0 && (
                        <div className="num-of-unread">
                            {chat.unread_count}
                        </div>
                    )}
                </div>
            </div>
        </>
    );
}

export default ChatPreview