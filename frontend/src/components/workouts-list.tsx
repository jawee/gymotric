import { useEffect, useState } from "react";
import { Link } from "react-router";
import { Workout } from "../models/workout";

const WorkoutsList = () => {
    const [workouts, setWorkouts] = useState<Workout[]>();

    useEffect(() => {
        const fetchWorkouts = async () => {
           const res = await fetch("http://localhost:8080/workouts");
            if (res.status === 200) {
                const resObj = await res.json();
                setWorkouts(resObj.workouts);
            }
        };

        fetchWorkouts();
    }, []);

    if (workouts === null || workouts === undefined) {
        return (<p>Loading..</p>)
    }

    return (
        <ul>
            {workouts.map(workout => {
                return (
                    <li key={workout.id}><Link to={"/workouts/" + workout.id}>{new Date(workout.created_on).toDateString()}: {workout.name}</Link></li>
                )
            })}
        </ul>
    );
}

export default WorkoutsList;
