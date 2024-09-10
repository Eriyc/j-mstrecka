import { ProductResponse } from "@/types";
import { Card, CardHeader, CardTitle, CardContent } from "./ui/card";
import { Badge } from "./ui/badge";

type Props = {
  product: ProductResponse | null;
};

export const ProductInfo = ({ product }: Props) => {
  if (!product) {
    return null;
  }

  return (
    <Card className="w-full max-w-lg flex-1">
      <CardHeader>
        <CardTitle className="flex justify-between items-center">
          <span>{product.product.name}</span>
          <Badge variant="outline">Product</Badge>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <span className="font-semibold">Internal Price:</span>
            <span className="text-xl">
              {product.price.internal_price.toFixed(2)} kr
            </span>
          </div>
          <div className="flex justify-between items-center">
            <span className="font-semibold">Stock:</span>
            <span>{product.product.total_stock} st</span>
          </div>
          <div className="text-sm text-muted-foreground">
            <div className="text-slate-400">
              External Price: <span className="text-white">{product.price.external_price.toFixed(2)} kr</span>
            </div>
            <div className="text-slate-400">
              Purchase Price: <span className="text-white">{product.price.purchase_price.toFixed(2)} kr</span>
            </div>
            <div className="text-slate-400">
              Valid from:{" "}
              <span className="text-white">
                {new Date(product.price.start_date).toLocaleDateString()}
              </span>{" "}
              to{" "}
              <span className="text-white">
                {new Date(product.price.end_date).toLocaleDateString()}
              </span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
