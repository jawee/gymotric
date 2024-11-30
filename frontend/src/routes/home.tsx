import { useState } from 'react'

function Home() {
  const [count, setCount] = useState(0)

  const fetchData = () => {
    fetch('http://localhost:8080/')
      .then(response => response.text())
      .then(data => setMessage(data))
      .catch(error => console.error('Error fetching data:', error));
  };

  const [message, setMessage] = useState<string>('');

  return (
    <>
      <h1>Home</h1>
      <div className="card">
        <button onClick={() => setCount((count) => count + 1)}>
          count is {count}
        </button>
        <p>
          Edit <code>src/routes/home.tsx</code> and save to test HMR
        </p>
      </div>
      <button onClick={fetchData}>
        Click to fetch from Go server
      </button>
      {message && (
        <div>
          <h2>Server Response:</h2>
          <p>{message}</p>
        </div>
      )}
    </>
  );
}

export default Home;
