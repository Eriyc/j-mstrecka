"'use client'";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useTransactions } from "@/lib/transaction-context";
import { ArrowDownIcon, ArrowUpIcon, Circle, MinusIcon } from "lucide-react";

export function Leaderboard() {
  const { leaderboard } = useTransactions();

  return (
    <Table className="flex-1">
      <TableHeader>
        <TableRow>
          <TableHead className="w-[100px]">Rank</TableHead>
          <TableHead>Player</TableHead>
          <TableHead className="text-right">Score</TableHead>
          <TableHead className="w-[100px] text-center">Change</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody className="text-white">
        {leaderboard.map((player) => (
          <TableRow key={player.user_id}>
            <TableCell className="font-medium">{player.current_rank}</TableCell>
            <TableCell>{player.user_name}</TableCell>
            <TableCell className="text-right">
              {player.total_transaction_count}
            </TableCell>
            <TableCell className="text-center">
              {player.rank_change_indicator === "↑" && (
                <ArrowUpIcon className="inline text-green-500" />
              )}
              {player.rank_change_indicator === "↓" && (
                <ArrowDownIcon className="inline text-red-500" />
              )}
              {player.rank_change_indicator === "=" && (
                <MinusIcon className="inline text-gray-500" />
              )}
              {player.rank_change_indicator === "New" && (
                <Circle className="inline text-green-500 fill-green-500" />
              )}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
