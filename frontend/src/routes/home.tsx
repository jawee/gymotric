import Statistics from '../components/statistics';
import WorkoutsList from '../components/workouts-list';

const Home = () => {
  return (
    <>
      <h2 className="text-xl mb-2">Statistics</h2>
      <Statistics />
      <h2 className="text-xl mt-2 mb-2">Workouts</h2>
      <WorkoutsList />
    </>
  );
}

export default Home;
