import './LangSelector.css'
import { useState } from 'react';
import langIcon from '../assets/languageIcon.png'
import langSelect from './lang-select.svg'

function LangSelector() {
    const [lang, setLang] = useState('rus')
    return(
        <>
            <div className="lang-wrapper">
                <img src={langIcon} alt="planet" className="lang-icon"/>
                <div className="language-selector">
                    <p className="selected-lang">Русский</p>
                    <img src={langSelect} alt="selector" />
                </div>
            </div>
        </>
    );
}

export default LangSelector