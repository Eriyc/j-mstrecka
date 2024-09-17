import React, { PropsWithChildren, useEffect, useState } from "react";

import { ChartDomain, Transaction, TransactionChartData } from "./types";
import { convertTransactionsToChartData } from "./transform";
import * as wails from "@wailsio/runtime";

/* @ts-ignore */
import * as Service from "@/../bindings/gostrecka/services/transactions/transactionservice";
import { TransactionLeaderboard } from "@/types";

type TransactionContextType = {
  transactions: {
    data: TransactionChartData;
    domain: ChartDomain;
    ticks: number[];
  };
  leaderboard: TransactionLeaderboard[];
  refetch: () => Promise<void>;
};

const TransactionContext = React.createContext<
  TransactionContextType | undefined
>(undefined);

export const TransactionProvider = ({ children }: PropsWithChildren) => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [leaderboard, setLeaderboard] = useState<TransactionLeaderboard[]>([]);

  const refetch = async () => {
    const response = await Service.GetLatestTransactions();
    const leaderboard = await Service.GetLeaderboard();

    if (!response) {
      return;
    }

    const transactions = response.map((transaction: Transaction) => ({
      ...transaction,
      transaction_date: new Date(transaction.transaction_date),
    }));
    setLeaderboard(leaderboard);
    setTransactions(transactions);
  };

  useEffect(() => {
    refetch();

    wails.Events.On("transaction_updated", refetch);
    return () => {
      wails.Events.Off("transaction_updated");
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
