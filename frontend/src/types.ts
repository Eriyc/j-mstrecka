export type User = {
  id: string;
  name: string;
};

export type Balance = {
  total_credits_earned: number;
  total_payments_made: number;
  total_debt_incurred: number;
  remaining_credits: number;
  debt_incurred: number;
};

export type UserResponse = {
  type: "user";
  user: User;
  balance: Balance;
};

export type Product = {
  id: number;
  name: string;
  total_stock: number;
};

export type ProductPrice = {
  id: number;
  product_id: number;
  purchase_price: number;
  internal_price: number;
  external_price: number;
  start_date: string;
  end_date: string;
};

export type ProductResponse = {
  type: "product";
  product: Product;
  price: ProductPrice;
};

export type TransactionLeaderboard = {
  id: number;
  user_id: string;
  product_id: number;
  cumulative_transaction_count: number;
  transaction_count: number;
  transaction_date: Date;

  plater_rank: number;
  current_rank: number;
  rank_change_indicator: "↑" | "↓" | "=" | "New";
  user_name: string;
  total_transaction_count: number;
};

export type Transaction = {
  id: number;
  user_id: string;
  product_id: number;
  cumulative_transaction_count: number;
  transaction_count: number;
  transaction_date: Date;
};
