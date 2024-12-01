import { useState } from "react";
import { Link } from "react-router";
import { workout } from "../models/workout";
import { dummyWorkouts } from "../models/dummy-data";

const WorkoutsList = () => {
    const [workouts] = useState<workout[]>(dummyWorkouts);
    return (
        <ul>
            {workouts.map(work => {
                return (
                    <li key={work.id}><Link to={"/workouts/" + work.id}>{work.name} {work.date.toISOString()}</Link></li>
                )
            })}
        </ul>
    );
}

export default WorkoutsList;
