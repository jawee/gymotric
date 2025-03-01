const checkIfUnauthorized = (res: Response) => {
  if (res.status === 401) {
    window.location.href = "/login";
  }
};
const login = async (username: string, password: string) => {
  const res = await fetch("/api/login", {
    method: "POST",
    body: JSON.stringify({ username: username, password: password }),
  });

  return res;
};

const fetchWorkouts = async () => {
  const res = await fetch("/api/workouts", {
    credentials: "include",
  });

  checkIfUnauthorized(res);

  return res;
};

const fetchWorkout = async (id: string) => {
  const res = await fetch("/api/workouts/" + id, {
    credentials: "include",
  });
  checkIfUnauthorized(res);
  return res;
};

const createWorkout = async (name: string) => {
  const res = await fetch("/api/workouts", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ name: name })
  });
  checkIfUnauthorized(res);
  return res;
};

const finishWorkout = async (id: string) => {
  const res = await fetch("/api/workouts/" + id + "/complete", {
    method: "PUT",
    credentials: "include",
  });
  checkIfUnauthorized(res);
  return res;

};

const deleteWorkout = async (workoutId: string) => {
  const res = await fetch("/api/workouts/" + workoutId, {
    method: "DELETE",
    credentials: "include",
  });
  checkIfUnauthorized(res);
  return res;
};

const fetchSets = async (workoutId: string, exerciseId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets", {
    credentials: "include",
  });
  checkIfUnauthorized(res);
  return res;
};

const deleteSet = async (workoutId: string, exerciseId: string, setId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets/" + setId, {
    method: "DELETE",
    credentials: "include",
  });

  checkIfUnauthorized(res);

  return res;
};

const createSet = async (workoutId: string, exerciseId: string, repetitions: number, weight: number) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ repetitions: repetitions, weight: weight })
  });

  checkIfUnauthorized(res);

  return res;
};

const fetchExerciseTypes = async () => {
  const res = await fetch("/api/exercise-types", {
    credentials: "include",
  });
  checkIfUnauthorized(res);
  return res;
};

const createExerciseType = async (name: string) => {
  const res = await fetch("/api/exercise-types", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ name: name })
  });
  checkIfUnauthorized(res);
  return res;
};

const deleteExerciseType = async (id: string) => {
  const res = await fetch("/api/exercise-types/" + id, {
    method: "DELETE",
    credentials: "include",
  });
  checkIfUnauthorized(res);
  return res;
};

const fetchExercises = async (workoutId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises", {
    credentials: "include",
  });
  checkIfUnauthorized(res);
  return res;
};

const createExercise = async (workoutId: string, exerciseTypeId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ exercise_type_id: exerciseTypeId })
  });
  checkIfUnauthorized(res);
  return res;
};

const deleteExercise = async (workoutId: string, exerciseId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId, {
    method: "DELETE",
    credentials: "include",
  });
  checkIfUnauthorized(res);
  return res;
};

const ApiService = {
  login,
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
