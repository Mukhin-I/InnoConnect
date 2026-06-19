import './CreationInput.css';

const CreationInput = ({
    value,
    onChange,
    placeholder = '',
    className = '',
    icon = null,
    type = 'text',
    ...props
}) => {
    const iconStyle = icon ? {
        backgroundImage: `url(${icon})`,
        backgroundRepeat: 'no-repeat',
        backgroundPosition: 'left 1.33vw center',
        backgroundSize: '2.66vw',
        paddingLeft: '6.4vw',
    } : {};

    return (
        <input
            type={type}
            value={value}
            onChange={onChange}
            placeholder={placeholder}
            className={`input-icon ${className}`}
            style={iconStyle}
            {...props}
        />
    );
};

export default CreationInput;