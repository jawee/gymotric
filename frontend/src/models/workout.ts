export type workout = {
    id: number
    name: string
    date: Date
    completed: boolean
    exercises: exercise[]
};

export type exercise = {
    id: number
    name: string
    sets: set[]
};

export type set = {
    weight: number
    reps: number
};
