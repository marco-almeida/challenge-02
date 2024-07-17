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
    "number_plate"    text NOT NULL UNIQUE,
    "current_weight" real NOT NULL DEFAULT 0
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


-- init vehicles
INSERT INTO public.vehicle (max_weight_capacity,number_plate) VALUES
	 (40000,'11-ad-92'),
	 (40000,'74-qg-23'),
	 (36000,'10-ae-12'),
	 (32000,'55-bk-30');

-- init orders
INSERT INTO public.order (weight,destination,observations) VALUES
	 (35000.23,'(47.2342,5.45)','to be delivered in gate A'),
	 (1000,'(38.71,-9.14)','Z Lisbon stuff'),
     (10000, '(40.97134, -5.66337)', 'A'),
     (10000, '(39.47590, -0.37696)', 'B'),
     (10000, '(41.15936, -8.62907)', 'C'),
     (10000, '(40.42382, -3.70254)', 'D');


