CREATE TABLE public.logs
(
    id serial NOT NULL,
	dni_request VARCHAR(255) NOT NULL,
    country_id VARCHAR(255) NOT NULL,
    status_request VARCHAR(255) NOT NULL,
    date_request date NOT NULL, 
    request_type VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);