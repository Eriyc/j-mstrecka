import { useCallback, useEffect, useRef } from "react";

export const useKeyboardListener = (
  callback: (upc: string) => Promise<void>
) => {
  const input = useRef<string>("");

  const listener = useCallback(
    async (e: KeyboardEvent) => {
      if (e.key === "Enter") {
        e.preventDefault();

        callback(input.current);

        input.current = "";
      } else if (e.key >= "0" && e.key <= "9") {
        e.preventDefault();

        input.current += e.key;
      }
    },
    [callback]
  );

  useEffect(() => {
    window.addEventListener("keydown", listener);

    return () => {
      window.removeEventListener("keydown", listener);
    };
  }, [listener]);
};
