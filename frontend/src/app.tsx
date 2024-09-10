import { useCallback, useState } from "react";
import DiscordStatus from "./components/discord-status";
import { DrinkChart } from "./components/comptetitive-chart";
import { Leaderboard } from "./components/leaderboard";
import { ScanUpc, Strecka } from "@/../wailsjs/go/desktop/App";
import { ProductResponse, UserResponse } from "./types";
import { UserInfo } from "./components/user-info";
import { ProductInfo } from "./components/product-info";
import { useToast } from "./lib/toast-context";
import { useTransactions } from "./lib/transaction-context";
import { useKeyboardListener } from "./hooks/use-keyboard-listener";

function App() {
  const { transactions } = useTransactions();
  const { addToast } = useToast();

  const [error, setError] = useState("");

  const [user, setUser] = useState<UserResponse | null>(null);
  const [product, setProduct] = useState<ProductResponse | null>(null);

  const onScan = useCallback(
    async (upc: string) => {
      const result = await ScanUpc(upc);
      if (result === null) {
        setError("Scan failed");
        return;
      }
      if (result.type === "user") {
        setUser(result);
      } else if (result.type === "product") {
        setProduct(result);
        const p = result as ProductResponse;
        if (user !== null) {
          const amount = 1;
          const {
            user: newUser,
            product: newProduct,
            balance,
            error,
          } = await Strecka(p.product.id, user.user.id, amount).catch((err) =>
            console.error("error strecka:", err)
          );

          console.log(error, newUser);

          setUser({
            balance,
            user: newUser,
            type: "user",
          });

          setProduct({
            ...product!,
            product: {
              ...product!.product,
              total_stock: newProduct.total_stock,
            },
            type: "product",
          });

          addToast(`Streckade ${p.product.name} för ${user.user.name}`);
        }
      }
    },
    [user, product]
  );

  useKeyboardListener(onScan);

  return (
    <div className="flex-1 flex flex-col dark:bg-indigo-950 dark:text-white">
      <header className="p-4 flex flex-row items-center justify-between">
        <div>
          <div className="bg-rainbow-gradient text-transparent bg-clip-text font-bold text-4xl">
            Strecka
          </div>
          <p className="text-red-500">{error}</p>
        </div>
        <DiscordStatus />
      </header>
      <section className="flex flex-1 gap-4 px-2 pb-8 items-baseline">
        <div className="gap-4 flex-col flex">
          <UserInfo user={user} onRemove={() => setUser(null)} />
          <ProductInfo product={product} />
        </div>
        <div className="flex-[3] grid grid-cols-2 gap-4">
          <DrinkChart {...transactions} />
          <Leaderboard />
        </div>
      </section>
    </div>
  );
}

export default App;
