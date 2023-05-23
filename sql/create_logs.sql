CREATE TABLE public.logs
(
    id serial NOT NULL,
    order_request date NOT NULL,
    status_request VARCHAR(255) NOT NULL,
    dni_request INTEGER NOT NULL, 
    PRIMARY KEY (id)
);