import { Link } from 'react-router-dom';
import React from 'react';
import './CardOfRequest.css';

function CardOfRequest({ request }) {
    // Date in such format: day.month, hour:minutes like 02.02, 15:05
    const formatDate = (dateString) => {
        const date = new Date(dateString);  
        const day = date.getDate().toString().padStart(2, '0');
        const month = (date.getMonth() + 1).toString().padStart(2, '0');
        const hours = date.getHours().toString().padStart(2, '0');
        const minutes = date.getMinutes().toString().padStart(2, '0');
        return `${day}.${month}, ${hours}:${minutes}`;
    };

    const typeOfReq = {
        "Помощь": "helpreq",
        "Вещи": "stuffreq",
        "Транспорт": "transportreq",
        "Прочее": "otherreq",
    }

    return (
        <Link 
            to={`/request/${request.request_id}`}
            className="request-card-container"
        >
            <div className="request-card-icon">
                <div className={`type-request-icon-card ${typeOfReq[request.type]}`}></div>
            </div>

            <div className="request-card-info">
                <p className="request-card-type">{request.type}</p>
                <p className="request-card-title">{request.title}</p>
                <p className="request-card-date">
                    {formatDate(request.deadline)}
                </p>
            </div>

            <div className="new-request-tag">
                <p>Новое</p>
            </div>
        </Link>
    );
}

export default CardOfRequest