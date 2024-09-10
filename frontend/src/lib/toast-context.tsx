import { createContext, PropsWithChildren, useContext, useState } from "react";
import { cn } from "./utils";
import { CheckCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

interface Toast {
  id: number;
  message: string;
  type: "error" | "success";
}

interface ToastContextType {
  toasts: Toast[];
  addToast: (message: string, type?: "error" | "success") => void;
}

const TOAST_DURATION = 5000; // 5 seconds

const ToastContext = createContext<ToastContextType>({} as ToastContextType);

export const ToastProvider: React.FC<PropsWithChildren> = ({ children }) => {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const removeToast = (id: number) => {
    setToasts((prev) => prev.filter((toast) => toast.id !== id));
  };

  const addToast = (message: string, type: "error" | "success" = "success") => {
    const id = Date.now();
    setToasts((prev) => [...prev, { id, message, type }]);
    setTimeout(() => {
      setToasts((prev) => prev.filter((toast) => toast.id !== id));
    }, TOAST_DURATION);
  };

  return (
    <ToastContext.Provider value={{ toasts, addToast }}>
      {children}
      <div className="fixed bottom-4 right-4 max-w-sm w-full flex flex-col items-end gap-2 p-4 pointer-events-none">
        {toasts.map((toast, index) => (
          <div
            key={toast.id}
            className={cn(
              "w-full p-4 rounded-lg shadow-lg text-white",
              toast.type === "error" && "bg-red-500",
              toast.type === "success" && "bg-green-500",
              "transform transition-all duration-500 ease-in-out",
              "translate-y-0 opacity-100"
            )}
            style={{
              zIndex: 1000 - index, // Ensure newer toasts appear on top
            }}
            role="alert"
          >
            <div className="flex items-center space-x-4">
              <CheckCircle className="w-6 h-6 text-white" />
              <div className="flex-1">
                <p className="font-medium">{toast.message}</p>
              </div>
            </div>
            <button
              onClick={() => removeToast(toast.id)}
              className="absolute top-2 right-2 text-white hover:text-opacity-75"
              aria-label="Close"
            >
              ×
            </button>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
};

export const useToast = () => useContext(ToastContext);
