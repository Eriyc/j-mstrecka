import React from "react";
import { createRoot } from "react-dom/client";
import "./style.css";
import App from "./app";
import { DiscordProvider } from "./lib/discord-context";
import { TransactionProvider } from "./lib/transaction-context";
import { ToastProvider } from "./lib/toast-context";

const container = document.getElementById("root");

const root = createRoot(container!);

root.render(
  <React.StrictMode>
    <ToastProvider>
      <DiscordProvider>
        <TransactionProvider>
          <App />
        </TransactionProvider>
      </DiscordProvider>
    </ToastProvider>
  </React.StrictMode>
);
