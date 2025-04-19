import { Input } from "@/components/ui/input";
import { useEffect, useId, useState } from "react";
import { useNavigate } from "react-router";
import { Workout } from "../models/workout";
import ApiService from "../services/api-service";
import WtDialog from "./wt-dialog";
import { Exercise } from "../models/exercise";
import Loading from "./loading";
import { Dumbbell } from "lucide-react";
import { Pagination, PaginationContent, PaginationItem, PaginationNext, PaginationPrevious } from "./ui/pagination";

const WorkoutsList = () => {
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [workouts, setWorkouts] = useState<Workout[]>([]);
  const [name, setName] = useState<string>("");
  const nameId = useId()

  const [pageSize] = useState<number>(10);
  const [page, setPage] = useState<number>(1);
  const [totalPages, setTotalPages] = useState<number>(0);

  const navigate = useNavigate();

  useEffect(() => {
    const fetchWorkouts = async () => {
      const res = await ApiService.fetchWorkouts(pageSize, page);

      if (res.status === 200) {
        const resObj = await res.json();
        setIsLoading(false);
        setWorkouts(resObj.data);
        setTotalPages(resObj.total_pages);
      }
    };

    fetchWorkouts();
  }, [pageSize, page]);

  const createWorkout = async () => {
    const res = await ApiService.createWorkout(name);

    if (res.status !== 201) {
      console.log("Error");
      return
    }

    const response = await res.json()

    navigate("/app/workouts/" + response.id);
  };

  if (isLoading) {
    return <Loading />
  }


  return (
    <div className="flex flex-col gap-4">
      <WtDialog openButtonTitle={<><Dumbbell /> Create workout</>} form={<Input id={nameId} value={name} onChange={e => setName(e.target.value)} type="text" placeholder="Name of workout" />} onSubmitButtonClick={createWorkout} onSubmitButtonTitle="Create workout" title="Create workout" />
      {workouts.map(workout => {
        return (
          <WorkoutListItem key={workout.id} workout={workout} />
        )
      })}
      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious className={page <= 1 ? "pointer-events-none opacity-50" : "cursor-pointer"} onClick={() => page !== 1 && setPage(page - 1)} />
          </PaginationItem>
          <PaginationItem>
            <PaginationNext className={page === totalPages ? "pointer-events-none opacity-50" : "cursor-pointer"} onClick={() => page !== totalPages && setPage(page + 1)} />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
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
        setExercises(resObj.data);
      }
    };

    fetchExercises();
  }, [workout.id]);

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
