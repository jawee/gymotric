import WorkoutsList from '../components/workouts-list';

const Home = () => {
  return (
    <>
      <h1 className="text-xl mb-2">Home</h1>
      <h2 className="text-l mb-2">Workouts</h2>
      <WorkoutsList />
    </>
  );
}

export default Home;
