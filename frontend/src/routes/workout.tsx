import { Input } from "@/components/ui/input";
import { useEffect, useState } from "react";
import { Workout } from "../models/workout";
import { Exercise } from "../models/exercise";
import { Set } from "../models/set";
import { ExerciseType } from "../models/exercise-type";
import { useLocation, useNavigate, useParams } from "react-router";
import ApiService from "../services/api-service";
import { Button } from "@/components/ui/button";
import { buttonVariants } from "@/components/ui/button"

import { cn } from "@/lib/utils"
import React from "react";
import WtDialog from "../components/wt-dialog";
import { Check, Copy, Plus, Trash2 } from "lucide-react";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import Loading from "../components/loading";

type ExerciseProps = {
  exercise: Exercise,
};

type EditableExerciseProps = {
  exercise: Exercise,
  deleteExerciseFunc: (exerciseId: string) => Promise<void>;
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
    <div className="border-2 border-black p-2 mt-2 mb-2">
      <p className="text-xl">{ex.name}</p>
      <Table className="mb-2">
        <TableHeader>
          <TableRow>
            <TableHead className="text-primary">Weight</TableHead>
            <TableHead className="text-primary">Reps</TableHead>
            <TableHead></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {sets.map((set, i) => {
            return (
              <TableRow key={ex.id + "" + i}>
                <TableCell>{set.weight}kg</TableCell>
                <TableCell>{set.repetitions}</TableCell>
                <TableCell className="text-right">
                </TableCell>
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
};

const EditableExercise = (props: EditableExerciseProps) => {
  const [ex] = useState<Exercise>(props.exercise);
  const [sets, setSets] = useState<Set[]>([]);
  const [lastWeight, setLastWeight] = useState<number | null>(null);
  const [lastReps, setLastReps] = useState<number | null>(null);
  const [maxWeight, setMaxWeight] = useState<number | null>(null);
  const [maxReps, setMaxReps] = useState<number | null>(null);

  const fetchMaxWeightAndReps = async () => {
    const res = await ApiService.fetchMaxWeightAndReps(ex.exercise_type_id);

    if (res.status === 200) {
      const resObj = await res.json();
      setMaxWeight(resObj.weight);
      setMaxReps(resObj.reps);
    }
  };

  const fetchLastWeightAndReps = async () => {
    const res = await ApiService.fetchLastWeightAndReps(ex.exercise_type_id);

    if (res.status === 200) {
      const resObj = await res.json();
      setLastWeight(resObj.weight);
      setLastReps(resObj.reps);
    }
  };

  useEffect(() => {
    fetchMaxWeightAndReps();
    fetchLastWeightAndReps();
  }, []);

  useEffect(() => {
    fetchSets(ex.workout_id, ex.id, setSets);
  }, [ex]);

  const deleteSet = async (setId: string) => {
    const confirmRes = confirm("Are you sure you want to delete this set?");
    if (!confirmRes) {
      return;
    }
    const res = await ApiService.deleteSet(ex.workout_id, ex.id, setId);

    if (res.status !== 204) {
      console.log("Error");
      return
    }

    setSets(s => s.filter(item => item.id !== setId));
  };

  const addSet = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const target = event.target as typeof event.target & {
      weight: { value: string };
      reps: { value: number };
    };

    const weight = +target.weight.value.replace(",", ".");
    const reps = +target.reps.value;

    const res = await ApiService.createSet(ex.workout_id, ex.id, reps, weight);

    if (res.status !== 201) {
      console.log("Error");
      return
    }

    const obj = await res.json();
    setSets([...sets, { id: obj.id, weight: weight, repetitions: reps, exercise_id: ex.id }]);
  };

  const deleteExercise = async () => {
    const res = confirm("Are you sure you want to delete this exercise?");
    if (!res) {
      return;
    }

    await props.deleteExerciseFunc(ex.id);
  };

  return (
    <div className="border-2 border-black p-2 mt-2 mb-2">
      <li key={ex.id}>
        <p className="text-xl">{ex.name}
          <Button
            onClick={deleteExercise}
            className={
              cn(buttonVariants({ variant: "default" }),
                "bg-red-500",
                "hover:bg-red-700",
                "ml-1",
              )
            }>
            <Trash2 />
          </Button>
        </p>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="text-primary">Weight</TableHead>
              <TableHead className="text-primary">Reps</TableHead>
              <TableHead></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sets.map((set, i) => {
              return (
                <TableRow key={ex.id + "" + i}>
                  <TableCell>{set.weight}kg</TableCell>
                  <TableCell>{set.repetitions}</TableCell>
                  <TableCell className="text-right">
                    <Button
                      onClick={() => deleteSet(set.id)}
                      className={
                        cn(buttonVariants({ variant: "default" }),
                          "bg-red-500",
                          "hover:bg-red-700",
                          "ml-1",
                        )
                      }>
                      <Trash2 />
                    </Button>
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
        <p className="mt-2">Add set</p>
        <div className="w-full">
          <form onSubmit={addSet} className="flex w-full max-w-md items-center space-x-2">
            <Input className="" id="weight" inputMode="decimal" type="text" pattern="^\d+([.,]0|[.,]5)?$" />
            <span className="flex-none mr-1">kg for</span>
            <Input className="" id="reps" inputMode="numeric" type="number" />
            <span className="mr-1">reps</span>
            <Button className="" type="submit"><Plus />Add</Button>
          </form>
        </div>
        <p className="font-bold">Last set: {lastWeight}kg for {lastReps}reps</p>
        <p className="font-bold">Max set: {maxWeight}kg for {maxReps}reps</p>
      </li >
    </div>
  );
};

const WorkoutComponent = () => {
  const params = useParams();
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [workout, setWorkout] = useState<Workout | null>(null);
  const [exercises, setExercises] = useState<Exercise[]>([])
  const [exerciseTypes, setExerciseTypes] = useState<ExerciseType[]>([]);

  const location = useLocation();

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
  }, [workout]);

  const fetchWorkout = async () => {
    if (id === undefined) {
      return;
    }

    const res = await ApiService.fetchWorkout(id);
    if (res.status === 200) {
      const resObj = await res.json();
      setWorkout(resObj.workout);
      setIsLoading(false);
      return
    }

    navigate("/app/workouts", { state: { error: "Workout not found" } });
  };

  useEffect(() => {
    fetchWorkout();
  }, [location]);

  const deleteWorkout = async () => {
    if (workout === null) {
      return;
    }

    const confirmRes = confirm("Are you sure you want to delete this workout?");
    if (!confirmRes) {
      return;
    }

    const res = await ApiService.deleteWorkout(workout.id);
    if (res.status !== 204) {
      console.log("Error", res.status, res.statusText);
      return;
    }

    navigate("/app/workouts");
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
  }, [workout]);

  const cloneWorkout = async () => {
    if (id === undefined) {
      return;
    }

    const res = await ApiService.cloneWorkout(id);
    if (res.status !== 201) {
      console.log("Error", res.status, res.statusText);
      return;
    }

    const obj = await res.json();
    navigate("/app/workouts/" + obj.id);
  };

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

  const [value, setValue] = useState("");

  if (workout === null || isLoading) {
    return (
      <Loading />
    );
  }

  if (workout.completed_on !== null) {
    return (
      <>
        <h1 className="text-2xl">Workout {workout.name}</h1>
        <h2 className="text-l font-bold">{new Date(workout.created_on).toDateString()}</h2>
        <Button onClick={cloneWorkout}><Copy />Clone</Button>
        <h3 className="text-2xl mt-3">Exercises</h3>
        <ul>
          {exercises.map(e => {
            return (
              <ExerciseComponent key={e.id} exercise={e} />
            );
          })}
        </ul>
        <div className="mt-2">
          <Button
            onClick={deleteWorkout}
            className={
              cn(buttonVariants({ variant: "default" }),
                "bg-red-500",
                "hover:bg-red-700"
              )
            }>
            <Trash2 />
            Delete workout
          </Button>
        </div>
      </>
    );
  }

  const addExercise = async () => {
    if (value !== "" && value !== null) {
      const exerciseTypeMatch = exerciseTypes.filter(et => et.name == value);

      if (exerciseTypeMatch.length === 1) {
        const exerciseType = exerciseTypeMatch[0];
        const res = await ApiService.createExercise(workout.id, exerciseType.id);

        if (res.status !== 201) {
          console.log("Error");
          return
        }

        setValue("");
        const obj = await res.json();

        setExercises([...exercises, { id: obj.id, exercise_type_id: exerciseType.id, workout_id: workout.id, name: exerciseType.name }]);
        return;
      }

      if (value === "") {
        return;
      }

      const exerciseTypeRes = await ApiService.createExerciseType(value);

      if (exerciseTypeRes.status !== 201) {
        console.log("Error");
        return;
      }

      let obj = await exerciseTypeRes.json();
      setExerciseTypes([...exerciseTypes, { id: obj.id, name: value }]);

      const res = await ApiService.createExercise(workout.id, obj.id);

      if (res.status !== 201) {
        console.log("Error");
        return
      }

      obj = await res.json();

      setExercises([...exercises, { id: obj.id, exercise_type_id: obj.id, workout_id: workout.id, name: value }]);

      setValue("");

      return;
    }
  }

  const finishWorkout = async () => {
    const confirmRes = confirm("Are you sure you want to finish this workout?");
    if (!confirmRes) {
      return;
    }

    const res = await ApiService.finishWorkout(workout.id);

    if (res.status !== 204) {
      console.log("Error", res.status, res.statusText);
      return;
    }

    navigate("/app");
  };


  return (
    <>
      <h1 className="text-2xl">Workout {workout.name}</h1>
      <h2 className="text-l font-bold">{new Date(workout.created_on).toDateString()}</h2>
      <div className="mt-2"><Button onClick={finishWorkout}><Check />Finish workout</Button></div>
      <h3 className="text-2xl mt-3">Exercises</h3>
      <ul>
        {exercises.map(e => {
          return (<EditableExercise key={e.id} exercise={e} deleteExerciseFunc={deleteExercise} />);
        })}
      </ul>
      <WtDialog openButtonTitle="Add Exercise" form={
        <>
          <Autocomplete value={value} setValue={setValue} suggestions={exerciseTypes.map(et => et.name)} />
        </>
      } onSubmitButtonClick={addExercise} onSubmitButtonTitle="Add exercise" title="Add Exercise" />
      <div className="mt-2">
        <Button
          onClick={deleteWorkout}
          className={
            cn(buttonVariants({ variant: "default" }),
              "bg-red-500",
              "hover:bg-red-700"
            )
          }>
          <Trash2 />
          Delete workout
        </Button>
      </div>
    </>
  );
};

type AutoCompleteProps = {
  value: string;
  setValue: React.Dispatch<React.SetStateAction<string>>;
  suggestions: string[];
};

const Autocomplete = ({ value, setValue, suggestions }: AutoCompleteProps) => {
  const [filteredSuggestions, setFilteredSuggestions] = useState<string[]>([]);
  const [showSuggestions, setShowSuggestions] = useState<boolean>(false);

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    let filteredSuggestions: string[] = [];

    if (value.length > 0) {
      const regex = new RegExp(`${value}`, "i");
      filteredSuggestions = suggestions.sort().filter(v => regex.test(v));
    }

    setFilteredSuggestions(filteredSuggestions);
    setShowSuggestions(true);
    setValue(value);
  };

  const onClick = (e: React.MouseEvent<HTMLLIElement>) => {
    setValue(e.currentTarget.innerText);
    setFilteredSuggestions([]);
    setShowSuggestions(false);
  };

  return (
    <div className="relative">
      <input type="text" value={value} onChange={onChange} className={cn(
        "w-full border-2 border-black rounded-md px-3 py-1 text-base outline-none",
      )} />
      {showSuggestions && value.length > 0 && (
        <ul className="absolute w-full z-10 bg-white border-2 border-black">
          {filteredSuggestions.map((suggestion, i) => {
            return (
              <Suggestion key={i} suggestion={suggestion} onClick={onClick} />
            );
          })}
        </ul>
      )}
    </div>
  );
};

const Suggestion = ({ suggestion, onClick }: { suggestion: string, onClick: (e: React.MouseEvent<HTMLLIElement>) => void }) => {
  return (
    <li onClick={onClick} className="p-2 cursor-pointer">{suggestion}</li>
  )
};

export default WorkoutComponent;
