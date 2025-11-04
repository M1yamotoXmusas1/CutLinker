import styles from './Link.module.css'

function Link() {
    
    return(
        <div className={styles.block}>
            <input className={styles.link} value="Вставьте вашу ссылку"/>
            <button className={styles.button}> Сократить </button>
        </div>
    )   

}

export default Link