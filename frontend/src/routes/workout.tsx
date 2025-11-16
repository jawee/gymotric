import { Input } from "@/components/ui/input";
import { Dispatch, useEffect, useState } from "react";
import { Workout } from "../models/workout";
import { Exercise } from "../models/exercise";
import { ExerciseItem } from "../models/exercise-item";
import { Set } from "../models/set";
import { ExerciseType } from "../models/exercise-type";
import { useLocation, useNavigate, useParams } from "react-router";
import ApiService from "../services/api-service";
import { Button } from "@/components/ui/button";
import { buttonVariants } from "@/components/ui/button"

import { cn } from "@/lib/utils"
import React from "react";
import WtDialog from "../components/wt-dialog";
import { Check, Copy, Key, Plus, Trash2 } from "lucide-react";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import Loading from "../components/loading";
import { Textarea } from "@/components/ui/textarea";

type ExerciseProps = {
  exercise: Exercise,
};

type EditableExerciseProps = {
  exercise: Exercise,
  exerciseItemId: string,
  deleteExerciseFunc: (exerciseItemId: string, exerciseId: string) => Promise<void>;
};

const fetchSets = async (wId: string, eId: string, setSets: Dispatch<React.SetStateAction<Set[]>>) => {
  const res = await ApiService.fetchSets(wId, eId);

  if (res.status === 200) {
    const resObj = await res.json();
    setSets(resObj.data);
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

const EditableExercise = ({ exercise, exerciseItemId, deleteExerciseFunc }: EditableExerciseProps) => {
  const [sets, setSets] = useState<Set[]>([]);
  const [lastWeight, setLastWeight] = useState<number | null>(null);
  const [lastReps, setLastReps] = useState<number | null>(null);
  const [maxWeight, setMaxWeight] = useState<number | null>(null);
  const [maxReps, setMaxReps] = useState<number | null>(null);
  const [formData, setFormData] = useState<{ weight: string | null, reps: number | null }>({ weight: null, reps: null });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value,
    });

    localStorage.setItem("form-data-" + exercise.workout_id + "-" + exercise.id, JSON.stringify({
      ...formData,
      [name]: value,
    }));
  };

  useEffect(() => {
    const storedFormData = localStorage.getItem("form-data-" + exercise.workout_id + "-" + exercise.id);
    if (storedFormData) {
      setFormData(JSON.parse(storedFormData));
    }

    const fetchMaxWeightAndReps = async () => {
      const res = await ApiService.fetchMaxWeightAndReps(exercise.exercise_type_id);

      if (res.status === 200) {
        const resObj = await res.json();
        setMaxWeight(resObj.data.weight);
        setMaxReps(resObj.data.reps);
      }
    };

    const fetchLastWeightAndReps = async () => {
      const res = await ApiService.fetchLastWeightAndReps(exercise.exercise_type_id);

      if (res.status === 200) {
        const resObj = await res.json();
        setLastWeight(resObj.data.weight);
        setLastReps(resObj.data.reps);
      }
    };

    fetchMaxWeightAndReps();
    fetchLastWeightAndReps();
  }, [exercise]);

  useEffect(() => {
    fetchSets(exercise.workout_id, exercise.id, setSets);
  }, [exercise]);

  const deleteSet = async (setId: string) => {
    const res = await ApiService.deleteSet(exercise.workout_id, exercise.id, setId);

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

    const res = await ApiService.createSet(exercise.workout_id, exercise.id, reps, weight);

    if (res.status !== 201) {
      console.log("Error");
      return
    }

    const obj = await res.json();
    setSets([...sets, { id: obj.id, weight: weight, repetitions: reps, exercise_id: exercise.id }]);
  };

  const deleteExercise = async () => {
    await deleteExerciseFunc(exerciseItemId, exercise.id);
  };

  return (
    <div className="border-2 border-black p-2 mt-2 mb-2">
      <div>
        <p className="text-xl">{exercise.name}
          <WtDialog openButtonTitle={<Trash2 />}
            openButtonClassName={cn(buttonVariants({ variant: "default" }),
              "bg-red-500",
              "hover:bg-red-700",
              "ml-1",
            )
            } onSubmitButtonClick={() => deleteExercise()} title="Delete Set"
            form={<p>Are you sure you want to delete this exercise?</p>}
            onSubmitButtonTitle="Delete"
          />
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
                <TableRow key={exercise.id + "" + i}>
                  <TableCell>{set.weight}kg</TableCell>
                  <TableCell>{set.repetitions}</TableCell>
                  <TableCell className="text-right">
                    <WtDialog openButtonTitle={<Trash2 />} openButtonClassName={cn(buttonVariants({ variant: "default" }),
                      "bg-red-500",
                      "hover:bg-red-700",
                      "ml-1",
                    )
                    } onSubmitButtonClick={() => deleteSet(set.id)} title="Delete Set"
                      form={<p>Are you sure you want to delete this set?</p>}
                      onSubmitButtonTitle="Delete"
                    />
                  </TableCell>
                </TableRow>
              );
            })}
          </TableBody>
        </Table>
        <p className="mt-2">Add set</p>
        <div className="w-full">
          <form onSubmit={addSet} className="flex w-full max-w-md items-center space-x-2">
            <Input id="weight" name="weight" inputMode="decimal" type="text" pattern="^\d+([.,](00|0|25|50|5|75))?$" value={formData.weight ?? ""} onChange={handleChange} />
            <span className="flex-none mr-1">kg for</span>
            <Input id="reps" name="reps" inputMode="numeric" type="number" value={formData.reps ?? ""} onChange={handleChange} />
            <span className="mr-1">reps</span>
            <Button type="submit"><Plus />Add</Button>
          </form>
        </div>
        <p className="font-bold">Last set: {lastWeight}kg for {lastReps}reps</p>
        <p className="font-bold">Max set: {maxWeight}kg for {maxReps}reps</p>
      </div>
    </div>
  );
};

const WorkoutComponent = () => {
  const params = useParams();
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [workout, setWorkout] = useState<Workout | null>(null);
  const [exerciseItems, setExerciseItems] = useState<ExerciseItem[]>([])
  const [exerciseTypes, setExerciseTypes] = useState<ExerciseType[]>([]);

  const [finishDialogOpen, setFinishDialogOpen] = useState<boolean>(false);
  const [addExerciseDialogOpen, setAddExerciseDialogOpen] = useState<boolean>(false);
  const [itemTypeDialogOpen, setItemTypeDialogOpen] = useState<boolean>(false);
  const [selectedExerciseItemId, setSelectedExerciseItemId] = useState<string | null>(null);

  const [note, setNote] = useState<string>("");

  const location = useLocation();

  const id = params.id;

  const navigate = useNavigate();

  useEffect(() => {
    const fetchExerciseTypes = async () => {
      const res = await ApiService.fetchExerciseTypes();
      if (res.status === 200) {
        const resObj = await res.json();
        setExerciseTypes(resObj.data);
      }
    };

    fetchExerciseTypes();
  }, [workout]);


  useEffect(() => {
    const fetchWorkout = async () => {
      if (id === undefined) {
        return;
      }

      const res = await ApiService.fetchWorkout(id);
      if (res.status === 200) {
        const resObj = await res.json();
        setWorkout(resObj.data);
        setIsLoading(false);
        setNote(resObj.data.note);
        return
      }

      navigate("/app", { state: { error: "Workout not found" } });
    };
    fetchWorkout();
  }, [location, id, navigate]);

  const deleteWorkout = async () => {
    if (workout === null) {
      return;
    }

    const res = await ApiService.deleteWorkout(workout.id);
    if (res.status !== 204) {
      console.log("Error", res.status, res.statusText);
      return;
    }

    navigate("/app");
  };


  useEffect(() => {
    const fetchExerciseItems = async () => {
      if (id === undefined) {
        return;
      }

      const res = await ApiService.fetchExerciseItems(id);
      if (res.status === 200) {
        const resObj = await res.json();
        setExerciseItems(resObj.data);
      }
    };
    fetchExerciseItems();
  }, [workout, id]);

  const reopen = async () => {
    if (id === undefined) {
      return;
    }

    const res = await ApiService.reopenWorkout(id);
    if (res.status !== 204) {
      console.log("Error", res.status, res.statusText);
      return;
    }

    navigate("/app/workouts/" + id);
  };
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

  const deleteExercise = async (exerciseItemId: string, exerciseId: string) => {
    if (id === undefined) {
      return;
    }

    const res = await ApiService.deleteExercise(id, exerciseItemId, exerciseId);
    if (res.status !== 204) {
      console.log("Error");
      return;
    }

    // Remove the exercise from the exercise item that contains it
    setExerciseItems(items => {
      const updatedItems = items.map(item => ({
        ...item,
        exercises: item.exercises.filter(e => e.id !== exerciseId)
      })).filter(item => item.exercises.length > 0);

      // If the exercise item is now empty, delete it from the backend
      const itemToDelete = items.find(item => 
        item.id === exerciseItemId && 
        item.exercises.length === 1 && 
        item.exercises[0].id === exerciseId
      );
      
      if (itemToDelete) {
        ApiService.deleteExerciseItem(id, exerciseItemId).catch(err => {
          console.log("Error deleting empty exercise item:", err);
        });
      }

      return updatedItems;
    });
  };

  const [value, setValue] = useState("");
  const [selectedExerciseTypeId, setSelectedExerciseTypeId] = useState<string | null>(null);

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
        <Button onClick={reopen} className="ml-1"><Key />Reopen</Button>
        <h3 className="text-2xl mt-3">Exercises</h3>
        <div>
          {exerciseItems.map(item => (
            item.type === 'superset' ? (
              <div key={item.id} className="mt-4 p-3 border-2 border-yellow-500 bg-yellow-50 rounded">
                <h4 className="font-bold text-lg">🔗 Superset</h4>
                {item.exercises.map(e => (
                  <ExerciseComponent key={e.id} exercise={e} />
                ))}
              </div>
            ) : (
              <div key={item.id} className="mt-2">
                {item.exercises.map(e => (
                  <ExerciseComponent key={e.id} exercise={e} />
                ))}
              </div>
            )
          ))}
        </div>
        <div className="mt-2">
          <WtDialog openButtonTitle={<><Trash2 /> Delete workout</>}
            openButtonClassName={cn(buttonVariants({ variant: "default" }),
              "bg-red-500",
              "hover:bg-red-700",
            )
            } onSubmitButtonClick={() => deleteWorkout()} title="Delete Workout"
            form={<p>Are you sure you want to delete this Workout?</p>}
            onSubmitButtonTitle="Delete"
          />
        </div>
        <h3>Note</h3>
        <Textarea className="border-2" value={workout.note} disabled />
      </>
    );
  }

  const addExerciseToItem = async (exerciseTypeId: string, exerciseItemId: string | null, itemType: string) => {
    if (!workout) return;
    
    const exerciseType = exerciseTypes.find(et => et.id === exerciseTypeId);
    if (!exerciseType) return;

    let finalExerciseItemId = exerciseItemId;

    // If no exercise item provided, create one
    if (!finalExerciseItemId) {
      const exerciseItemRes = await ApiService.createExerciseItem(workout.id, itemType);
      if (exerciseItemRes.status !== 201) {
        console.log("Error creating exercise item");
        return;
      }
      const exerciseItemObj = await exerciseItemRes.json();
      finalExerciseItemId = exerciseItemObj.id;
    }

    // Create exercise with the exercise_item_id (finalExerciseItemId is now guaranteed to be a string)
    const res = await ApiService.createExercise(workout.id, exerciseTypeId, finalExerciseItemId as string);
    if (res.status !== 201) {
      console.log("Error creating exercise");
      return;
    }

    const exerciseObj = await res.json();

    // Update state
    setExerciseItems(items => {
      const existingItemIndex = items.findIndex(item => item.id === finalExerciseItemId);
      
      if (existingItemIndex >= 0) {
        // Add to existing item
        const updatedItems = [...items];
        updatedItems[existingItemIndex] = {
          ...updatedItems[existingItemIndex],
          exercises: [...updatedItems[existingItemIndex].exercises, {
            id: exerciseObj.id,
            exercise_type_id: exerciseTypeId,
            workout_id: workout.id,
            name: exerciseType.name
          }]
        };
        return updatedItems;
      } else {
        // Create new item
        const newExerciseItem: ExerciseItem = {
          id: finalExerciseItemId as string,
          type: itemType,
          user_id: "",
          workout_id: workout.id,
          created_on: new Date().toISOString(),
          updated_on: new Date().toISOString(),
          exercises: [{
            id: exerciseObj.id,
            exercise_type_id: exerciseTypeId,
            workout_id: workout.id,
            name: exerciseType.name
          }]
        };
        return [...items, newExerciseItem];
      }
    });

    setValue("");
    setAddExerciseDialogOpen(false);
    setItemTypeDialogOpen(false);
    setSelectedExerciseTypeId(null);
    setSelectedExerciseItemId(null);
  };

  const addExercise = async () => {
    if (value === "" || value === null) {
      return;
    }

    const exerciseTypeMatch = exerciseTypes.filter(et => et.name === value);

    if (exerciseTypeMatch.length === 1) {
      // Exercise type exists
      const exerciseTypeId = exerciseTypeMatch[0].id;
      
      // If adding to existing superset, add directly without type dialog
      if (selectedExerciseItemId) {
        const existingItem = exerciseItems.find(item => item.id === selectedExerciseItemId);
        if (existingItem) {
          addExerciseToItem(exerciseTypeId, selectedExerciseItemId, existingItem.type);
          return;
        }
      }
      
      // Otherwise show type selection dialog
      setSelectedExerciseTypeId(exerciseTypeId);
      setItemTypeDialogOpen(true);
      return;
    }

    // Create new exercise type
    const exerciseTypeRes = await ApiService.createExerciseType(value);
    if (exerciseTypeRes.status !== 201) {
      console.log("Error creating exercise type");
      return;
    }

    const exerciseTypeObj = await exerciseTypeRes.json();
    setExerciseTypes([...exerciseTypes, { id: exerciseTypeObj.id, name: value }]);
    
    // If adding to existing superset, add directly without type dialog
    if (selectedExerciseItemId) {
      const existingItem = exerciseItems.find(item => item.id === selectedExerciseItemId);
      if (existingItem) {
        addExerciseToItem(exerciseTypeObj.id, selectedExerciseItemId, existingItem.type);
        return;
      }
    }
    
    // Otherwise show type selection dialog
    setSelectedExerciseTypeId(exerciseTypeObj.id);
    setItemTypeDialogOpen(true);
  }

  const finishWorkout = async () => {
    const res = await ApiService.finishWorkout(workout.id);

    if (res.status !== 204) {
      console.log("Error", res.status, res.statusText);
      return;
    }

    setFinishDialogOpen(false);

    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i);
      if (key && key.startsWith("form-data-" + workout.id + "-")) {
        localStorage.removeItem(key);
      }
    }

    navigate("/app");
  };

  const updateNote = async () => {
    const res = await ApiService.updateWorkout(workout.id, note);
    if (res.status !== 204) {
      console.log("Error", res.status, res.statusText);
      return;
    }
  };

  return (
    <>
      <h1 className="text-2xl">Workout {workout.name}</h1>
      <h2 className="text-l font-bold">{new Date(workout.created_on).toDateString()}</h2>
      <div className="mt-2"><Button onClick={() => setFinishDialogOpen(true)}><Check />Finish workout</Button></div>
      <h3 className="text-2xl mt-3">Exercises</h3>
      <div>
        {exerciseItems.map(item => (
          item.type === 'superset' ? (
            <div key={item.id} className="mt-4 p-3 border-2 border-yellow-500 bg-yellow-50 rounded">
              <div className="flex justify-between items-center mb-2">
                <h4 className="font-bold text-lg">🔗 Superset</h4>
                <Button 
                  size="sm"
                  onClick={() => {
                    setSelectedExerciseItemId(item.id);
                    setItemTypeDialogOpen(false);
                    setAddExerciseDialogOpen(true);
                  }}
                >
                  <Plus /> Add to Superset
                </Button>
              </div>
              {item.exercises.map(e => (
                <EditableExercise key={e.id} exercise={e} exerciseItemId={item.id} deleteExerciseFunc={deleteExercise} />
              ))}
            </div>
          ) : (
            <div key={item.id} className="mt-2">
              {item.exercises.map(e => (
                <EditableExercise key={e.id} exercise={e} exerciseItemId={item.id} deleteExerciseFunc={deleteExercise} />
              ))}
            </div>
          )
        ))}
      </div>
      <div className="mt-4">
        <WtDialog 
          openButtonTitle={<><Plus /> Add Exercise</>} 
          form={
            <>
              <Autocomplete value={value} setValue={setValue} suggestions={exerciseTypes.map(et => et.name)} />
            </>
          } 
          onSubmitButtonClick={addExercise} 
          onSubmitButtonTitle={selectedExerciseItemId ? "Add" : "Next"} 
          title="Add Exercise"
          dialogProps={{ open: addExerciseDialogOpen, onOpenChange: setAddExerciseDialogOpen }}
        />
      </div>
      
      {/* Type Selection Dialog for new exercise items - programmatically opened */}
      {itemTypeDialogOpen && (
        <WtDialog
          openButtonTitle={null}
          form={
            <div className="space-y-4">
              <p>What type of exercise item do you want to create?</p>
              <div className="space-y-2">
                <Button 
                  onClick={() => {
                    if (selectedExerciseTypeId) {
                      addExerciseToItem(selectedExerciseTypeId, selectedExerciseItemId, 'exercise');
                    }
                  }}
                  className="w-full"
                >
                  Single Exercise
                </Button>
                <Button 
                  onClick={() => {
                    if (selectedExerciseTypeId) {
                      addExerciseToItem(selectedExerciseTypeId, selectedExerciseItemId, 'superset');
                    }
                  }}
                  className="w-full bg-yellow-500 hover:bg-yellow-600"
                >
                  Superset
                </Button>
              </div>
            </div>
          }
          onSubmitButtonTitle=""
          onSubmitButtonClick={() => {}}
          title="Exercise Item Type"
          hideSubmitButton={true}
          dialogProps={{ open: itemTypeDialogOpen, onOpenChange: setItemTypeDialogOpen }}
        />
      )}
      <div className="mt-2">
        <WtDialog openButtonTitle={<><Trash2 /> Delete workout</>}
          openButtonClassName={cn(buttonVariants({ variant: "default" }),
            "bg-red-500",
            "hover:bg-red-700"
          )
          } onSubmitButtonClick={() => deleteWorkout()} title="Delete Workout"
          form={<p>Are you sure you want to delete this Workout?</p>}
          onSubmitButtonTitle="Delete"
        />
      </div>
      <h3>Note</h3>
      <Textarea className="border-2" value={note} onChange={(e) => setNote(e.currentTarget.value)} onBlur={() => updateNote()} />
      <WtDialog onSubmitButtonClick={finishWorkout} title="Finish Workout" form={
        <p>Are you sure you want to finish this workout? You won't be able to add more exercises or sets.</p>
      } onSubmitButtonTitle="Finish Workout"
        dialogProps={{ open: finishDialogOpen, onOpenChange: setFinishDialogOpen }}
      />
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
