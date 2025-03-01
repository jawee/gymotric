import { useLocation } from "react-router";
import WorkoutsList from "../components/workouts-list";

interface WorkoutsState {
  error: string | null;
}

const workouts = () => {
    const location = useLocation();
    const state = location.state as WorkoutsState | null;

    return (
        <>
            <h1>Workouts</h1>
            {state?.error && <p>{state?.error}</p>}
            <WorkoutsList />
        </>
    );
};

export default workouts;
