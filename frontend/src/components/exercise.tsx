import { useState } from "react";
import { exercise } from "../models/workout";

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

export default Exercise
