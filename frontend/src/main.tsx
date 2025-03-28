import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router'
import Home from './routes/home'
import Default from './routes/default';
import Landing from './routes/landing';
import Login from './routes/login';
import Workout from './routes/workout';
import ExerciseTypes from './routes/exercise-types';
import Logout from './routes/logout';
import './index.css';
import Profile from './routes/profile';
createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Landing />} />
                <Route path="/login" element={<Login />} />
                <Route path="/app" element={<Default />}>
                    <Route index element={<Home />} />
                    <Route path="workouts/:id" element={<Workout />} />
                    <Route path="exercise-types" element={<ExerciseTypes />} />
                    <Route path="profile" element={<Profile />} />
                    <Route path="*" element={<Home />} />
                    <Route path="logout" element={<Logout />} />
                </Route>
            </Routes>
        </BrowserRouter>
    </StrictMode>
);
