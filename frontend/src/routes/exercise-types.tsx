import { Input } from "@/components/ui/input";
import { useEffect, useId, useState } from "react";
import { ExerciseType } from "../models/exercise-type";
import ApiService from "../services/api-service";
import { Button, buttonVariants } from "@/components/ui/button";
import { Pencil, Plus, Trash2 } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table";
import { cn } from "@/lib/utils";
import Loading from "@/components/loading";
import WtDialog from "@/components/wt-dialog";

const ExerciseTypes = () => {
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [exerciseTypes, setExerciseTypes] = useState<ExerciseType[]>([]);
  const exerciseNameId = useId();
  const [exerciseName, setExerciseName] = useState<string>("");

  const fetchExerciseTypes = async () => {
    const res = await ApiService.fetchExerciseTypes(); 
    if (res.status === 200) {
      const resObj = await res.json();
      setIsLoading(false);
      setExerciseTypes(resObj.data);
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

  if (isLoading) {
    return <Loading />
  }

  return (
    <>
      <h1 className="text-xl">Exercises</h1>
      <form onSubmit={addExercise}>
        Add new:
        <Input id={exerciseNameId} value={exerciseName} onChange={e => setExerciseName(e.target.value)} type="text" />
        <Button className="mt-2 mb-2" type="submit"><Plus />Add</Button>
      </form>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Exercise</TableHead>
            <TableHead></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {exerciseTypes.map((et) => {
            return (
              <ExerciseTypeRow 
                key={et.id} 
                exerciseType={et} 
                setExerciseTypes={setExerciseTypes}
              />
            );
          })}
        </TableBody>
      </Table>
    </>
  );
};

const ExerciseTypeRow = ({ exerciseType, setExerciseTypes }: { exerciseType: ExerciseType, setExerciseTypes: React.Dispatch<React.SetStateAction<ExerciseType[]>>  }) => {
  const [exerciseName, setExerciseName] = useState<string>(exerciseType.name);
  const [exerciseTypeName, setExerciseTypeName] = useState<string>(exerciseType.name);

  const handleNameChange = async () => {
    if (exerciseName === exerciseTypeName) {
      return;
    }

    const res = await ApiService.updateExerciseType(exerciseType.id, exerciseName);
    if (res.status !== 204) {
      console.log("Error updating exercise type");
      return;
    }

    setExerciseTypeName(exerciseName);
  };

  const deleteExerciseType = async (id: string) => {
    const confirm = window.confirm("Are you sure you want to delete this exercise?");
    if (!confirm) {
      return;
    }

    const res = await ApiService.deleteExerciseType(id);

    if (res.status !== 204) {
      console.log("Error");
    }

    setExerciseTypes(l => l.filter(item => item.id !== id));
  };


  return (
    <TableRow>
      <TableCell>{exerciseTypeName}</TableCell>
      <TableCell className="text-right">
        <WtDialog 
          openButtonTitle={<Pencil />} 
          form={<Input value={exerciseName} onChange={e => setExerciseName(e.target.value)} type="text" placeholder="Name" />} 
          onSubmitButtonClick={() => handleNameChange()} 
          onSubmitButtonTitle="Save" 
          title="Change name" />
        <Button className={
          cn(
            buttonVariants({ variant: "default" }),
            "ml-1",
            "bg-red-500",
            "hover:bg-red-700"
          )
        } onClick={() => deleteExerciseType(exerciseType.id)}>
          <Trash2 />
        </Button>
      </TableCell>
    </TableRow>
  );
}

export default ExerciseTypes;
