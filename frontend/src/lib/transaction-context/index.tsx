import React, { PropsWithChildren, useEffect, useState } from "react";
import {
  GetLatestTransactions,
  GetLeaderboard,
} from "@/../wailsjs/go/desktop/App";
import { models } from "wailsjs/go/models";
import { ChartDomain, Transaction, TransactionChartData } from "./types";
import { convertTransactionsToChartData } from "./transform";

type TransactionContextType = {
  transactions: {
    data: TransactionChartData;
    domain: ChartDomain;
    ticks: number[];
  };
  leaderboard: models.TransactionLeaderboard[];
  refetch: () => Promise<void>;
};

const TransactionContext = React.createContext<
  TransactionContextType | undefined
>(undefined);

export const TransactionProvider = ({ children }: PropsWithChildren) => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [leaderboard, setLeaderboard] = useState<
    models.TransactionLeaderboard[]
  >([]);

  const refetch = async () => {
    const response = await GetLatestTransactions();
    const leaderboard = await GetLeaderboard();

    if (!response) {
      return;
    }

    const transactions = response.map((transaction) => ({
      ...transaction,
      transaction_date: new Date(transaction.transaction_date),
    }));
    setLeaderboard(leaderboard);
    setTransactions(transactions);
  };

  useEffect(() => {
    refetch();

    window.runtime.EventsOn("transaction_updated", refetch);
    return () => {
      window.runtime.EventsOff("transaction_updated");
    };
  }, []);

  return (
    <TransactionContext.Provider
      value={{
        transactions: convertTransactionsToChartData(transactions),
        leaderboard,
        refetch,
      }}
    >
      {children}
    </TransactionContext.Provider>
  );
};

export const useTransactions = () => {
  const context = React.useContext(TransactionContext);

  if (context === undefined) {
    throw new Error(
      "useTransactions must be used within a TransactionProvider"
    );
  }

  return context;
};
