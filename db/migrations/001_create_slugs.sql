CREATE TABLE IF NOT EXISTS  slugs(
    id BIGINT PRIMARY KEY,
    url TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), 
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT url_length_check CHECK ( LENGTH(url) <= 2000 ),
    CONSTRAINT slug_length_check CHECK ( LENGTH(slug) <= 10 ) 
);
COMMENT ON TABLE slugs IS 'Сокращенные ссылки';

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_slugs_updated_at BEFORE UPDATE ON slugs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();