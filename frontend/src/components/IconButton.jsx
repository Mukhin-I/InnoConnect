import './IconButton.css';

const IconButton = ({ icon, alt, className = '' }) => {
    return (
        <button 
            className={`icon-button ${className}`} 
        >
            <img src={icon} alt={alt || ''} />
        </button>
    );
};

export default IconButton;