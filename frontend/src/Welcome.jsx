import CategoryWelcome from './components/CategoryWelcome'
import "./Welcome.css"
import liveMapIcon from './assets/live-map.png'
import helpIcon from './assets/help.png'
import peopleIcon from './assets/people.png'
import welcomeImage from './assets/welcome-car.svg'
import { useNavigate } from 'react-router-dom';

function Welcome() {
  const navigate = useNavigate();

  const handleRegisterClick = () => {
    navigate('/register');
  };

  const handleLoginClick = () => {
    navigate('/login');
  };

  return (
    <>
     <div className="welcome-page">
      <div className="welcome-content">
        <div className="header-top">
          <h2>InnoConnect</h2>
        </div>
        <div className="title-container">
          <h1>Твой город -<br />твои люди</h1>
          <p>Удобный способ находить события и помогать соседям в Иннополисе в один клик.</p>
        </div>
        <img src={welcomeImage} alt="" className="welcome-image" />
        <section id="categories">
          <div className="categories-row">
            <CategoryWelcome
              icon={liveMapIcon}
              label="Живая карта"
            />
            <CategoryWelcome
              icon={helpIcon}
              label="Взаимопощь"
            />
            <CategoryWelcome
              icon={peopleIcon}
              label="Свои люди"
            />
          </div>
        </section>
        <div className="buttons-container">
          <button type="button" className="register-button"
          onClick={handleRegisterClick}>Регистрация</button>
          <button type="button" className="login-button"
          onClick={handleLoginClick}>Войти</button>
        </div>
      </div>
    </div>
    </>
  )
}

export default Welcome