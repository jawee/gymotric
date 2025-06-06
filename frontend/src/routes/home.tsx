import Statistics from '../components/statistics';
import WorkoutsList from '../components/workouts-list';

const Home = () => {
  return (
    <>
      <h1 className="text-xl mb-2">Welcome</h1>
      <Statistics showWeekly={true} showMonthly={false} showYearly={false} />
      <h2 className="text-xl mt-4 mb-2">Workouts</h2>
      <WorkoutsList />
    </>
  );
}

export default Home;
