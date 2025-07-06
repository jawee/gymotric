import { ChartColumnIncreasing, Dumbbell } from "lucide-react";
import { Link } from "react-router";

const Landing = () => {
  return (
    <div className="font-sans flex flex-col min-h-screen bg-gray-100 text-gray-800">
      <header className="bg-gray-800 text-white p-6 md:px-10 flex justify-between items-center shadow-md">
        <h1 className="text-3xl md:text-4xl font-bold"><Dumbbell className="size-8 inline" /> Gymotric</h1>
        <nav className="space-x-4">
          <Link to="/login"
            className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg transition-colors duration-300"
          >
            Login
          </Link>
          <Link to="/register"
            className="bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded-lg transition-colors duration-300"
          >
            Register
          </Link>
        </nav>
      </header>

      <section className="flex-grow flex justify-center items-center text-center p-8 md:p-16 bg-white">
        <div className="max-w-4xl">
          <h2 className="text-4xl md:text-6xl font-extrabold mb-5 text-gray-900 leading-tight">
            Track Your Progress
          </h2>
          <p className="text-xl md:text-2xl leading-relaxed mb-8 text-gray-600">
            Our gym app makes it easy to log your workouts. Take
            control of your health today!
          </p>
          <Link to="/register"
            className="bg-red-500 hover:bg-red-600 text-white font-bold py-3 px-8 rounded-xl text-lg shadow-lg transform transition-transform duration-300 hover:scale-105"
          >
            Start Tracking Now!
          </Link>
        </div>
      </section>

      <section className="p-8 md:p-16 bg-gray-50 text-center border-t border-gray-200">
        <h3 className="text-3xl md:text-5xl font-bold mb-12 text-gray-900">
          Key Features
        </h3>
        <div className="flex flex-wrap justify-center gap-8">
          <div className="flex-shrink-0 w-full sm:w-80 bg-white p-8 rounded-xl shadow-lg text-left transform transition-transform duration-300 hover:scale-105 hover:shadow-xl">
            <h4 className="text-2xl font-semibold mb-3 text-gray-800">
              <Dumbbell className="inline" /> Track Workouts
            </h4>
            <p className="text-gray-600">
              No predefined workouts. Add and remove exercises as you go.
            </p>
          </div>
          <div className="flex-shrink-0 w-full sm:w-80 bg-white p-8 rounded-xl shadow-lg text-left transform transition-transform duration-300 hover:scale-105 hover:shadow-xl">
            <h4 className="text-2xl font-semibold mb-3 text-gray-800">
              <ChartColumnIncreasing className="inline" /> Statistics
            </h4>
            <p className="text-gray-600">
              Track your workouts over time.
            </p>
          </div>
        </div>
      </section>

      <footer className="bg-gray-800 text-gray-300 text-center p-5 mt-auto shadow-inner">
        <p>&copy; {new Date().getFullYear()} Gymotric. All rights reserved.</p>
      </footer>
    </div>
  );
};

export default Landing;
