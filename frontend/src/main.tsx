import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router'
import Home from './routes/home'
import Workout from './routes/workout';
import Workouts from './routes/workouts';
import Default from './routes/default';
import ExerciseTypes from './routes/exercise-types';

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <BrowserRouter>
            <Routes>
                <Route element={<Default />}>
                    <Route index element={<Home />} />
                    <Route path="/workouts" element={<Workouts />} />
                    <Route path="/workouts/:id" element={<Workout />} />
                    <Route path="/exercise-types" element={<ExerciseTypes />} />
                </Route>
            </Routes>
        </BrowserRouter>
    </StrictMode>
);
