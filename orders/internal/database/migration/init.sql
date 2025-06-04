-- Таблица заказов (Orders)
CREATE TABLE Orders (
    UUID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    CustomerID UUID NOT NULL,
    OrderDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    TotalAmount NUMERIC(10, 2) NOT NULL,
    Status VARCHAR(20) NOT NULL DEFAULT 'Created',
    FOREIGN KEY (CustomerID) REFERENCES Customers(UUID),
    CHECK (Status IN ('Created', 'ReadyForPickup', 'Received'))
);

-- Таблица позиций заказов (OrderItems)
CREATE TABLE OrderItems (
    UUID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    OrderID UUID NOT NULL,
    GoodUUID UUID NOT NULL,
    Quantity INT NOT NULL,
    FOREIGN KEY (OrderID) REFERENCES Orders(UUID),
    FOREIGN KEY (GoodUUID) REFERENCES Goods(UUID)
);

CREATE TABLE Bag(
    UUID UUID PRIMARY KEY DEFAULT uuid_generate_v4(), 
    CustomerID UUID NOT NULL,
    Goods []UUID
)
