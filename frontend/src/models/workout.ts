export type Workout = {
    id: string
    name: string
    created_at: Date
    updated_at: Date
    completed: boolean
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
    reps: number
};
