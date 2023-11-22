CREATE TABLE message_queue (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);

CREATE ROLE keda WITH LOGIN SUPERUSER PASSWORD 'pass';

DO $$
BEGIN
    FOR i IN 1..10 LOOP
        EXECUTE 'CREATE ROLE producer' ||i|| ' WITH LOGIN SUPERUSER PASSWORD ''pass''';
    END LOOP;
END $$;

DO $$
BEGIN
    FOR i IN 1..100 LOOP
        EXECUTE 'CREATE ROLE consumer' ||i|| ' WITH LOGIN SUPERUSER PASSWORD ''pass''';
    END LOOP;
END $$;