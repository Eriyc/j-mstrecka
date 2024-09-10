import {
  CartesianGrid,
  Line,
  LineChart,
  XAxis,
  YAxis,
  ResponsiveContainer,
} from "recharts";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { format } from "date-fns";
import { ChartLegend } from "./ui/chart";
import {
  ChartDomain,
  TransactionChartData,
} from "@/lib/transaction-context/types";

type ChartProps = {
  data: TransactionChartData;
  domain: ChartDomain;
  ticks: number[];
};

export function DrinkChart({ data, domain, ticks }: ChartProps) {
  // Generate ticks for every 30 minutes

  console.log(data);

  return (
    <Card className="">
      <CardHeader>
        <CardTitle>Bästa alkoholisten idag</CardTitle>
        <CardDescription className="text-lg">
          Antal enheter över tid
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="h-96">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              margin={{
                top: 20,
                right: 30,
                left: 20,
                bottom: 60,
              }}
            >
              <CartesianGrid />
              <XAxis
                dataKey="timestamp"
                type="number"
                domain={[domain.x[0], domain.x[1]]}
                tickFormatter={(unixTime) =>
                  format(new Date(unixTime), "HH:mm")
                }
                ticks={ticks}
                interval={1}
              />
              <YAxis domain={domain.y} />

              {Object.entries(data).map(([key, line], index) => {
                return (
                  <Line
                    key={line.config.label}
                    type="basis"
                    strokeLinejoin="round"
                    data={line.data}
                    dataKey="value"
                    stroke={line.config.color}
                    strokeWidth={2}
                    dot={false}
                  />
                );
              })}

              <ChartLegend />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </CardContent>
      <CardFooter>
        <div className="flex items-center gap-2 leading-none text-slate-300 text-lg">
          Seså, gå och köp en enhet
        </div>
      </CardFooter>
    </Card>
  );
}
