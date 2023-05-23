CREATE DATABASE labora-project-2;

CREATE TABLE public.wallet
(
    id serial NOT NULL,
    dni_request INTEGER NOT NULL, 
    country_id VARCHAR(255) NOT NULL,
    order_request date NOT NULL,
    PRIMARY KEY (id)
);
