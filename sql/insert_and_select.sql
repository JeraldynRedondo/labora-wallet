INSERT INTO public.items (customer_name , order_date , product , quantity , price)
VALUES ('John Doe', '01/02/2023', 'Televisor', 2, 150),
       ('Jane Smith', '07/02/2023', 'Celular', 3, 200),
       ('Michael Smith', '27/01/2023', 'Estufa', 1, 55),
       ('Harry Styles', '15/02/2023', 'Ventilador', 3, 70),
       ('Zayn Malik', '16/01/2023', 'Mouse', 3, 45),
       ('Bob Johnson', '05/03/2023', 'Lavadora', 1, 49);

SELECT * FROM public.items
WHERE (quantity>2) AND (price>50);