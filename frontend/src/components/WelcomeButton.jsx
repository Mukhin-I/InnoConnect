import './WelcomeButton.css';

const WelcomeButton = ({ onClick, children }) => {
  return (
    <button 
      type="button" 
      className="welcome-button"
      onClick={onClick}
    >
      {children || 'Регистрация'}
    </button>
  );
};

export default WelcomeButton;