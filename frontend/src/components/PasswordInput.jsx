import './PasswordInput.css';
import { useState } from 'react';
import eyeOpenIcon from '../assets/eye-open.png';
import eyeClosedIcon from '../assets/eye-closed.png';

const PasswordInput = ({
    value,
    onChange,
    placeholder = 'Пароль',
    className = '',
    icon = null,
    required = false,
    ...props
}) => {
    const [showPassword, setShowPassword] = useState(false);

    const iconStyle = icon ? {
        backgroundImage: `url(${icon})`,
        backgroundRepeat: 'no-repeat',
        backgroundPosition: 'left 1.33vw center',
        backgroundSize: '3.5vw',
        paddingLeft: '6.4vw',
    } : {};

    return (
        <div className="password-input-container">
            <input
                type={showPassword ? 'text' : 'password'}
                value={value}
                onChange={onChange}
                placeholder={placeholder}
                className={`input-icon ${className}`}
                style={iconStyle}
                required={required}
                {...props}
            />
            <button
                type="button"
                className="eye-button"
                onClick={() => setShowPassword(!showPassword)}
            >
                <img 
                    src={showPassword ? eyeOpenIcon : eyeClosedIcon} 
                    alt=''
                />
            </button>
        </div>
    );
};

export default PasswordInput;