export type Transaction = {
  user_id: string;
  user_name: string;
  // Go type: time
  transaction_date: Date;
  cumulative_transaction_count: number;
};

export type TransactionChartData = {
  [key: string]: {
    config: LineConfig;
    data: {
      timestamp: Date;
      value: number;
    }[];
  };
};

export type LineConfig = { label: string; color: string };
export type ChartDomain = { x: [number, number]; y: [number, number] };
