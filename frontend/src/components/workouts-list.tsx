import { useState } from "react";
import { Link } from "react-router";
import { Workout } from "../models/workout";
import { dummyWorkouts } from "../models/dummy-data";

const WorkoutsList = () => {
    const [workouts] = useState<Workout[]>(dummyWorkouts);
    return (
        <ul>
            {workouts.map(workout => {
                return (
                    <li key={workout.id}><Link to={"/workouts/" + workout.id}>{workout.name} {workout.created_at.toISOString()}</Link></li>
                )
            })}
        </ul>
    );
}

export default WorkoutsList;
