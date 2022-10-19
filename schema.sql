-- schema.sql
-- Since we might run the import many times we'll drop if exists
DROP DATABASE IF EXISTS currencies;

CREATE DATABASE currencies;

-- Make sure we're using our `currencies` database
\c currencies;


-- We can create our requests table
CREATE TABLE IF NOT EXISTS requests_status (
  id SERIAL PRIMARY KEY,
  time_elapsed VARCHAR,
  url VARCHAR,
  status TEXT,
  requested_at TIMESTAMP
);

INSERT INTO requests_status
(
  time_elapsed,
  url,
  status,
  requested_at
)
VALUES 
('1s','test.com','success','2022-10-06T14:23:34');


-- We can create our currencies_values table
CREATE TABLE IF NOT EXISTS currencies_values (
  id SERIAL PRIMARY KEY,
  name VARCHAR,
  request_id INTEGER REFERENCES requests_status(id),
  value NUMERIC (10, 4),
  last_updated_at TIMESTAMP
);



INSERT INTO currencies_values
  (
    name,
    request_id,
    value,
    last_updated_at
  )
VALUES 
('MXN',1,'20.19','2022-10-06T14:23:34'),
('MXN',1,'20.19','2022-10-07T14:23:34'),
('MXN',1,'20.19','2022-10-08T14:23:34'),
('MXN',1,'20.19','2022-10-09T14:23:34'),
('MXN',1,'20.19','2022-10-013T14:23:34'),
('MXN',1,'20.19','2022-10-015T14:23:34'),
('MXN',1,'20.19','2022-10-018T14:23:34'),
('RUP',1,'34.19','2022-10-06T14:23:34'),
('RUP',1,'34.19','2022-10-07T14:23:34'),
('RUP',1,'34.19','2022-10-08T14:23:34'),
('RUP',1,'34.19','2022-10-09T14:23:34'),
('RUP',1,'34.19','2022-10-013T14:23:34'),
('RUP',1,'34.19','2022-10-015T14:23:34'),
('RUP',1,'34.19','2022-10-018T14:23:34'),
('AAA',1,'4.19','2022-10-06T14:23:34'),
('AAA',1,'4.19','2022-10-07T14:23:34'),
('AAA',1,'4.19','2022-10-08T14:23:34'),
('AAA',1,'4.19','2022-10-09T14:23:34'),
('AAA',1,'4.19','2022-10-013T14:23:34'),
('AAA',1,'4.19','2022-10-015T14:23:34'),
('AAA',1,'4.19','2022-10-018T14:23:34');



    -- # docker exec -it postgres_container /bin/sh
    -- # psql -U postgres currencies

    -- SELECT * FROM currencies_values;
    -- SELECT * FROM requests_status;


