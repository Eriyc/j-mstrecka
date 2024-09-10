import { type ClassValue, clsx } from "clsx";
import { useEffect, useRef, useState } from "react";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const useCountdown = (callback: () => void) => {
  const [timeLeft, setTimeLeft] = useState(8);
  const countdown = useRef<NodeJS.Timeout | null>(null);

  // Effect to start the countdown when timeLeft is set
  useEffect(() => {
    if (timeLeft > 0 && !countdown.current) {
      countdown.current = setInterval(() => {
        setTimeLeft((t) => t - 1);
      }, 1000);
    }

    // Cleanup on unmount or timeLeft change
    return () => {
      console.log("unmount");
      clearInterval(countdown.current!);
    };
  }, []);

  useEffect(() => {
    if (timeLeft <= 0) {
      callback();
      cancel();
    }
  }, [timeLeft, callback]);

  const cancel = () => {
    clearInterval(countdown.current!);
    countdown.current = null;
  };

  const start = (time: number) => {
    cancel();
    setTimeLeft(time);

    countdown.current = setInterval(() => {
      setTimeLeft((t) => t - 1);
    }, 1000);
  };

  return [timeLeft, cancel, start] as const;
};
