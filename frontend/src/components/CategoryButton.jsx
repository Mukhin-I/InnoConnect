import './CategoryButton.css';

const CategoryButton = ({
    icon,
    label,
    isActive = false,
    onClick,
    className = '',
    ...props
}) => {
    return (
        <button
            type="button"
            className={`category-button ${isActive ? 'active' : ''} ${className}`}
            onClick={onClick}
            {...props}
        >
            <img 
                src={icon} 
                alt={label} 
                className="category-icon"
            />
            <span className="category-label">{label}</span>
        </button>
    );
};

export default CategoryButton;