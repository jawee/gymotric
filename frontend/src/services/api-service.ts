const timeout = (ms: number) => {
  return new Promise(resolve => setTimeout(resolve, ms));
}

const checkIfUnauthorized = async (res: Response, isRetry: boolean = false): Promise<boolean> => {
  if (res.status === 401) {
    if (isRetry) {
      localStorage.setItem("refreshing", "false");
      window.location.href = "/login";
      return false
    }

    if (localStorage.getItem("refreshing") === "true") {
      await timeout(500);
      return true;
    }
    localStorage.setItem("refreshing", "true");
    const refreshRes = await refreshToken();
    localStorage.setItem("refreshing", "false");
    if (refreshRes.status === 200) {
      return true;
    }
    window.location.href = "/login";
  }

  return false
};

const login = async (username: string, password: string) => {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    body: JSON.stringify({ username: username, password: password }),
  });

  return res;
};

const refreshToken = async () => {
  const res = await fetch("/api/auth/token", {
    method: "POST",
    credentials: "include",
  });

  return res;
};

const fetchWorkouts = async (isRetry: boolean = false) => {
  const res = await fetch("/api/workouts", {
    credentials: "include",
  });

  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await fetchWorkouts(true);
  }

  return res;
};

const fetchWorkout = async (id: string, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts/" + id, {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await fetchWorkout(id, true);
  }
  return res;
};

const createWorkout = async (name: string, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ name: name })
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await createWorkout(name, true);
  }
  return res;
};

const finishWorkout = async (id: string, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts/" + id + "/complete", {
    method: "PUT",
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await finishWorkout(id, true);
  }
  return res;
};

const deleteWorkout = async (workoutId: string, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts/" + workoutId, {
    method: "DELETE",
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await deleteWorkout(workoutId, true);
  }
  return res;
};

const fetchSets = async (workoutId: string, exerciseId: string, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets", {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await fetchSets(workoutId, exerciseId, true);
  }
  return res;
};

const deleteSet = async (workoutId: string, exerciseId: string, setId: string, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets/" + setId, {
    method: "DELETE",
    credentials: "include",
  });

  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await deleteSet(workoutId, exerciseId, setId, true);
  }

  return res;
};

const createSet = async (workoutId: string, exerciseId: string, repetitions: number, weight: number, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ repetitions: repetitions, weight: weight })
  });

  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await createSet(workoutId, exerciseId, repetitions, weight, true);
  }

  return res;
};

const fetchExerciseTypes = async (isRetry: boolean = false) => {
  const res = await fetch("/api/exercise-types", {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await fetchExerciseTypes(true);
  }
  return res;
};

const createExerciseType = async (name: string, isRetry: boolean = false) => {
  const res = await fetch("/api/exercise-types", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ name: name })
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await createExerciseType(name, true);
  }
  return res;
};

const deleteExerciseType = async (id: string, isRetry: boolean = false) => {
  const res = await fetch("/api/exercise-types/" + id, {
    method: "DELETE",
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await deleteExerciseType(id, true);
  }
  return res;
};

const fetchExercises = async (workoutId: string, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises", {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await fetchExercises(workoutId, true);
  }
  return res;
};

const createExercise = async (workoutId: string, exerciseTypeId: string, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ exercise_type_id: exerciseTypeId })
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await createExercise(workoutId, exerciseTypeId, true);
  }
  return res;
};

const deleteExercise = async (workoutId: string, exerciseId: string, isRetry: boolean = false) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId, {
    method: "DELETE",
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await deleteExercise(workoutId, exerciseId, true);
  }
  return res;
};

const logout = async () => {
  const res = await fetch("/api/logout", {
    method: "POST",
    credentials: "include",
  });
  return res;
};

const fetchMaxWeightAndReps = async (exercise_type_id: string, isRetry: boolean = false) => {
  const res = await fetch("/api/exercise-types/" + exercise_type_id + "/max", {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await fetchExerciseTypes(true);
  }
  return res;
};
const fetchLastWeightAndReps = async (exercise_type_id: string, isRetry: boolean = false) => {
  const res = await fetch("/api/exercise-types/" + exercise_type_id + "/last", {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res, isRetry);
  if (shouldRetry && !isRetry) {
    return await fetchExerciseTypes(true);
  }
  return res;
};

const ApiService = {
  login,
  logout,
  refreshToken,
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
  fetchMaxWeightAndReps,
  fetchLastWeightAndReps
};

export default ApiService;
