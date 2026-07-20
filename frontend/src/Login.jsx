import './Login.css'
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import signUpBackground from './assets/signup-bg.svg'
import CreationInput from './components/CreationInput'
import avatarIcon from './assets/avatar.png'
import lockIcon from './assets/lock.png'
import CreateButton from './components/CreateButton'
import PasswordInput from './components/PasswordInput';
import logoIcon from './assets/logo.svg';

function Login() {
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
        const response = await fetch(`${API_URL}/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: email,
                password: password
            }),
        });

        const data = await response.json();

        if (response.status === 200) {
            localStorage.setItem('token', data.token);
            const expiresInMs = Date.now() + (data.expiresIn || 3600) * 1000;
            localStorage.setItem('expiresIn', expiresInMs.toString());

            setSuccess('Успешный вход!');
            setEmail('');
            setPassword('');
            navigate('/');
        } else if (response.status === 401) {
            setError(data.message);
        }
    } catch (error) {
        console.log(error);
        setEmail('');
        setPassword('');
        setError('Проблемы с подключением');
    } finally {
        setIsLoading(false);
    }
    };

  const handleRegisterClick = () => {
    navigate('/register');
  };

  return (
    <>
      <div className="login-page">
        <div className="login-page-content">
            <div className="header-top header-login">
                    <img src={logoIcon} alt="Logo" style={{ width: 65, height: 15 }} />
            </div>
            <div className="text-header">
                <h1>С возвращением<br />в InnoConnect</h1>
                <p>Рады видеть вас снова! Войдите, чтобы продолжить помогать соседям.</p>
            </div>
            <div className="login-card">
                <form onSubmit={handleSubmit}>
                    <div className="login-inputs">
                        <CreationInput
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        placeholder="Email"
                        maxLength={50}
                        icon={avatarIcon}
                        required
                        disabled={isLoading}
                        fontSize="3.2vw"
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
                    <button className="forget-password">Забыли пароль?</button>
                   {error && <div className="error-message">{error}</div>}
                   {success && <div className="success-message">{success}</div>}
                   <CreateButton type="submit" disabled={isLoading}>
                    {isLoading ? 'Авторизация..' : 'Войти'}
                   </CreateButton>
                </form>
                <div className="no-account">
                    <p>Нет аккаунта?</p>
                    <button type="button" className="no-account-button"
                    onClick={handleRegisterClick}>Зарегистрироваться</button>
                </div>
                <div className="img-container">
                    <img src={signUpBackground} alt=""></img>
                </div>
            </div>
        </div>
      </div>
    </>
  )
}

export default Login