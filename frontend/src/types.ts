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
