import React from 'react';
import './CardOfMeeting.css';
import calendarIcon from '../assets/calendar.png';
import locationIcon from '../assets/location.png';
import sportIcon from '../assets/football.svg';
import socialIcon from '../assets/coffee-cup.svg';
import studyIcon from '../assets/books.svg';
import gantelIcon from '../assets/guntell.svg';

const CATEGORY_TAG_ICON = {
    'Спорт': sportIcon,
    'Соц': socialIcon,
    'Учеба': studyIcon,
};

const CATEGORY_BIG_ICON = {
    'Спорт': gantelIcon,
    'Соц': socialIcon,
    'Учеба': studyIcon,
};

function CardOfMeeting({ meeting, onClick }) {
    const formatDate = (dateString) => {
        if (!dateString) return '';
        const start = new Date(dateString);
        const now = new Date();
        const tomorrow = new Date(now);
        tomorrow.setDate(now.getDate() + 1);

        const hours = start.getHours().toString().padStart(2, '0');
        const minutes = start.getMinutes().toString().padStart(2, '0');

        let prefix;
        if (start.toDateString() === now.toDateString()) {
            prefix = 'Сегодня';
        } else if (start.toDateString() === tomorrow.toDateString()) {
            prefix = 'Завтра';
        } else {
            const day = start.getDate().toString().padStart(2, '0');
            const month = (start.getMonth() + 1).toString().padStart(2, '0');
            prefix = `${day}.${month}`;
        }

        return `${prefix}, ${hours}:${minutes}`;
    };

    return (
        <div
            className="meeting-card-container"
            onClick={() => onClick && onClick(meeting.id)}
        >
            <div className="meeting-card-icon">
                <img src={CATEGORY_BIG_ICON[meeting.type] || gantelIcon} alt="" />
            </div>

            <div className="meeting-card-info">
                <p className="meeting-card-title">{meeting.title}</p>
                <div className="meeting-card-row">
                    <img src={calendarIcon} alt="" />
                    <span className="meeting-card-row-calendar">{formatDate(meeting.meeting_time)}</span>
                </div>
                <div className="meeting-card-row row-2">
                    <img src={locationIcon} alt="" />
                    <span className="meeting-card-row-location">{meeting.address || 'Адрес не указан'}</span>
                </div>
            </div>
        </div>
    );
}

export default CardOfMeeting;
