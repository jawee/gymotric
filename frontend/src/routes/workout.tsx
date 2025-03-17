import { Input } from "@/components/ui/input";
import { useEffect, useState } from "react";
import { Workout } from "../models/workout";
import { Exercise } from "../models/exercise";
import { Set } from "../models/set";
import { ExerciseType } from "../models/exercise-type";
import { useNavigate, useParams } from "react-router";
import ApiService from "../services/api-service";
import { Button } from "@/components/ui/button";
import { buttonVariants } from "@/components/ui/button"

import { cn } from "@/lib/utils"
import React from "react";
import WtDialog from "../components/wt-dialog";

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
    <li key={ex.id}>{ex.name}
      <ul className="border-2 border-black m-2 p-2">
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
    <div className="border-2 border-black m-2 p-2">
      <li key={ex.id}>{ex.name} <Button onClick={deleteExercise}>Delete exercise</Button>
        <ul>
          {sets.map((set, i) => {
            return (
              <li key={ex.id + " " + i}>{set.weight}kg for {set.repetitions} reps <Button onClick={() => deleteSet(set.id)}>Delete set</Button></li>
            );
          })}
        </ul>
        <p className="font-bold">Add set</p>
        <form onSubmit={addSet} className="flex w-full max-w-sm items-center space-x-2">
          <Input id="weight" inputMode="decimal" type="text" pattern="^\d+([.,]0|[.,]5)?$" />
          <span className="mr-1">kg for</span>
          <Input id="reps" inputMode="numeric" type="number" />
          <span className="mr-1">reps</span>
          <Button className="" type="submit">Add</Button>
        </form>
        <p className="font-bold">Last set: {lastWeight}kg for {lastReps}reps</p>
        <p className="font-bold">Max set: {maxWeight}kg for {maxReps}reps</p>
      </li >
    </div>
  );
};

const WorkoutComponent = () => {
  const params = useParams();
  const [workout, setWorkout] = useState<Workout | null>(null);
  const [exercises, setExercises] = useState<Exercise[]>([])
  const [exerciseTypes, setExerciseTypes] = useState<ExerciseType[]>([]);

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
      return
    }

    navigate("/app/workouts", { state: { error: "Workout not found" } });
  };

  useEffect(() => {
    fetchWorkout();
  }, []);

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

  const [value, setValue] = useState("");

  if (workout === null) {
    return (
      <p>Loading</p>
    );
  }

  if (workout.completed_on !== null) {
    return (
      <>
        <h1 className="text-2xl">Workout {workout.name}</h1>
        <h2 className="text-l font-bold">{new Date(workout.created_on).toDateString()}</h2>
        <h3 className="text-2xl mt-3">Exercises</h3>
        <ul>
          {exercises.map(e => {
            return (
              <ExerciseComponent key={e.id} exercise={e} />
            );
          })}
        </ul>
        <Button onClick={deleteWorkout} className={cn(buttonVariants({ variant: "default" }), "bg-red-500", "hover:bg-red-700")}>Delete workout</Button>
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

      setExercises([...exercises, { id: obj.id, exercise_type_id: obj.id, workout_id: workout.id, name: value}]);

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

    await fetchWorkout();
  };


  return (
    <>
      <h1 className="text-2xl">Workout {workout.name}</h1>
      <h2 className="text-l font-bold">{new Date(workout.created_on).toDateString()}</h2>
      <Button onClick={finishWorkout}>Finish workout</Button>
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
      <div>
        <Button onClick={deleteWorkout} className={cn(buttonVariants({ variant: "default" }), "bg-red-500", "hover:bg-red-700")}>Delete workout</Button>
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
      const regex = new RegExp(`^${value}`, "i");
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
              <li key={i} onClick={onClick}>{suggestion}</li>
            );
          })}
        </ul>
      )}

    </div>
  );
};

export default WorkoutComponent;
