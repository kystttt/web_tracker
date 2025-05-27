import { useEffect, useState } from 'react';

export default function Dashboard({ onLogout }) {
    const [habits, setHabits] = useState([]);

    useEffect(() => {
        const token = localStorage.getItem('token');
        fetch('/api/habits', {
            headers: { Authorization: `Bearer ${token}` },
        })
            .then((res) => {
                if (res.status === 403) {
                    alert('Доступ запрещён');
                }
                return res.json();
            })
            .then((data) => setHabits(data));
    }, []);

    return (
        <div>
            <h2>Ваши привычки</h2>
            <button onClick={onLogout}>Выйти</button>
            <ul>
                {habits.map((habit) => (
                    <li key={habit.id}>{habit.name}</li>
                ))}
            </ul>
        </div>
    );
}
