import Menu from '../components/menu';
import WorkoutsList from '../components/workouts-list';

const Home = () => {
    const newWorkout = () => {
        console.log("new workout");
    };

    return (
        <>
            <Menu />
            <h1>Home</h1>

            <button onClick={newWorkout}>New workout</button>
            <h3>Workouts</h3>
            <WorkoutsList />
        </>
    );
}

export default Home;
