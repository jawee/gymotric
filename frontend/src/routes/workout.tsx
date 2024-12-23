import { useEffect, useId, useState } from "react";
import { Exercise, Workout, Set, ExerciseType } from "../models/workout";
import { useParams } from "react-router";

type ExerciseProps = {
    exercise: Exercise
};

const fetchSets = async (wId: string, eId: string, setSets: React.Dispatch<React.SetStateAction<Set[]>>) => {
    const res = await fetch("http://localhost:8080/workouts/" + wId + "/exercises/" + eId + "/sets");
    if (res.status === 200) {
        const resObj = await res.json();
        setSets(resObj.sets);
    }
};

const ExerciseComponent = (props: ExerciseProps) => {
    const [ex] = useState<Exercise>(props.exercise);
    const [sets, setSets] = useState<Set[]>([]);

    useEffect(() => {
        fetchSets(ex.workout_id, ex.id, setSets);
    }, [ex]);

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
    const [ex] = useState<Exercise>(props.exercise);
    const [weight, setWeight] = useState<number>(0);
    const [sets, setSets] = useState<Set[]>([]);
    const [reps, setReps] = useState<number>(0);

    const weightId = useId();
    const repsId = useId();

    useEffect(() => {
        fetchSets(ex.workout_id, ex.id, setSets);
    }, [ex]);

    const addSet = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        // setExercise({ ...exercise, sets: [...exercise.sets, { weight: weight, reps: reps }] });
    };

    return (
        <>
            <li key={ex.id}>{ex.name}
                <ul>
                    {sets.map((set, i) => {
                        return (
                            <li key={ex.id + " " + i}>{set.weight}kg for {set.reps} reps</li>
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
    const params = useParams();
    const [workout, setWorkout] = useState<Workout | null>(null);
    const [exercises, setExercises] = useState<Exercise[]>([])
    const [exerciseTypes, setExerciseTypes] = useState<ExerciseType[]>([]);
    //form
    const [exerciseName, setExerciseName] = useState<string>("");
    const [exerciseTypeId, setExerciseTypeId] = useState<string>();

    const id = params.id;

    useEffect(() => {
        const fetchExerciseTypes = async () => {
            const res = await fetch("http://localhost:8080/exercise-types");
            if (res.status === 200) {
                const resObj = await res.json();
                setExerciseTypes(resObj.exercise_types);
            }
        };

        fetchExerciseTypes();
    }, []);

    useEffect(() => {
        const fetchWorkout = async () => {
            const res = await fetch("http://localhost:8080/workouts/" + id);
            if (res.status === 200) {
                const resObj = await res.json();
                setWorkout(resObj.workout);
            }
        };

        fetchWorkout();
    }, []);

    useEffect(() => {
        const fetchExercises = async () => {
            const res = await fetch("http://localhost:8080/workouts/" + id + "/exercises");
            if (res.status === 200) {
                const resObj = await res.json();
                setExercises(resObj.exercises);
            }
        };

        fetchExercises();
    }, []);

    const exerciseNameId = useId();
    const existingExerciseTypeSelectName = "exerciseTypeSelect";

    if (workout === null) {
        return (
            <p>Loading</p>
        );
    }


    if (workout.completed_on !== null) {
        return (
            <>
                <h1>Workout {workout.name}</h1>
                <h2>{new Date(workout.created_on).toDateString()}</h2>
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

    const addExercise = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        if (exerciseTypeId !== "None") {
            const exerciseType = exerciseTypes.filter(et => et.id == exerciseTypeId)[0];

            const res = await fetch("http://localhost:8080/workouts/" + workout.id + "/exercises", {
                method:"POST",
                body: JSON.stringify({ exercise_type_id: exerciseType.id })
            });
            if (res.status !== 201) {
                console.log("Error");
                return
            }
            const response = await res.json()

        }

        if (exerciseName === "") {
            return;
        }

        console.log(exerciseName)
        // setExercise({ ...exercise, sets: [...exercise.sets, { weight: weight, reps: reps }] });
        // setWorkout({...workout, exercises: [...workout.exercises, { id: 0, name: exerciseName, sets: [] }] });
        setExerciseName("");
    };


    return (
        <>
            <h1>Workout {workout.name}</h1>
            <h2>{new Date(workout.created_on).toDateString()}</h2>
            <h3>Exercises</h3>
            <ul>
                {exercises.map(e => {
                    return (<EditableExercise key={e.id} exercise={e} />);
                })}
            </ul>
            <form onSubmit={addExercise}>
                Add new: <input id={exerciseNameId} value={exerciseName} onChange={e => setExerciseName(e.target.value)} type="text" />
                <select name={existingExerciseTypeSelectName} onChange={e => setExerciseTypeId(e.target.value)}>
                    <option>None</option>
                    {exerciseTypes.map(e => {
                        return (<option key={e.id} value={e.id}>{e.name}</option>);
                    })}
                </select>
                <button type="submit">Add exercise</button>
            </form>
        </>
    );
};

export default WorkoutComponent;
