export type Workout = {
    id: string
    name: string
    created_on: Date
    updated_on: Date
    completed_on: Date | null
};

export type Exercise = {
    id: string
    name: string
    workout_id: string
    exercise_type_id: string
};

export type Set = {
    id: string
    exercise_id: string
    weight: number
    repetitions: number
};

export type ExerciseType = {
    id: string
    name: string
};
