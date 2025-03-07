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
  TableRow,
} from "@/components/ui/table"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog"

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
    debugger;
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
      <Table>
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
                <TableCell>
                  {(workout.completed_on === null ? "In progress" : new Date(workout.completed_on).toDateString())}
                </TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>
      <Dialog>
        <DialogTrigger><Button>Create workout</Button></DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create workout</DialogTitle>
            <DialogDescription>
            </DialogDescription>
          </DialogHeader>
          <Input id={nameId} value={name} onChange={e => setName(e.target.value)} type="text" placeholder="Name of workout" />
          <DialogFooter>
            <Button onClick={createWorkout}>Create workout</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  );
}

export default WorkoutsList;
