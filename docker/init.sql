-- Создаём таблицу пользователей
CREATE TABLE IF NOT EXISTS public.users
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    username text COLLATE pg_catalog."default" NOT NULL,
    password_hash text COLLATE pg_catalog."default" NOT NULL,
    coins integer NOT NULL DEFAULT 1000,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_username_key UNIQUE (username)
);

-- Создаём таблицу инвентаря
CREATE TABLE IF NOT EXISTS public.inventories
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    type text COLLATE pg_catalog."default" NOT NULL,
    quantity integer NOT NULL DEFAULT 0,
    CONSTRAINT inventories_pkey PRIMARY KEY (id),
    CONSTRAINT inventories_user_id_type_key UNIQUE (user_id, type)
);

-- Создаём таблицу транзакций
CREATE TABLE IF NOT EXISTS public.transactions
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    from_user uuid NOT NULL,
    to_user uuid NOT NULL,
    amount integer NOT NULL,
    CONSTRAINT transactions_pkey PRIMARY KEY (id)
);

-- Добавляем внешние ключи
ALTER TABLE IF EXISTS public.inventories
    ADD CONSTRAINT inventories_user_id_fkey FOREIGN KEY (user_id)
    REFERENCES public.users (id) ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.transactions
    ADD CONSTRAINT transactions_from_user_fkey FOREIGN KEY (from_user)
    REFERENCES public.users (id) ON DELETE CASCADE;

ALTER TABLE IF EXISTS public.transactions
    ADD CONSTRAINT transactions_to_user_fkey FOREIGN KEY (to_user)
    REFERENCES public.users (id) ON DELETE CASCADE;
