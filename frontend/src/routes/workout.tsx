import { useEffect, useId, useState } from "react";
import { Exercise, Workout, Set, ExerciseType } from "../models/workout";
import { useNavigate, useParams } from "react-router";
import ApiService from "../services/api-service";

type ExerciseProps = {
    exercise: Exercise,
};

type EditableExerciseProps = {
    exercise: Exercise,
    deleteExerciseFunc: Function
};

const fetchSets = async (wId: string, eId: string, setSets: React.Dispatch<React.SetStateAction<Set[]>>) => {
    const res = await ApiService.fetchSets(wId, eId);

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
                        <li key={ex.id + " " + i}>{set.weight}kg for {set.repetitions} reps</li>
                    );
                })}
            </ul>
        </li>
    );
};

const EditableExercise = (props: EditableExerciseProps) => {
    const [ex] = useState<Exercise>(props.exercise);
    const [weight, setWeight] = useState<number>(0);
    const [sets, setSets] = useState<Set[]>([]);
    const [reps, setReps] = useState<number>(0);

    const weightId = useId();
    const repsId = useId();

    useEffect(() => {
        fetchSets(ex.workout_id, ex.id, setSets);
    }, [ex]);

    const deleteSet = async (setId: string) => {
        const res = await ApiService.deleteSet(ex.workout_id, ex.id, setId);

        if (res.status !== 204) {
            console.log("Error");
            return
        }

        setSets(s => s.filter(item => item.id !== setId));
    };

    const addSet = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const res = await ApiService.createSet(ex.workout_id, ex.id, reps, weight);

        if (res.status !== 201) {
            console.log("Error");
            return
        }

        const obj = await res.json();
        setSets([...sets, { id: obj.id, weight: weight, repetitions: reps, exercise_id: ex.id }]);
    };

    return (
        <>
            <li key={ex.id}>{ex.name}<button onClick={async () => { await props.deleteExerciseFunc(ex.id) }}>Delete</button>
                <ul>
                    {sets.map((set, i) => {
                        return (
                            <li key={ex.id + " " + i}>{set.weight}kg for {set.repetitions} reps<button onClick={() => deleteSet(set.id)}>Delete</button></li>
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
    const [exerciseTypeId, setExerciseTypeId] = useState<string | null>(null);

    const id = params.id;

    const navigate = useNavigate();

    useEffect(() => {
        const fetchExerciseTypes = async () => {
            const res = await ApiService.fetchExerciseTypes();
            if (res.status === 200) {
                const resObj = await res.json();
                setExerciseTypes(resObj.exercise_types);
            }
        };

        fetchExerciseTypes();
    }, []);

    const fetchWorkout = async () => {
        if (id === undefined) {
            return;
        }

        const res = await ApiService.fetchWorkout(id);
        if (res.status === 200) {
            const resObj = await res.json();
            setWorkout(resObj.workout);
        }
    };

    useEffect(() => {
        fetchWorkout();
    }, []);

    const deleteWorkout = async () => {
        if (workout === null) {
            return;
        }

        const res = await ApiService.deleteWorkout(workout.id);
        if (res.status !== 204) {
            console.log("Error", res.status, res.statusText);
            return;
        }

        navigate("/workouts");
    };

    const fetchExercises = async () => {
        if (id === undefined) {
            return;
        }

        const res = await ApiService.fetchExercises(id);
        if (res.status === 200) {
            const resObj = await res.json();
            setExercises(resObj.exercises);
        }
    };

    useEffect(() => {
        fetchExercises();
    }, []);

    const deleteExercise = async (exerciseId: string) => {
        if (id === undefined) {
            return;
        }

        const res = await ApiService.deleteExercise(id, exerciseId);
        if (res.status !== 204) {
            console.log("Error");
            return;
        }

        setExercises(l => l.filter(item => item.id !== exerciseId));
    };

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
                <button onClick={deleteWorkout}>Delete workout</button>
            </>
        );
    }

    const addExercise = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        if (exerciseTypeId !== "None" && exerciseTypeId !== null) {
            const exerciseType = exerciseTypes.filter(et => et.id == exerciseTypeId)[0];

            const res = await ApiService.createExercise(workout.id, exerciseType.id);

            if (res.status !== 201) {
                console.log("Error");
                return
            }

            const obj = await res.json();

            setExercises([...exercises, { id: obj.id, exercise_type_id: exerciseType.id, workout_id: workout.id, name: exerciseType.name }]);
            return;
        }

        if (exerciseName === "") {
            return;
        }

        const exerciseTypeRes = await ApiService.createExerciseType(exerciseName);

        if (exerciseTypeRes.status !== 201) {
            console.log("Error");
            return;
        }

        let obj = await exerciseTypeRes.json();
        setExerciseTypes([...exerciseTypes, { id: obj.id, name: exerciseName }]);

        const res = await ApiService.createExercise(workout.id, obj.id);

        if (res.status !== 201) {
            console.log("Error");
            return
        }

        obj = await res.json();

        setExercises([...exercises, { id: obj.id, exercise_type_id: obj.id, workout_id: workout.id, name: exerciseName }]);

        setExerciseName("");
    };


    const finishWorkout = async () => {
        const res = await ApiService.finishWorkout(workout.id);

        if (res.status !== 204) {
            console.log("Error", res.status, res.statusText);
            return;
        }

        await fetchWorkout();
    };

    return (
        <>
            <h1>Workout {workout.name}</h1>
            <h2>{new Date(workout.created_on).toDateString()}</h2>
            <button onClick={finishWorkout}>Finish workout</button>
            <h3>Exercises</h3>
            <ul>
                {exercises.map(e => {
                    return (<EditableExercise key={e.id} exercise={e} deleteExerciseFunc={deleteExercise} />);
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

            <button onClick={deleteWorkout}>Delete workout</button>
        </>
    );
};

export default WorkoutComponent;
