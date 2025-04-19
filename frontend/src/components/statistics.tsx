import { Statistics } from "../models/statistics";
import { useEffect, useState } from "react";
import Loading from "./loading";
import ApiService from "@/services/api-service";

const StatisticsComponent = () => {
  const [stats, setStats] = useState<Statistics | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchStats = async () => {
      const res = await ApiService.fetchStatistics();

      if (res.status === 200) {
        const resObj = await res.json();
        setStats(resObj.data);
        setIsLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (isLoading || !stats) {
    return <Loading />;
  }
  return (
    <div>
      <p>Workouts this week: {stats.week}</p>
      <p>Workouts this month: {stats.month}</p>
      <p>Workouts this year: {stats.year}</p>
    </div>
  );
};

export default StatisticsComponent;
