import { useId, useState } from "react";
import { dummyExercises, dummySets, dummyWorkouts } from "../models/dummy-data";
import { Exercise, Workout, Set } from "../models/workout";

type ExerciseProps = {
    exercise: Exercise
};

const ExerciseComponent = (props: ExerciseProps) => {
    const [ex] = useState<Exercise>(props.exercise);
    const dSets: Set[] = dummySets.get(ex.id) ?? []
    const [sets] = useState<Set[]>(dSets);
    return (
        <li key={ex.id}>{ex.name}
            <ul>
                {sets.map((set, i) => {
                    return (
                        <li key={ex.id + " " + i}>{set.weight}kg for {set.reps} reps</li>
                    );
                })}
            </ul>
        </li>
    );
};

const EditableExercise = (props: ExerciseProps) => {
    const [exercise] = useState<Exercise>(props.exercise);
    const [weight, setWeight] = useState<number>(0);
    const dSets: Set[] = dummySets.get(exercise.id) ?? []
    const [sets] = useState<Set[]>(dSets);
    const [reps, setReps] = useState<number>(0);

    const weightId = useId();
    const repsId = useId();

    const addSet = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        // setExercise({ ...exercise, sets: [...exercise.sets, { weight: weight, reps: reps }] });
    };

    return (
        <>
            <li key={exercise.id}>{exercise.name}
                <ul>
                    {sets.map((set, i) => {
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

const WorkoutComponent = () => {
    const [workout] = useState<Workout>(dummyWorkouts[1]);
    const [exerciseName, setExerciseName] = useState<string>("");
    const ex = dummyExercises.get(workout.id) ?? [];
    const [exercises] = useState<Exercise[]>(ex)

    const exerciseNameId = useId();

    if (workout.completed) {
        return (
            <>
                <h1>Workout {workout.name}</h1>
                <h2>{workout.created_at.toDateString()}</h2>
                <h3>Exercises</h3>
                <ul>
                    {exercises.map(e => {
                        return (
                            <ExerciseComponent key={e.id} exercise={e} />
                        );
                    })}
                </ul>
            </>
        );
    }

    const addExercise = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        // setExercise({ ...exercise, sets: [...exercise.sets, { weight: weight, reps: reps }] });
        // setWorkout({...workout, exercises: [...workout.exercises, { id: 0, name: exerciseName, sets: [] }] });
        setExerciseName("");
    };


    return (
        <>
            <h1>Workout {workout.name}</h1>
            <h3>Exercises</h3>
            <ul>
                {exercises.map(e => {
                    return (<EditableExercise key={e.id} exercise={e} />);
                })}
            </ul>
            <form onSubmit={addExercise}>
                name: <input id={exerciseNameId} value={exerciseName} onChange={e => setExerciseName(e.target.value)} type="text" />
                <button type="submit">Add exercise</button>
            </form>
        </>
    );
};

export default WorkoutComponent;
