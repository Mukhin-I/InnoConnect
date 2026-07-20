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
import BottomMenu from './components/BottomMenu.jsx';

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
    const [currentUser, setCurrentUser] = useState(null);
    const API_URL = import.meta.env.VITE_API_URL;

    useEffect(() => {
        const fetchRequests = async () => {
            setLoading(true);
            try {
                const response = await fetch(`${API_URL}/requests`);

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

        const fetchCurrentUser = async () => {
            const token = localStorage.getItem("token");

            if (!token) return;

            try {
                const response = await fetch(`${API_URL}/me`, {
                    headers: {
                        Authorization: `Bearer ${token}`,
                    },
                });

                if (!response.ok) {
                    throw new Error("Не удалось получить пользователя");
                }

                const data = await response.json();
                setCurrentUser(data);
            } catch (error) {
                console.error(error);
            }
        };

        fetchRequests();
        fetchCurrentUser();
    }, [API_URL]);


    if (loading) {
        return(
            <>
                <div className="requests-page">
                     <header className="map-header">
                               <div className="logo-container">
                             <img src={logoIcon} alt="Logo" style={{ width: 108, height: 25 }} />          </div>
                                       {/* <div className="header-icons">
                                         <img src={notificationIcon} alt="Notifications" className="header-icon" />
                                           <img src={settingsIcon} alt="Settings" className="header-icon" />
                                    </div> */}
                               </header>
                     
                <div className="requests-page-content">
    
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
            <BottomMenu initialSelected={'requests'} />
            </>
        );
    }

    if (error) {
        return(
            <>
                <div className="requests-page">
                     <header className="map-header">
                               <div className="logo-container">
                     <img src={logoIcon} alt="Logo" style={{ width: 108, height: 25 }} />          </div>
                               {/* <div className="header-icons">
                                 <img src={notificationIcon} alt="Notifications" className="header-icon" />
                                   <img src={settingsIcon} alt="Settings" className="header-icon" />
                               </div> */}
                             </header>
                     
                <div className="requests-page-content">
    
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
            <BottomMenu initialSelected={'requests'} />
            </>
        );
    }

    const filteredRequests = (requests || []).filter(request => {
        // Мои
        if (
            filter === "mine" &&
            request.creator_id !== currentUser?.id
        ) {
            return false;
        }

        // Категории
        if (
            requestType !== "allTypes" &&
            request.type !== requestType
        ) {
            return false;
        }

        return true;
    });

    return(
        <>
                <div className="requests-page">
                     <header className="map-header">
                                              <div className="logo-container">
                                                <img src={logoIcon} alt="Logo" className="logo-icon" />
                                              </div>
                                              {/* <div className="header-icons">
                                                <img src={notificationIcon} alt="Notifications" className="header-icon" />
                                                <img src={settingsIcon} alt="Settings" className="header-icon" />
                                              </div> */}
                                    </header>
                <div className="requests-page-content">
    
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
                        <div className={`type-of-request ${requestType === 'Помощь' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('Помощь')}
                        >
                            <div className="type-request-icon helpreq"></div>
                            <p>Помощь</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'Вещи' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('Вещи')}
                        >
                            <div className="type-request-icon stuffreq"></div>
                            <p>Вещи</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'Транспорт' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('Транспорт')}
                        >
                            <div className="type-request-icon transportreq"></div>
                            <p>Транспорт</p>
                        </div>
                        <div className={`type-of-request ${requestType === 'Прочее' ? 'type-selected' : ''}`}
                            onClick={() => setRequestType('Прочее')}
                        >
                            <div className="type-request-icon otherreq"></div>
                            <p>Прочее</p>
                        </div>
                    </div>


                    {/* карточки просьб */}
                    <div className="list-of-reqs-container">
                        {filteredRequests.map((request) => (
                            <CardOfRequest
                                key={request.request_id}
                                request={request}
                            />
                        ))}
                    </div>
                </div>
            </div>
            <BottomMenu initialSelected={'requests'} />
            </>
    );
}

export default Requests