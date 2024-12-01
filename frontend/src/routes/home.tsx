import WorkoutsList from '../components/workouts-list';

const Home = () => {
    // const [count, setCount] = useState(0)

    // const fetchData = () => {
    //   fetch('http://localhost:8080/')
    //     .then(response => response.text())
    //     .then(data => setMessage(data))
    //     .catch(error => console.error('Error fetching data:', error));
    // };

    // const [message, setMessage] = useState<string>('');

    const newWorkout = () => {
        console.log("new workout");
    };
    return (
        <>
            <h1>Home</h1>

            <button onClick={newWorkout}>New workout</button>
            <h3>Workouts</h3>
            <WorkoutsList />
        </>
    );
}

export default Home;
