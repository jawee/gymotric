import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, Route, Routes } from 'react-router'
import Home from './routes/home'
import Default from './routes/default';
import Landing from './routes/landing';
import Login from './routes/login';
import Workouts from './routes/workouts';
import Workout from './routes/workout';
import ExerciseTypes from './routes/exercise-types';
import Logout from './routes/logout';
import './index.css';
// const router = createBrowserRouter([
//   {
//     // errorElement: <ErrorPage />,
//     children: [
//       {
//         path: "/login",
//         element: <Login />,
//       },
//       // {
//       //   path: "/logout",
//       //   element: <LogOut />,
//       // },
//       {
//         path: "/",
//         element: <Landing />,
//       },
//       {
//         path: "/app",
//         // element: <ProtectedRoute />,
//         element: <Default />,
//         children: [
//           // {
//           //   path: "/partners/:id",
//           //   element: <PartnerView />,
//           // },
//           // {
//           //   path: "/partners/:partnerId/customers/:id",
//           //   element: <CustomerView />,
//           // },
//           // {
//           //   path: "/partners/:partnerId/users/:id",
//           //   element: <UserView />,
//           // },
//           // {
//           //   path: "/*",
//           //   element: <Home />,
//           // },
//           {
//             path: "*",
//             element: <Home />,
//           },
//         ]
//       }
//     ]
//   },
// ]);
createRoot(document.getElementById('root')!).render(
    <StrictMode>
      {/* <RouterProvider router={router} /> */}
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Landing />} />
                <Route path="/login" element={<Login />} />
                <Route path="/app" element={<Default />}>
                    <Route index element={<Home />} />
                    <Route path="workouts" element={<Workouts />} />
                    <Route path="workouts/:id" element={<Workout />} />
                    <Route path="exercise-types" element={<ExerciseTypes />} />
                    <Route path="*" element={<Home />} />
                    <Route path="logout" element={<Logout />} />
                </Route>
            </Routes>
        </BrowserRouter>
    </StrictMode>
);
