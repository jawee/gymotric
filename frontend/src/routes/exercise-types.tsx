import { useEffect, useId, useState } from "react";
import { ExerciseType } from "../models/workout";

const ExerciseTypes = () => {
    const [exerciseTypes, setExerciseTypes] = useState<ExerciseType[]>([]);

    const fetchExerciseTypes = async () => {
        const res = await fetch("http://localhost:8080/exercise-types");
        if (res.status === 200) {
            const resObj = await res.json();
            setExerciseTypes(resObj.exercise_types);
        }
    };

    useEffect(() => {

        fetchExerciseTypes();
    }, []);

    const addExercise = async (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();

        if (exerciseName === "") {
            console.log("exerciseName empty");
            return;
        }

        if (exerciseTypes.filter(et => et.name === exerciseName).length > 0) {
            console.log("Already exists");
            return;
        }

        const exerciseTypeRes = await fetch("http://localhost:8080/exercise-types", {
            method: "POST",
            body: JSON.stringify({ name: exerciseName })
        });

        if (exerciseTypeRes.status !== 201) {
            console.log("Error");
            return;
        }

        const obj = await exerciseTypeRes.json();
        setExerciseTypes([...exerciseTypes, { id: obj.id, name: exerciseName }]);
    };

    const deleteExerciseType = async (id: string) => {
        const res = await fetch("http://localhost:8080/exercise-types/" + id, {
            method: "DELETE",
        });

        if (res.status !== 204) {
            console.log("Error");
        }

        setExerciseTypes(l => l.filter(item => item.id !== id));
    };

    const exerciseNameId = useId();
    const [exerciseName, setExerciseName] = useState<string>("");
    return (
        <>
            <h1>Exercise types</h1>
            <ul>
                {exerciseTypes.map((et) => {
                    return (<li key={et.id}>{et.name}<button onClick={() => deleteExerciseType(et.id)}>Delete</button></li>);
                })}
            </ul>
            <form onSubmit={addExercise}>
                Add new: <input id={exerciseNameId} value={exerciseName} onChange={e => setExerciseName(e.target.value)} type="text" />
                <button type="submit">Add exercise type</button>
            </form>
        </>
    );
};

export default ExerciseTypes;
