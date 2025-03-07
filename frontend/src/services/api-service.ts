const timeout = (ms: number) => {
  return new Promise(resolve => setTimeout(resolve, ms));
}

const checkIfUnauthorized = async (res: Response): Promise<boolean> => {
  if (res.status === 401) {
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

const fetchWorkouts = async () => {
  const res = await fetch("/api/workouts", {
    credentials: "include",
  });

  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await fetchWorkouts();
  }

  return res;
};

const fetchWorkout = async (id: string) => {
  const res = await fetch("/api/workouts/" + id, {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await fetchWorkout(id);
  }
  return res;
};

const createWorkout = async (name: string) => {
  const res = await fetch("/api/workouts", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ name: name })
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await createWorkout(name);
  }
  return res;
};

const finishWorkout = async (id: string) => {
  const res = await fetch("/api/workouts/" + id + "/complete", {
    method: "PUT",
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await finishWorkout(id);
  }
  return res;
};

const deleteWorkout = async (workoutId: string) => {
  const res = await fetch("/api/workouts/" + workoutId, {
    method: "DELETE",
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await deleteWorkout(workoutId);
  }
  return res;
};

const fetchSets = async (workoutId: string, exerciseId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets", {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await fetchSets(workoutId, exerciseId);
  }
  return res;
};

const deleteSet = async (workoutId: string, exerciseId: string, setId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets/" + setId, {
    method: "DELETE",
    credentials: "include",
  });

  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await deleteSet(workoutId, exerciseId, setId);
  }

  return res;
};

const createSet = async (workoutId: string, exerciseId: string, repetitions: number, weight: number) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId + "/sets", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ repetitions: repetitions, weight: weight })
  });

  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await createSet(workoutId, exerciseId, repetitions, weight);
  }

  return res;
};

const fetchExerciseTypes = async () => {
  const res = await fetch("/api/exercise-types", {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await fetchExerciseTypes();
  }
  return res;
};

const createExerciseType = async (name: string) => {
  const res = await fetch("/api/exercise-types", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ name: name })
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await createExerciseType(name);
  }
  return res;
};

const deleteExerciseType = async (id: string) => {
  const res = await fetch("/api/exercise-types/" + id, {
    method: "DELETE",
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await deleteExerciseType(id);
  }
  return res;
};

const fetchExercises = async (workoutId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises", {
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await fetchExercises(workoutId);
  }
  return res;
};

const createExercise = async (workoutId: string, exerciseTypeId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises", {
    method: "POST",
    credentials: "include",
    body: JSON.stringify({ exercise_type_id: exerciseTypeId })
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await createExercise(workoutId, exerciseTypeId);
  }
  return res;
};

const deleteExercise = async (workoutId: string, exerciseId: string) => {
  const res = await fetch("/api/workouts/" + workoutId + "/exercises/" + exerciseId, {
    method: "DELETE",
    credentials: "include",
  });
  const shouldRetry = await checkIfUnauthorized(res);
  if (shouldRetry) {
    return await deleteExercise(workoutId, exerciseId);
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
};

export default ApiService;
