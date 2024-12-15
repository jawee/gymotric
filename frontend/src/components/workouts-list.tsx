import { useState } from "react";
import { Link } from "react-router";
import { Workout } from "../models/workout";
import { dummyWorkouts } from "../models/dummy-data";

const WorkoutsList = () => {
    const [workouts] = useState<Workout[]>(dummyWorkouts);
    return (
        <ul>
            {workouts.map(work => {
                return (
                    <li key={work.id}><Link to={"/workouts/" + work.id}>{work.name} {work.created_at.toISOString()}</Link></li>
                )
            })}
        </ul>
    );
}

export default WorkoutsList;
