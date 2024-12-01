export const dummyWorkouts = [
    {
        id: 1,
        name: "Legs",
        date: new Date(),
        completed: true,
        exercises: [
            {
                id: 1,
                name: "Squats",
                sets: [
                    {
                        weight: 85,
                        reps: 6,
                    },
                    {
                        weight: 85,
                        reps: 5,
                    }
                ]

            },
            {
                id: 2,
                name: "Leg extensions",
                sets: [
                    {
                        weight: 79,
                        reps: 11,
                    },
                    {
                        weight: 79,
                        reps: 11,
                    }
                ]

            }
        ]
    },
    {
        id: 2,
        name: "Back",
        date: new Date(),
        completed: false,
        exercises: [
            {
                id: 1,
                name: "Deadlifts",
                sets: [
                    {
                        weight: 110,
                        reps: 6,
                    },
                    {
                        weight: 110,
                        reps: 5,
                    }
                ]

            },
            {
                id: 2,
                name: "Barbell row",
                sets: [
                    {
                        weight: 75,
                        reps: 8,
                    },
                    {
                        weight: 75,
                        reps: 8,
                    }
                ]

            }
        ]
    }
];
