CREATE TABLE IF NOT EXISTS "StockMovement" (
    "id" TEXT PRIMARY KEY,
    "partName" TEXT NOT NULL,
    "partModel" TEXT NOT NULL,
    "delta" INTEGER NOT NULL,
    "reason" TEXT NOT NULL,
    "referenceId" TEXT,
    "createdAt" TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_stock_movement_part ON "StockMovement"("partName", "partModel");
CREATE INDEX IF NOT EXISTS idx_stock_movement_reference ON "StockMovement"("referenceId");
