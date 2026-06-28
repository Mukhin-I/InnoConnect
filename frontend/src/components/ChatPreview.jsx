import React from 'react';
import './ChatPreview.css';
import chatAvatar from './chat-avatar.svg'

function ChatPreview() {
    return(
        <>
            <div className="chat-card">
                <div className="chat-left">
                    <img src={chatAvatar} alt="avatar" className="avatar" />
                    <div className="chat-info">
                        <div className="chat-info-top">
                            <h2>Иван Петров</h2>
                            <div className="request-associated"><p>Одолжите дрель</p></div>
                        </div>
            
                        <p className="last-message">Спасибо большое!</p>
                    </div>
                </div>
                <div className="chat-right">
                    <p className="sent-at">15:45</p>
                    <div className="num-of-unread">2</div>
                </div>
            </div>
        </>
    );
}

export default ChatPreview