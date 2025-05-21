import { useEffect } from 'react';

function App() {
    useEffect(() => {
        fetch('/api/hello')
            .then(res => res.json())
            .then(data => console.log(data.message));
    }, []);

    return (
        <h1>Фронтенд работает</h1>
    );
}

export default App;
