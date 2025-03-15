import { Input } from "@/components/ui/input";
import { useEffect, useId, useState } from "react";
import { useNavigate } from "react-router";
import { Workout } from "../models/workout";
import ApiService from "../services/api-service";
import WtDialog from "./wt-dialog";

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
  const navigate = useNavigate();
  return (
    <div onClick={() => navigate("/app/workouts/" + workout.id)} className="cursor-pointer p-4 border border-gray-200">
      <h1 className="font-medium">{workout.name}</h1>
      <p>Date: {new Date(workout.created_on).toLocaleString()}</p>
      <p>{(workout.completed_on === null ? "In progress" : "Completed on:" + new Date(workout.completed_on).toLocaleString())}</p>
    </div>
  );
};
export default WorkoutsList;
