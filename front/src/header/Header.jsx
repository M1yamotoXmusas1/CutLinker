import styles from './Header.module.css'

function Header() {
    return(
        <header className={styles.header}>
            <p className={styles.name}>CutLink.io</p>
        </header>
    );
}

export default Header