import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useDiscord } from "@/lib/discord-context";
import { BotIcon } from "lucide-react";

export default function DiscordStatus() {
  const application = useDiscord();

  return (
    <Card className="w-[300px]">
      <CardContent className="flex flex-row items-center space-x-4 py-4">
        <Avatar className="h-8 w-8">
          <AvatarImage src={application?.icon_url} alt="Discord Bot" />
          <AvatarFallback>
            <BotIcon className="h-6 w-6" />
          </AvatarFallback>
        </Avatar>
        <div className="flex-1 items-cen">
          <CardTitle>{application?.name ?? "JämK botten"}</CardTitle>
        </div>
        <div className="relative">
          <div
            className={`h-3 w-3 rounded-full ${
              !!application ? "bg-green-500" : "bg-red-500"
            }`}
          />
          <div
            className={`absolute inset-0 rounded-full ${
              !!application ? "bg-green-500" : "bg-red-500"
            } animate-ping`}
            style={{ animationDuration: "1.5s" }}
          />
        </div>
      </CardContent>
    </Card>
  );
}
