import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import './Requests.css';
import logoIcon from './assets/logo.svg';
import notificationIcon from './assets/notifications.svg';
import settingsIcon from './assets/settings.svg';
import IconButton from './components/IconButton';
import allIcon from './assets/all.svg';
import mineIcon from './assets/mine.svg';
import allRequestsIcon from './assets/allrequests.svg'
import CardOfRequest from './components/CardOfRequest.jsx'

function Requests() {
    const navigate = useNavigate();
    const handleClose = () => {
        navigate('/');
    };

    const [filter, setFilter] = useState('all');
    const [requestType, setRequestType] = useState('allTypes');

    const [requests, setRequests] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchRequests = async () => {
            setLoading(true);
            try {
                const response = await fetch("http://localhost:8080/requests");

                if (!response.ok) {
                    throw new Error("Ошибка загрузки");
                }

                const data = await response.json();
                setRequests(data.requests);
            } catch (error) {
                setError(error.message);
            } finally {
                setLoading(false);
            }
        };

        fetchRequests();
    }, []);

    if (loading) {
        return(
            <>
                <div className="requests-page">
                <div className="requests-page-content">
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

                    <h1 className="requests-header">Просьбы</h1>

                    <div className="requests-top-bar">
                        <div className="requests-filter-container">
                            <div 
                                className={`request-filter-item ${filter === 'all' ? 'request-selected' : ''}`}
                                onClick={() => setFilter('all')}
                            >
                                <img src={allIcon} alt="all" />
                                <p>Все</p>
                            </div>

                            <div 
                                className={`request-filter-item ${filter === 'mine' ? 'request-selected' : ''}`}
                                onClick={() => setFilter('mine')}
                            >
                                <img src={mineIcon} alt="mine" />
                                <p>Мои</p>
                            </div>
                        </div>

                        <Link to="/create-request" className="add-request">
                            <span>+</span>
                        </Link>
                    </div>

                    <div className="type-of-requests-filter">
                        <div className={`type-of-request ${requestType === 'allTypes' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('allTypes')}
                        >
                            <div className="type-request-icon allrequests"></div>
                            <p>Все категории</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'help' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('help')}
                        >
                            <div className="type-request-icon helpreq"></div>
                            <p>Помощь</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'stuff' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('stuff')}
                        >
                            <div className="type-request-icon stuffreq"></div>
                            <p>Вещи</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'transport' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('transport')}
                        >
                            <div className="type-request-icon transportreq"></div>
                            <p>Транспорт</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'other' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('other')}
                        >
                            <div className="type-request-icon otherreq"></div>
                            <p>Прочее</p>
                        </div>
                    </div>


                    {/* карточки просьб */}
                    <div className="list-of-reqs-container">
                        <p>Загрузка...</p>
                    </div>
                </div>
            </div>
            </>
        );
    }

    if (error) {
        return(
            <>
                <div className="requests-page">
                <div className="requests-page-content">
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

                    <h1 className="requests-header">Просьбы</h1>

                    <div className="requests-top-bar">
                        <div className="requests-filter-container">
                            <div 
                                className={`request-filter-item ${filter === 'all' ? 'request-selected' : ''}`}
                                onClick={() => setFilter('all')}
                            >
                                <img src={allIcon} alt="all" />
                                <p>Все</p>
                            </div>

                            <div 
                                className={`request-filter-item ${filter === 'mine' ? 'request-selected' : ''}`}
                                onClick={() => setFilter('mine')}
                            >
                                <img src={mineIcon} alt="mine" />
                                <p>Мои</p>
                            </div>
                        </div>

                        <Link to="/create-request" className="add-request">
                            <span>+</span>
                        </Link>
                    </div>

                    <div className="type-of-requests-filter">
                        <div className={`type-of-request ${requestType === 'allTypes' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('allTypes')}
                        >
                            <div className="type-request-icon allrequests"></div>
                            <p>Все категории</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'help' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('help')}
                        >
                            <div className="type-request-icon helpreq"></div>
                            <p>Помощь</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'stuff' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('stuff')}
                        >
                            <div className="type-request-icon stuffreq"></div>
                            <p>Вещи</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'transport' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('transport')}
                        >
                            <div className="type-request-icon transportreq"></div>
                            <p>Транспорт</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'other' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('other')}
                        >
                            <div className="type-request-icon otherreq"></div>
                            <p>Прочее</p>
                        </div>
                    </div>


                    {/* карточки просьб */}
                    <div className="list-of-reqs-container">
                        <p>Ошибка: {error}</p>
                    </div>
                </div>
            </div>
            </>
        );
    }

    const filteredRequests = requests.filter(request => {
        if (requestType === 'allTypes') {
            return true;
        }

        return request.type.toLowerCase() === requestType;
    });

    return(
        <>
            <div className="requests-page">
                <div className="requests-page-content">
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

                    <h1 className="requests-header">Просьбы</h1>

                    <div className="requests-top-bar">
                        <div className="requests-filter-container">
                            <div 
                                className={`request-filter-item ${filter === 'all' ? 'request-selected' : ''}`}
                                onClick={() => setFilter('all')}
                            >
                                <img src={allIcon} alt="all" />
                                <p>Все</p>
                            </div>

                            <div 
                                className={`request-filter-item ${filter === 'mine' ? 'request-selected' : ''}`}
                                onClick={() => setFilter('mine')}
                            >
                                <img src={mineIcon} alt="mine" />
                                <p>Мои</p>
                            </div>
                        </div>

                        <Link to="/create-request" className="add-request">
                            <span>+</span>
                        </Link>
                    </div>

                    <div className="type-of-requests-filter">
                        <div className={`type-of-request ${requestType === 'allTypes' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('allTypes')}
                        >
                            <div className="type-request-icon allrequests"></div>
                            <p>Все категории</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'help' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('help')}
                        >
                            <div className="type-request-icon helpreq"></div>
                            <p>Помощь</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'things' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('things')}
                        >
                            <div className="type-request-icon stuffreq"></div>
                            <p>Вещи</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'transport' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('transport')}
                        >
                            <div className="type-request-icon transportreq"></div>
                            <p>Транспорт</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'other' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('other')}
                        >
                            <div className="type-request-icon otherreq"></div>
                            <p>Прочее</p>
                        </div>
                    </div>


                    {/* карточки просьб */}
                    <div className="list-of-reqs-container">
                        {filteredRequests.map(request => (
                            <CardOfRequest
                                key={request.request_id}
                                request={request}
                            />
                        ))}
                    </div>
                </div>
            </div>
        </>
    );
}

export default Requests