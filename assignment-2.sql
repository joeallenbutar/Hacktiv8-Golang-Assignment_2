CREATE DATABASE assignment-2;

CREATE TABLE orders(
    order_id serial,
    customer_name character varying(100) NOT NULL,
    ordered_at timestamp without time zone,
	PRIMARY KEY(order_id)
);

CREATE TABLE items(
    item_id serial,
    item_code character varying(50) NOT NULL,
    description character varying(255),
    quantity integer,
    order_id integer,
	PRIMARY KEY(item_id)
);