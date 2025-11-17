import { Exercise } from "./exercise";

export type ExerciseItem = {
  id: string
  type: string
  user_id: string
  workout_id: string
  created_on: string
  updated_on: string
  exercises: Exercise[]
};
