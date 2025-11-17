import { useEffect, useState } from "react";

type TimerDisplayProps = {
  isActive: boolean;
  lastSetTime: number | null;
};

const TimerDisplay = ({ isActive, lastSetTime }: TimerDisplayProps) => {
  const [elapsed, setElapsed] = useState<number>(0);

  useEffect(() => {
    if (!isActive || !lastSetTime) {
      setElapsed(0);
      return;
    }

    const interval = setInterval(() => {
      const now = Date.now();
      const secondsElapsed = Math.floor((now - lastSetTime) / 1000);
      setElapsed(secondsElapsed);
    }, 100);

    return () => clearInterval(interval);
  }, [isActive, lastSetTime]);

  const formatTime = (seconds: number) => {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;

    if (hours > 0) {
      return `${hours}:${String(minutes).padStart(2, "0")}:${String(secs).padStart(2, "0")}`;
    }
    return `${minutes}:${String(secs).padStart(2, "0")}`;
  };

  if (!isActive) {
    return null;
  }

  return (
    <div className="mt-3 text-lg text-gray-800 font-semibold">⏱️ Time since last set: {formatTime(elapsed)}</div>
  );
};

export default TimerDisplay;
