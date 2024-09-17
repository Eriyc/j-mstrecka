import {
  createContext,
  PropsWithChildren,
  useContext,
  useEffect,
  useState,
} from "react";
import * as wails from "@wailsio/runtime";

type DiscordApplication = {
  icon_url: string;
  name: string;
};
const DiscordContext = createContext<DiscordApplication | null>(null);

export const DiscordProvider = ({ children }: PropsWithChildren) => {
  const [discord, setDiscord] = useState<DiscordApplication | null>(null);

  useEffect(() => {
    wails.Events.Emit({ name: "discord_check", data: {} });
    wails.Events.On(
      "discord_ready",
      ({ data }: { data: DiscordApplication | null }) => {
        console.log(data);

        setDiscord(data);
      }
    );

    return () => {
      wails.Events.Off("discord_ready");
    };
  }, []);

  return (
    <DiscordContext.Provider value={discord}>
      {children}
    </DiscordContext.Provider>
  );
};

export const useDiscord = () => {
  return useContext(DiscordContext);
};
