import { useId, useState } from "react";
import { dummyWorkouts } from "../models/dummy-data";
import { exercise, workout } from "../models/workout";

type ExerciseProps = {
    exercise: exercise
};

const Exercise = (props: ExerciseProps) => {
    const [ex] = useState<exercise>(props.exercise);
    return (
        <li key={ex.id}>{ex.name}
            <ul>
                {ex.sets.map((set, i) => {
                    return (
                        <li key={ex.id + " " + i}>{set.weight}kg for {set.reps} reps</li>
                    );
                })}
            </ul>
        </li>
    );
};

const EditableExercise = (props: ExerciseProps) => {
    const [exercise, setExercise] = useState<exercise>(props.exercise);
    const [weight, setWeight] = useState<number>(0);
    const [reps, setReps] = useState<number>(0);

    const weightId = useId();
    const repsId = useId();

    const addSet = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        setExercise({...exercise, sets: [...exercise.sets, { weight: weight, reps: reps }]});
    };

    return (
        <>
            <li key={exercise.id}>{exercise.name}
                <ul>
                    {exercise.sets.map((set, i) => {
                        return (
                            <li key={exercise.id + " " + i}>{set.weight}kg for {set.reps} reps</li>
                        );
                    })}
                </ul>
                <form onSubmit={addSet}>
                    <input id={weightId} value={weight} onChange={e => setWeight(+e.target.value)} step=".5" type="number" />kg for <input value={reps} onChange={e => setReps(+e.target.value)} id={repsId} type="number" /> reps<br />
                    <button type="submit">Add set</button>
                </form>
            </li >
        </>
    );
};

const Workout = () => {
    const [workout] = useState<workout>(dummyWorkouts[1]);

    if (workout.completed) {
        return (
            <>
                <h1>Workout {workout.name}</h1>
                <h2>{workout.date.toDateString()}</h2>
                <h3>Exercises</h3>
                <ul>
                    {workout.exercises.map(e => {
                        return (
                            <Exercise key={e.id} exercise={e} />
                        );
                    })}
                </ul>
            </>
        );
    }

    return (
        <>
            <h1>Workout {workout.name}</h1>
            <h3>Exercises</h3>
            <ul>
                {workout.exercises.map(e => {
                    return (<EditableExercise key={e.id} exercise={e} />);
                })}
            </ul>
        </>
    );
};

export default Workout;
