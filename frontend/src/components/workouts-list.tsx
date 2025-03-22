import { Input } from "@/components/ui/input";
import { useEffect, useId, useState } from "react";
import { useNavigate } from "react-router";
import { Workout } from "../models/workout";
import ApiService from "../services/api-service";
import WtDialog from "./wt-dialog";
import { Exercise } from "../models/exercise";

const WorkoutsList = () => {
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [name, setName] = useState<string>("");
  const nameId = useId()

  const navigate = useNavigate();

  useEffect(() => {
    const fetchWorkouts = async () => {
      const res = await ApiService.fetchWorkouts();

      if (res.status === 200) {
        const resObj = await res.json();
        setWorkouts(resObj.workouts);
      }
    };

    fetchWorkouts();
  }, []);

  const createWorkout = async () => {
    const res = await ApiService.createWorkout(name);

    if (res.status !== 201) {
      console.log("Error");
      return
    }

    const response = await res.json()

    navigate("/app/workouts/" + response.id);
  };


  return (
    <div className="flex flex-col gap-4">
      {workouts.map(workout => {
        return (
          <WorkoutListItem key={workout.id} workout={workout} />
        )
      })}
      <WtDialog openButtonTitle="Create workout" form={<Input id={nameId} value={name} onChange={e => setName(e.target.value)} type="text" placeholder="Name of workout" />} onSubmitButtonClick={createWorkout} onSubmitButtonTitle="Create workout" title="Create workout" />
    </div>
  );
}

type WorkoutListItemProps = {
  workout: Workout
}

const WorkoutListItem = ({ workout }: WorkoutListItemProps) => {
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchExercises = async () => {
      const res = await ApiService.fetchExercises(workout.id);

      if (res.status === 200) {
        const resObj = await res.json();
        setExercises(resObj.exercises);
      }
    };

    fetchExercises();
  }, []);

  return (
    <div onClick={() => navigate("/app/workouts/" + workout.id)} className="cursor-pointer p-4 border border-gray-200
      flex flex-col md:flex-row gap-4">
      <div className="flex-1">
        <h3 className="font-medium text-xl">{workout.name}{(workout.completed_on === null) ? <span className="text-sm text-green-500"> In progress</span> : ""}</h3>
        <p>{new Date(workout.created_on).toLocaleString()}</p>
      </div>
      <div className="flex-1 hidden md:block">
        <p className="font-medium">Exercises:</p>
        <ul>
          {exercises.map(exercise => {
            return (
              <li key={exercise.id}>{exercise.name}</li>
            )
          })}
        </ul>
      </div>
    </div>
  );
};
export default WorkoutsList;
