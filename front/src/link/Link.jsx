import styles from './Link.module.css'
import { useState } from 'react';

function Link() {

    const [link, setLink] = useState("");

    function linkChange(e) {
        setLink(e.target.value);
    }
    
    async function sendLink() {

        const response = await fetch('#URL', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ link })
            });
  
        if (response.ok) {
            alert('Ссылка отправлена в Go!');
        }
    }   

    return(
        <div className={styles.block}>
            <input 
                className={styles.link}  
                onChange={linkChange}
                placeholder="Введите текст..."
            />
            <button className={styles.button} onClick={sendLink}> Сократить </button>
        </div>
    )   

}



export default Link