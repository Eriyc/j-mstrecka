import {
  createContext,
  PropsWithChildren,
  useContext,
  useEffect,
  useState,
} from "react";

type DiscordApplication = {
  icon_url: string;
  name: string;
};
const DiscordContext = createContext<DiscordApplication | null>(null);

export const DiscordProvider = ({ children }: PropsWithChildren) => {
  const [discord, setDiscord] = useState<DiscordApplication | null>(null);

  useEffect(() => {
    const isOnWails = !!window.runtime;

    if (isOnWails) {
      window.runtime.EventsEmit("discord_check");
      window.runtime.EventsOn(
        "discord_ready",
        (bot: DiscordApplication | null) => {
          setDiscord(bot);
        }
      );
    }

    return () => {
      if (isOnWails) {
        window.runtime.EventsOff("discord_ready");
      }
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
