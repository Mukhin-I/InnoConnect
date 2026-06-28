import './BottomMenu.css';
import { Link } from 'react-router-dom';
import { useState } from 'react';
import toMap from './map-nav.svg'
import toReqs from './requests-nav.svg'
import toChats from './chat-nav.svg'
import toProfile from './profile-nav.svg'

function BottomMenu( {initialSelected} ) {

    const [selected, setSelected] = useState(initialSelected)

    return (
        <>
            <div className="bottom-menu">
                <div className="menu-links">
                    <Link to="/" onClick={() => setSelected('map')}>
                        <img src={toMap} alt="map" className={`${selected === 'map' ? "page-selected" : ""}`}/>
                        {selected === 'map' && <div className="selected-underline"></div>}
                    </Link>
                    <Link to="/requests" onClick={() => setSelected('requests')}>
                        <img src={toReqs} alt="requests" className={`${selected === 'requests' ? "page-selected" : ""}`}/>
                        {selected === 'requests' && <div className="selected-underline"></div>}
                    </Link>
                    <Link to="/chats" onClick={() => setSelected('chats')}>
                        <img src={toChats} alt="chats" className={`${selected === 'chats' ? "page-selected" : ""}`}/>
                        {selected === 'chats' && <div className="selected-underline"></div>}
                    </Link>
                    <Link to="/profile" onClick={() => setSelected('profile')}>
                        <img src={toProfile} alt="profile" className={`${selected === 'profile' ? "page-selected" : ""}`}/>
                        {selected === 'profile' && <div className="selected-underline"></div>}
                    </Link>
                </div>
            </div>
        </>
    );
}

export default BottomMenu;