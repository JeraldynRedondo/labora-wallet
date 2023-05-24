CREATE TABLE public.logs
(
    id serial NOT NULL,
	dni_request VARCHAR(255) NOT NULL,
    status_request VARCHAR(255) NOT NULL,
    order_request date NOT NULL, 
    PRIMARY KEY (id)
);