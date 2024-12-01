import Exercise from "../components/exercise";
import { dummyWorkouts } from "../models/dummy-data";

const Workout = () => {

    const workout = dummyWorkouts[0];

    return (
        <>
            <h1>Workout {workout.name}</h1>
            {workout.completed && <h2>{workout.date.toISOString()}</h2>}
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
};

export default Workout;
