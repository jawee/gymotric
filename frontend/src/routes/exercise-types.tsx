import { Input } from "@/components/ui/input";
import { useEffect, useId, useState } from "react";
import { ExerciseType } from "../models/exercise-type";
import ApiService from "../services/api-service";
import { Button } from "@/components/ui/button";

const ExerciseTypes = () => {
  const [exerciseTypes, setExerciseTypes] = useState<ExerciseType[]>([]);

  const fetchExerciseTypes = async () => {
    const res = await ApiService.fetchExerciseTypes();
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

    const exerciseTypeRes = await ApiService.createExerciseType(exerciseName);

    if (exerciseTypeRes.status !== 201) {
      console.log("Error");
      return;
    }

    const obj = await exerciseTypeRes.json();
    setExerciseTypes([...exerciseTypes, { id: obj.id, name: exerciseName }]);
    setExerciseName("");
  };

  const deleteExerciseType = async (id: string) => {
    const res = await ApiService.deleteExerciseType(id);

    if (res.status !== 204) {
      console.log("Error");
    }

    setExerciseTypes(l => l.filter(item => item.id !== id));
  };

  const exerciseNameId = useId();
  const [exerciseName, setExerciseName] = useState<string>("");
  return (
    <>
      <h1 className="text-xl">Exercises</h1>
      <ul>
        {exerciseTypes.map((et) => {
          return (<li className="m-2" key={et.id}>{et.name} <Button onClick={() => deleteExerciseType(et.id)}>Delete exercise</Button></li>);
        })}
      </ul>
      <form onSubmit={addExercise}>
        Add new: <Input id={exerciseNameId} value={exerciseName} onChange={e => setExerciseName(e.target.value)} type="text" />
        <Button type="submit">Add exercise</Button>
      </form>
    </>
  );
};

export default ExerciseTypes;
