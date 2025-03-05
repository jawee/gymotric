import { Input } from "@/components/ui/input";
import { useEffect, useId, useState } from "react";
import { useNavigate } from "react-router";
import { Workout } from "../models/workout";
import ApiService from "../services/api-service";
import { Button } from "./ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

const CreateWorkoutForm = () => {
  const [name, setName] = useState<string>("");
  const nameId = useId()

  const navigate = useNavigate();

  const createWorkout = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const res = await ApiService.createWorkout(name);

    if (res.status !== 201) {
      console.log("Error");
      return
    }

    const response = await res.json()

    navigate("/app/workouts/" + response.id);
  };

  return (
    <>
      <form className="mt-5" onSubmit={createWorkout}>
        <Input id={nameId} value={name} onChange={e => setName(e.target.value)} type="text" placeholder="Name of workout" />
        <Button type="submit">Create workout</Button>
      </form>
    </>
  );
};

const WorkoutsList = () => {
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [isCreateWorkoutMode, setIsCreateWorkoutMode] = useState<boolean>(false);

  const navigate = useNavigate();

  const addWorkout = () => {
    setIsCreateWorkoutMode(true);
  };

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


  return (
    <>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Workout name</TableHead>
            <TableHead>Date</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {workouts.map(workout => {
            return (
              <TableRow key={workout.id} onClick={() => navigate("/app/workouts/" + workout.id)}>
                <TableCell className="font-medium">
                    {workout.name}
                </TableCell>
                <TableCell>
                  {new Date(workout.created_on).toDateString()}
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>
      {!isCreateWorkoutMode && <Button className="mt-5" onClick={addWorkout}>Create workout</Button>}
      {isCreateWorkoutMode && <CreateWorkoutForm />}
    </>
  );
}

export default WorkoutsList;
