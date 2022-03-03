DROP TABLE IF EXISTS pings;

CREATE TABLE pings (
  id            SERIAL PRIMARY KEY,
  created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
  status        INT NOT NULL,
  ping_time     FLOAT NOT NULL,
);

INSERT INTO pings (status) VALUES (0, 65.70),(1, 65.70),(1, 53.70),(1, 21.70),(1, 85.74),(1, 165.70),(0,365.70),(0, 565.76),(0, 465.02);
