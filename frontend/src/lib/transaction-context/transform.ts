import { addHours } from "date-fns";
import { ChartDomain, Transaction, TransactionChartData } from "./types";

function convertTransactionsToChartData(transactions: Transaction[]): {
  data: TransactionChartData;
  domain: ChartDomain;
  ticks: number[];
} {
  // Sort transactions by date
  const sortedTransactions = [...transactions].sort(
    (a, b) => a.transaction_date.getTime() - b.transaction_date.getTime()
  );

  const users = [...new Set(transactions.map((item) => item.user_name))];

  const chartData: TransactionChartData = users.reduce<TransactionChartData>(
    (acc, user, index) => {
      acc[user] = {
        config: {
          label: user,
          color: `hsl(${(index * 360) / users.length}, 70%, 50%)`,
        },
        data: sortedTransactions
          .filter((transaction) => transaction.user_name === user)
          .map((transaction) => ({
            timestamp: new Date(transaction.transaction_date),
            value: transaction.cumulative_transaction_count,
          })),
      };
      return acc;
    },
    {}
  );

  const timestamps = Object.values(chartData).flatMap((userData) =>
    userData.data.map((dataPoint) => dataPoint.timestamp.getTime())
  );
  const minTimestamp = Math.min(...timestamps);
  const maxTimestamp = Math.max(...timestamps);

  const timeRange = maxTimestamp - minTimestamp;
  const padding = 0.1; // Adjust the padding factor as needed

  const startTime = new Date(minTimestamp - padding * timeRange);
  const endTime = new Date(maxTimestamp + padding * timeRange);

  const domain = {
    x: [startTime.getTime(), endTime.getTime()],
    y: [0, 16], // Adjust the y-domain based on your requirements
  } as ChartDomain;

  const ticks = [];
  const tickInterval = Math.ceil(timeRange / 10); // Adjust the number of ticks as needed

  let tickTime = Math.floor(minTimestamp / tickInterval) * tickInterval;
  while (tickTime <= maxTimestamp) {
    ticks.push(tickTime);
    tickTime += tickInterval;
  }

  return { data: chartData, domain, ticks };
}

export { convertTransactionsToChartData };
