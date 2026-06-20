import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import './CardOfRequest.css';

function CardOfRequest() {
    // Date in such format: day.month, hour:minutes like 02.02, 15:05
    const formatDate = (dateString) => {
        const date = new Date(dateString);  
        const day = date.getDate().toString().padStart(2, '0');
        const month = (date.getMonth() + 1).toString().padStart(2, '0');
        const hours = date.getHours().toString().padStart(2, '0');
        const minutes = date.getMinutes().toString().padStart(2, '0');
        return `${day}.${month}, ${hours}:${minutes}`;
    };

    const [request, setRequest] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    return(
        <>
            <div className="request-card-container">
                <div className="request-card-icon">
                    <div className="type-request-icon-card helpreq"></div>
                </div>

                <div className="request-card-info">
                    <p className="request-card-type">Помощь</p>
                    <p className="request-card-title">Помочь с выносом мусора</p>
                    <p className="request-card-date">Сегодня, 15:45</p>
                </div>

                <div className="new-request-tag">
                    <p>Новое</p>
                </div>
            </div>
        </>
    );
}

export default CardOfRequest