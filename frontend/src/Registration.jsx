import './Registration.css'
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import signUpBackground from './assets/signup-bg.svg'
import CreationInput from './components/CreationInput'
import avatarIcon from './assets/avatar.png'
import lockIcon from './assets/lock.png'
import CreateButton from './components/CreateButton'
import mailIcon from './assets/mailIcon.png'
import languageIcon from './assets/languageIcon.png'
import PasswordInput from './components/PasswordInput';

function Registration() {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const navigate =useNavigate();
  const API_URL = import.meta.env.VITE_API_URL;

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setIsLoading(true);

    try {
        const response = await fetch(`${API_URL}/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: name,
                email: email,
                password: password
            }),
        });

        const data = await response.json();

        if (response.status === 201) {
            setSuccess(data.message);
            setName('');
            setEmail('');
            setPassword('');
            navigate('/login');
        } else if (response.status === 400) {
            setError(data.message);
        }
    } catch (error) {
        console.log(error);
        setName('');
        setEmail('');
        setPassword('');
        setError('Проблемы с подключением');
    } finally {
        setIsLoading(false);
    }
    };

  const handleLoginClick = () => {
    navigate('/login');
  };

  return (
    <>
      <div className="registration-page">
        <div className="registration-page-content">
            <div className="header-top">
                <h2>InnoConnect</h2>
                <button className="language-change-button">
                  <img src={languageIcon} alt=""/>
                  Русский
                </button>
            </div>
            <div className="text-header header-register">
                <h1>Добро пожаловать<br />в InnoConnect</h1>
                <p>Создайте аккаунт, чтобы находить события, помогать соседям и делать Иннополис лучше каждый день.</p>
            </div>
            <div className="img-container">
              <img src={signUpBackground} alt=""></img>
            </div>
            <div className="register-card">
                <div className="card-header">
                    <h2>Создать аккаунт</h2>
                </div>
                <form onSubmit={handleSubmit}>
                    <div className="register-inputs">
                        <CreationInput
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        placeholder="Ваше имя"
                        maxLength={50}
                        icon={avatarIcon}
                        required
                        disabled={isLoading}
                        />
                        <CreationInput
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="Email"
                        maxLength={50}
                        icon={mailIcon}
                        type="email"
                        required
                        disabled={isLoading}
                        />
                        <PasswordInput
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            placeholder="Пароль"
                            maxLength={50}
                            icon={lockIcon}
                            required
                            disabled={isLoading}
                        />
                    </div>
                   {error && <div className="error-message">{error}</div>}
                   {success && <div className="success-message">{success}</div>}
                   <CreateButton type="submit" disabled={isLoading}>
                    {isLoading ? 'Регистрация...' : 'Зарегистрироваться'}
                   </CreateButton>
                </form>
                <div className="already-account">
                    <p>Уже есть аккаунт?</p>
                    <button type="button" className="already-login-button"
                    onClick={handleLoginClick}>Войти</button>
                </div>
            </div>
        </div>
      </div>
    </>
  )
}

export default Registration