CREATE TABLE public.users (
	id serial8 NOT NULL,
	full_name varchar(255) NOT NULL,
	gender varchar(10) NOT NULL,
	dob timestamp NOT NULL,
	email varchar(255) NOT NULL,
	"password" varchar(255) NOT NULL,
	photo_url text NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT users_pk PRIMARY KEY (id),
	CONSTRAINT users_unique UNIQUE (email)
);

CREATE TABLE public.products (
	id serial4 NOT NULL,
	code varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT products_pk PRIMARY KEY (id),
	CONSTRAINT products_unique UNIQUE (code)
);

CREATE TABLE public.user_products (
	id serial4 NOT NULL,
	user_id int8 NOT NULL,
	product_id int4 NOT NULL,
	CONSTRAINT user_products_pk PRIMARY KEY (id),
	CONSTRAINT user_products_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id),
	CONSTRAINT user_products_products_fk FOREIGN KEY (product_id) REFERENCES public.products(id)
);

CREATE TABLE public.orders (
	id serial4 NOT NULL,
	user_id int8 NOT NULL,
	product_id serial4 NOT NULL,
	CONSTRAINT orders_pk PRIMARY KEY (id),
	CONSTRAINT orders_products_fk FOREIGN KEY (product_id) REFERENCES public.products(id),
	CONSTRAINT orders_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public."subscription" (
	id serial8 NOT NULL,
	user_id int8 NOT NULL,
	expired_at timestamp NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT subscription_pk PRIMARY KEY (id),
	CONSTRAINT subscription_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id)
);

CREATE TABLE public.user_subscribed_products (
	id serial8 NOT NULL,
	user_id int8 NOT NULL,
	expired_at timestamp NOT NULL,
	is_active boolean NOT NULL,
	product_code varchar(255) NOT NULL,
	product_name varchar(255) NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT subscribed_products_pk PRIMARY KEY (id),
	CONSTRAINT subscribed_products_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id)
);