CREATE TABLE "order" (
    "id"    bigserial PRIMARY KEY,
    "weight"    real NOT NULL,
    "destination"    POINT NOT NULL,
    "observations"    text,
    "finished"    boolean NOT NULL DEFAULT FALSE
);

CREATE TABLE "vehicle" (
    "id"    bigserial PRIMARY KEY,
    "max_weight_capacity"    real NOT NULL,
    "number_plate"    text NOT NULL UNIQUE
);

CREATE TABLE "assigned_order" (
    "id"    bigserial PRIMARY KEY,
    "assigned_at"    timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
    "vehicle_id"    bigint NOT NULL,
    "order_id"    bigint NOT NULL
);

ALTER TABLE "assigned_order" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicle" ("id");
ALTER TABLE "assigned_order" ADD FOREIGN KEY ("order_id") REFERENCES "order" ("id");

CREATE INDEX ON "assigned_order" ("vehicle_id");
