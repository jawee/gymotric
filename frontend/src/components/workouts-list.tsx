import { useEffect, useId, useState } from "react";
import { Link, useNavigate } from "react-router";
import { Workout } from "../models/workout";

const CreateWorkoutForm = () => {
    const [name, setName] = useState<string>("");
    const nameId = useId()

    const navigate = useNavigate();

    const createWorkout = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const res = await fetch("http://localhost:8080/workouts", {
            method: "POST",
            body: JSON.stringify({ name: name })
        });

        if (res.status !== 201) {
            console.log("Error");
        }

        const response = await res.json()

        navigate("/workouts/" + response.id);
    };

    return (
        <>
            <form onSubmit={createWorkout}>
                <input id={nameId} value={name} onChange={e => setName(e.target.value)} type="text" placeholder="Name of workout" />
                <button type="submit">Create workout</button>
            </form>
        </>
    );
};

const WorkoutsList = () => {
    const [workouts, setWorkouts] = useState<Workout[]>();
    const [isCreateWorkoutMode, setIsCreateWorkoutMode] = useState<boolean>(false);

    const addWorkout = () => {
        setIsCreateWorkoutMode(true);
    };
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
        <>
            <ul>
                {workouts.map(workout => {
                    return (
                        <li key={workout.id}><Link to={"/workouts/" + workout.id}>{new Date(workout.created_on).toDateString()}: {workout.name}</Link></li>
                    )
                })}
            </ul>
            {!isCreateWorkoutMode && <button onClick={addWorkout}>Create workout</button>}
            {isCreateWorkoutMode && <CreateWorkoutForm />}
        </>
    );
}

export default WorkoutsList;
