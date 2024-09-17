import * as Runtime from "@wailsio/runtime";
import { Balance, Product, ProductResponse, User, UserResponse } from "./types";

declare global {
  interface Window {
    runtime: typeof Runtime;
  }
}

declare module "@/../bindings/gostrecka/services/transactions/transactionservice" {
  export interface TransactionService {
    GetLatestTransactions(): Promise<any>;
    GetLeaderboard(): Promise<any>;
    ScanUpc(upc: string): Promise<any>;
    Strecka(
      ProductID: number,
      UserID: string,
      amount: number
    ): Promise<{
      user: User;
      product: Product;
      balance: Balance;
      error: Error | null;
    }>;
  }
}
