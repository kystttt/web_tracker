import { Link } from 'react-router-dom';

export default function Home() {
    return (
        <div>
            <h1>Это веб-приложение - трекер привычек</h1>
            <Link to="/login">Войти</Link> | <Link to="/register">Зарегистрироваться</Link>
        </div>
    );
}
