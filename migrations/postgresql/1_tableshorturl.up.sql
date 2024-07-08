CREATE TABLE IF NOT EXISTS public.shorturl
(
    key text NOT NULL
        CONSTRAINT shorturl_pk
            PRIMARY KEY,
    short text,
    original text CONSTRAINT shorturl_pk_2
        unique
);

COMMENT ON TABLE public.shorturl IS 'Сокрашённые ссылки URL';
COMMENT ON COLUMN public.shorturl.short IS 'Сокращённая';
COMMENT ON COLUMN public.shorturl.original IS 'Оригинальная ';