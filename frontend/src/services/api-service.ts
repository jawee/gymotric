
const fetchWorkouts = async () => {
    const res = await fetch("/api/workouts", {
        headers: {
            "ApiKey": "1234",
        },
    });
    return res;
};

const fetchWorkout = async (id: string) => {
    const res = await fetch("/api/workouts/" + id, {
        headers: {
            "ApiKey": "1234",
        },
    });
    return res;
};

const createWorkout = async (name: string) => {
    const res = await fetch("/api/workouts", {
        method: "POST",
        headers: {
            "ApiKey": "1234",
        },
        body: JSON.stringify({ name: name })
    });
    return res;
};

const finishWorkout = async (id: string) => {
    const res = await fetch("/api/workouts/" + id + "/complete", {
        method: "PUT",
        headers: {
            "ApiKey": "1234",
        },
    });
    return res;

};

const deleteWorkout = async (workoutId: string) => {
    const res = await fetch("/api/workouts/" + workoutId, {
        method: "DELETE",
        headers: {
            "ApiKey": "1234",
        },
    });
    return res;
};

const fetchSets = async (workoutId: string, exerciseId: string) => {
    const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets", {
        headers: {
            "ApiKey": "1234",
        },
    });
    return res;
};

const deleteSet = async (workoutId: string, exerciseId: string, setId: string) => {
    const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets/" + setId, {
        method: "DELETE",
        headers: {
            "ApiKey": "1234",
        },
    });

    return res;

};

const createSet = async (workoutId: string, exerciseId: string, repetitions: number, weight: number) => {
    const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets", {
        method: "POST",
        headers: {
            "ApiKey": "1234",
        },
        body: JSON.stringify({ repetitions: repetitions, weight: weight })
    });

    return res;
};

const fetchExerciseTypes = async () => {
    const res = await fetch("/api/exercise-types", {
        headers: {
            "ApiKey": "1234",
        },
    });
    return res;
};

const createExerciseType = async (name: string) => {
    const res = await fetch("/api/exercise-types", {
        method: "POST",
        headers: {
            "ApiKey": "1234",
        },
        body: JSON.stringify({ name: name })
    });
    return res;
};

const deleteExerciseType = async (id: string) => {
    const res = await fetch("/api/exercise-types/" + id, {
        method: "DELETE",
        headers: {
            "ApiKey": "1234",
        },
    });
    return res;

};

const fetchExercises = async (workoutId: string) => {
    const res = await fetch("/api/workouts/" + workoutId + "/exercises", {
        headers: {
            "ApiKey": "1234",
        },
    });
    return res;

};

const createExercise = async (workoutId: string, exerciseTypeId: string) => {
    const res = await fetch("/api/workouts/" + workoutId + "/exercises", {
        method: "POST",
        headers: {
            "ApiKey": "1234",
        },
        body: JSON.stringify({ exercise_type_id: exerciseTypeId })
    });
    return res;
};

const deleteExercise = async (workoutId: string, exerciseId: string) => {
    const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId, {
        method: "DELETE",
        headers: {
            "ApiKey": "1234",
        },
    });
    return res;
};

const ApiService = {
    fetchWorkouts,
    fetchWorkout,
    createWorkout,
    finishWorkout,
    deleteWorkout,
    fetchSets,
    deleteSet,
    createSet,
    fetchExerciseTypes,
    createExerciseType,
    deleteExerciseType,
    fetchExercises,
    createExercise,
    deleteExercise,
};

export default ApiService;
