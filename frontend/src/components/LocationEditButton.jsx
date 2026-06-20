import './LocationEditButton.css';

const LocationEditButton = ({ onClick, children = 'Изменить', className = '' }) => {
    return (
        <button
            type="button"
            onClick={onClick}
            className={`edit-button ${className}`}
        >
            {children}
        </button>
    );
};

export default LocationEditButton;