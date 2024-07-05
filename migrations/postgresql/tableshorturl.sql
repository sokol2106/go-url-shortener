CREATE TABLE public.shorturl
(
    key text NOT NULL
        CONSTRAINT shorturl_pk
            PRIMARY KEY,
    short text,
    original text
);

COMMENT ON TABLE public.shorturl IS 'Сокрашённые ссылки URL';
COMMENT ON COLUMN public.shorturl.short IS 'Сокращённая';
COMMENT ON COLUMN public.shorturl.original IS 'Оригинальная ';