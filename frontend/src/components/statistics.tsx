import { Statistics } from "../models/statistics";
import { useEffect, useState } from "react";
import Loading from "./loading";
import ApiService from "@/services/api-service";

type StatisticsComponentProps = {
  showWeekly?: boolean;
  showMonthly?: boolean;
  showYearly?: boolean;
};
const StatisticsComponent = ({ showWeekly = true, showMonthly = true, showYearly = true } : StatisticsComponentProps) => {
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
      {showWeekly && <p>Workouts this week: {stats.week}</p>}
      {showWeekly && <p>Workouts last week: {stats.previous_week}</p>}
      {showMonthly && <p>Workouts this month: {stats.month}</p>}
      {showMonthly && <p>Workouts last month: {stats.previous_month}</p>}
      {showYearly && <p>Workouts this year: {stats.year}</p>}
      {showYearly && <p>Workouts last year: {stats.previous_year}</p>}
    </div>
  );
};

export default StatisticsComponent;
