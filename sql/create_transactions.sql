CREATE TABLE public.transactions
(
    id serial NOT NULL,
    wallet_id integer NOT NULL,
    transaction_type VARCHAR(255) NOT NULL,
    amount NUMERIC NOT NULL,
    date_transaction date NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (wallet_id) REFERENCES wallets (id)
);
