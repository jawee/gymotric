import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router'
import Home from './routes/home'
import Workout from './routes/workout';
import Workouts from './routes/workouts';

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <BrowserRouter>
            <Routes>
                <Route index element={<Home />} />
                <Route path="/workouts" element={<Workouts />} />
                <Route path="/workouts/:id" element={<Workout />} />
            </Routes>
        </BrowserRouter>
    </StrictMode>
);
