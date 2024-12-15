import { Workout, Exercise, Set } from "./workout";

export const dummyWorkouts: Workout[] = [
    {
        id: "workout-1",
        name: "Legs",
        created_at: new Date(),
        updated_at: new Date(),
        completed: true,
    },
    {
        id: "workout-2",
        name: "Back",
        created_at: new Date(),
        updated_at: new Date(),
        completed: false,
    }
];
export const dummyExercises: Map<string, Exercise[]> = new Map<string, Exercise[]>([
    [
        "workout-1",
        [
            {
                id: "a",
                workout_id: "workout-1",
                name: "Squats",
                exercise_type_id: "squats"
            },
            {
                id: "b",
                name: "Leg extensions",
                workout_id: "workout-1",
                exercise_type_id: "leg-extensions"
            }
        ]
    ],
    [
        "workout-2",
        [
            {
                id: "c",
                name: "Deadlifts",
                workout_id: "workout-2",
                exercise_type_id: "deadlifts",
            },
            {
                id: "d",
                name: "Barbell row",
                workout_id: "workout-2",
                exercise_type_id: "barbell-row"
            }
        ]
    ]
]);

export const dummySets: Map<string, Set[]> = new Map<string, Set[]>([
    [
        "a",
        [
            {
                id: "a-set-1",
                reps: 2,
                weight: 12.5,
                exercise_id: "a",
            },
            {
                id: "a-set-2",
                reps: 2,
                weight: 12.5,
                exercise_id: "a",
            }
        ],
    ],
    [
        "b",
        [
            {
                id: "b-set-1",
                reps: 2,
                weight: 12.5,
                exercise_id: "b",
            },
            {
                id: "b-set-2",
                reps: 2,
                weight: 12.5,
                exercise_id: "b",
            }
        ],
    ],
    [
        "c",
        [
            {
                id: "c-set-1",
                reps: 2,
                weight: 12.5,
                exercise_id: "c",
            },
            {
                id: "c-set-2",
                reps: 2,
                weight: 12.5,
                exercise_id: "c",
            }
        ],
    ],
    [
        "d",
        [
            {
                id: "d-set-1",
                reps: 2,
                weight: 12.5,
                exercise_id: "d",
            },
            {
                id: "d-set-2",
                reps: 2,
                weight: 12.5,
                exercise_id: "d",
            }
        ],
    ]
]);
