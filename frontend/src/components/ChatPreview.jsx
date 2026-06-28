import React from 'react';
import './ChatPreview.css';
import chatAvatar from './chat-avatar.svg'

function ChatPreview( {chat} ) {
    const formatTime = (dateString) => {
        const date = new Date(dateString);

        return date.toLocaleTimeString([], {
            hour: "2-digit",
            minute: "2-digit",
        });
    };
    
    return(
        <>
            <div className="chat-card">
                <div className="chat-left">
                    <img src={chatAvatar} alt="avatar" className="avatar" />
                    <div className="chat-info">
                        <div className="chat-info-top">
                            <h2>{chat.participants[0].name}</h2>
                            <div className="request-associated"><p>{chat.request_title}</p></div>
                        </div>
            
                        <p className="last-message">{chat.last_message.text}</p>
                    </div>
                </div>
                <div className="chat-right">
                    <p className="sent-at">{formatTime(chat.last_message.sent_at)}</p>
                    {
                        chat.unread_count > 0 && (
                            <div className="num-of-unread">
                                {chat.unread_count}
                            </div>
                        )
                    }
                </div>
            </div>
        </>
    );
}

export default ChatPreview