import React, { useState, useEffect } from "react";
import { Link, useParams, useNavigate, BrowserRouter as Router, Routes, Route } from "react-router-dom";
import pdfMake from "pdfmake/build/pdfmake";
import { vfs } from "pdfmake/build/vfs_fonts";

pdfMake.vfs = vfs;

const API_BASE = "http://localhost:3001/api";

const last21Days = Array.from({ length: 21 }, (_, i) => {
    const d = new Date();
    d.setDate(d.getDate() - (20 - i));
    return d.toISOString().slice(0, 10);
});

function Home() {
    return (
        <div style={{ display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", height: "100vh" }}>
            <h1>Трекер привычек</h1>
            <Link to="/habits">Перейти к привычкам</Link>
        </div>
    );
}

function HabitPage({ onUpdate }) {
    const { name } = useParams();
    const navigate = useNavigate();

    const markToday = async () => {
        await fetch(`${API_BASE}/habits/${encodeURIComponent(name)}`, { method: "POST" });
        onUpdate();
    };

    return (
        <div style={{ padding: 20 }}>
            <h2>Привычка: {name}</h2>
            <button onClick={markToday} style={{ marginBottom: 20 }}>Отметить за сегодня</button>
            <br />
            <button onClick={() => navigate("/habits")}>Назад</button>
        </div>
    );
}

function HabitsList() {
    const [habits, setHabits] = useState({});
    const [newHabit, setNewHabit] = useState("");

    useEffect(() => {
        fetchHabits();
    }, []);

    const fetchHabits = async () => {
        const res = await fetch(`${API_BASE}/habits`);
        const data = await res.json();
        setHabits(data);
    };

    const addHabit = async () => {
        if (!newHabit.trim()) return;
        if (habits[newHabit.trim()]) {
            alert("Такая привычка уже есть");
            return;
        }
        await fetch(`${API_BASE}/habits`, {
            method: "POST",
            headers: { "Content-Type": "application/json; charset=utf-8" },
            body: JSON.stringify({ name: newHabit.trim() }),
        });
        setNewHabit("");
        fetchHabits();
    };

    const downloadPDF = () => {
        const headerRow = [
            { text: "Привычка / Дата", bold: true },
            ...last21Days.map((d) => ({ text: `${d.slice(8)}.${d.slice(5, 7)}`, bold: true })),
        ];

        const body = [headerRow];
        Object.entries(habits).forEach(([name, dates]) => {
            const row = [{ text: name }];
            last21Days.forEach((d) => row.push({ text: dates.includes(d) ? "+" : "" }));
            body.push(row);
        });
        if (body.length === 1) {
            body.push([
                { text: "Нет привычек, добавьте новую!", colSpan: 22, alignment: "center" },
                ...Array(21).fill({}),
            ]);
        }

        const docDefinition = {
            pageOrientation: 'landscape',
            pageSize: 'A4',
            content: [
                { text: "Трекер привычек", style: "header" },
                {
                    table: {
                        headerRows: 1,
                        widths: ['auto', ...Array(21).fill('*')],
                        body,
                    },
                },
            ],
            defaultStyle: {
                font: "Roboto",
                fontSize: 8,
            },
            styles: {
                header: { fontSize: 18, bold: true, margin: [0, 0, 0, 10] },
            },
        };

        pdfMake.createPdf(docDefinition).download("habits_tracker.pdf");
    };

    const printPDF = () => {
        const headerRow = [
            { text: "Привычка / Дата", bold: true },
            ...last21Days.map((d) => ({ text: `${d.slice(8)}.${d.slice(5, 7)}`, bold: true })),
        ];

        const body = [headerRow];
        Object.entries(habits).forEach(([name, dates]) => {
            const row = [{ text: name }];
            last21Days.forEach((d) => row.push({ text: dates.includes(d) ? "+" : "" }));
            body.push(row);
        });
        if (body.length === 1) {
            body.push([
                { text: "Нет привычек, добавьте новую!", colSpan: 22, alignment: "center" },
                ...Array(21).fill({}),
            ]);
        }

        const docDefinition = {
            pageOrientation: 'landscape',
            pageSize: 'A4',
            content: [
                { text: "Трекер привычек", style: "header" },
                {
                    table: {
                        headerRows: 1,
                        widths: ['auto', ...Array(21).fill('*')],
                        body,
                    },
                },
            ],
            defaultStyle: {
                font: "Roboto",
                fontSize: 8,
            },
            styles: {
                header: { fontSize: 18, bold: true, margin: [0, 0, 0, 10] },
            },
        };

        pdfMake.createPdf(docDefinition).print();
    };

    return (
        <div style={{ padding: 20, textAlign: "center" }}>
            <h2>Мои привычки</h2>
            <div style={{ marginBottom: 10 }}>
                <input
                    type="text"
                    placeholder="Новая привычка"
                    value={newHabit}
                    onChange={(e) => setNewHabit(e.target.value)}
                    style={{ padding: 5, marginRight: 10 }}
                />
                <button onClick={addHabit}>Добавить</button>
                <button onClick={downloadPDF} style={{ marginLeft: 10 }}>
                    Скачать PDF
                </button>
                <button onClick={printPDF} style={{ marginLeft: 10 }}>
                    Печать PDF
                </button>
            </div>

            <table border="1" cellPadding="5" style={{ margin: "0 auto", textAlign: "center" }}>
                <thead>
                <tr>
                    <th style={{ textAlign: "left" }}>Привычка / Дата</th>
                    {last21Days.map((d) => (
                        <th key={d}>{`${d.slice(8)}.${d.slice(5, 7)}`}</th>
                    ))}
                </tr>
                </thead>
                <tbody>
                {Object.entries(habits).map(([name, dates]) => (
                    <tr key={name}>
                        <td style={{ textAlign: "left" }}>
                            <Link to={`/habit/${encodeURIComponent(name)}`}>{name}</Link>
                        </td>
                        {last21Days.map((d) => (
                            <td key={d}>{dates.includes(d) ? "＋" : ""}</td>
                        ))}
                    </tr>
                ))}
                {Object.keys(habits).length === 0 && (
                    <tr>
                        <td colSpan={22}>Нет привычек, добавьте новую!</td>
                    </tr>
                )}
                </tbody>
            </table>
        </div>
    );
}

export default function App() {
    const [updateFlag, setUpdateFlag] = useState(false);
    const triggerUpdate = () => setUpdateFlag((f) => !f);

    return (
        <Router>
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/habits" element={<HabitsList key={updateFlag} />} />
                <Route path="/habit/:name" element={<HabitPage onUpdate={triggerUpdate} />} />
            </Routes>
        </Router>
    );
}
