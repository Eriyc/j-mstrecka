import { UserResponse } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "./ui/card";
import { Badge } from "@/components/ui/badge";
import { useCountdown } from "@/lib/utils";
import { useEffect } from "react";

type Props = {
  user: UserResponse | null;
  onRemove: () => void;
};

export const UserInfo = ({ user, onRemove }: Props) => {
  const [timeLeft, cancel, start] = useCountdown(() => onRemove());

  useEffect(() => {
    start(8);

    return () => {
      cancel();
    };
  }, [user?.user.id]);

  if (!user) {
    return null;
  }

  return (
    <Card className="max-w-md w-full">
      <CardHeader>
        <CardTitle className="flex justify-between items-center">
          <span>
            {user.user.name} <span>({timeLeft} s)</span>
          </span>

          <Badge variant="outline">User</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <span className="font-semibold">Remaining Credits:</span>
            <span className="text-2xl font-bold text-green-600">
              {user.balance.remaining_credits} kr
            </span>
          </div>
          <div className="flex justify-between items-center">
            <span className="font-semibold">Debt Incurred:</span>
            <span className="text-2xl font-bold text-red-600">
              {user.balance.debt_incurred} kr
            </span>
          </div>
          <div className="text-sm text-muted-foreground">
            <p className="text-slate-300">
              Total Credits Earned:{" "}
              <span className="text-white">
                {user.balance.total_credits_earned} kr
              </span>
            </p>
            <p className="text-slate-300">
              Total Payments Made:{" "}
              <span className="text-white">
                {user.balance.total_payments_made} kr
              </span>
            </p>
            <p className="text-slate-300">
              Total Debt Incurred:{" "}
              <span className="text-white">
                {user.balance.total_debt_incurred} kr
              </span>
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
