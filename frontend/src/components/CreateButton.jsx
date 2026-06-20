import './CreateButton.css';

const CreateButton = ({ 
    children,
    onClick, 
    className = '', 
    type = 'button',
    disabled = false,
    ...props 
}) => {
    return (
        <button
            type={type}
            onClick={onClick}
            disabled={disabled}
            className={`create-button ${className}`}
            {...props}
        >
            {children}
        </button>
    );
};

export default CreateButton;