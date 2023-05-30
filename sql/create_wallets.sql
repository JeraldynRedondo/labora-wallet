CREATE DATABASE labora-project-2;

CREATE TABLE public.wallets
(
    id serial NOT NULL,
    dni_request VARCHAR(255) NOT NULL UNIQUE, 
    country_id VARCHAR(255) NOT NULL,
    created_date date NOT NULL,
    balance NUMERIC NOT NULL CHECK(balance >= 0),
    PRIMARY KEY (id)
);

