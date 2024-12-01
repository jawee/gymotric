import { useState } from "react";
import { Link } from "react-router";
import { workout } from "../models/workout";
import { dummyWorkouts } from "../models/dummy-data";

const WorkoutsList = () => {
    const [workouts] = useState<workout[]>(dummyWorkouts);
    return (
        <>
            {workouts.map(work => {
                return (
                    <Link to={"/workouts/" + work.id} key={work.id}>{work.name} {work.date.toISOString()}</Link>
                )
            })}
        </>
    );
}

export default WorkoutsList;
