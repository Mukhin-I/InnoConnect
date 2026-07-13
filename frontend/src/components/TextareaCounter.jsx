import './TextareaCounter.css';
import descriptionIcon from '../assets/descriptionIcon.png';

const TextareaCounter = ({
    value,
    onChange,
    placeholder = '',
    maxLength = 500,
    rows = 5,
    className = '',
    icon = descriptionIcon,
    required = false,
    ...props
}) => {
    const iconStyle = {
        backgroundImage: `url(${icon})`,
        backgroundRepeat: 'no-repeat',
        backgroundPosition: 'left 1.33vw top 1.2vh',
        backgroundSize: '2.66vw',
        paddingLeft: '6.4vw',
        paddingTop: '1.5vh',
        paddingBottom: '3vh',
    };

    return (
        <div className={`textarea-wrapper ${className}`}>
            <textarea
                value={value}
                onChange={onChange}
                placeholder={placeholder}
                maxLength={maxLength}
                rows={rows}
                required={required}
                className="textarea-counter"
                style={icon ? iconStyle : {}}
                {...props}
            />
            <div className="char-counter">
                {value.length}/{maxLength}
            </div>
        </div>
    );
};

export default TextareaCounter;