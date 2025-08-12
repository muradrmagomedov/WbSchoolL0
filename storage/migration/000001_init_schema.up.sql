CREATE TABLE "orders" (
  "id" serial PRIMARY KEY,
  "order_uid" varchar NOT NULL UNIQUE,
  "track_number" varchar NOT NULL UNIQUE,
  "entry" varchar NOT NULL,
  "locale" varchar NOT NULL,
  "internal_signature" varchar,
  "customer_id" varchar NOT NULL,
  "delivery_service" varchar NOT NULL,
  "shardkey" varchar NOT NULL,
  "sm_id" bigint NOT NULL,
  "date_created" timestamp DEFAULT CURRENT_TIMESTAMP,
  "oof_shard" varchar NOT NULL
);

CREATE TABLE "delivery" (
  "id" serial PRIMARY KEY,
  "order_uid" varchar NOT NULL,
  "name" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "zip" varchar NOT NULL,
  "city" varchar NOT NULL,
  "address" varchar NOT NULL,
  "region" varchar NOT NULL,
  "email" varchar NOT NULL,
  "date_created" timestamp DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (order_uid) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE "payment" (
  "id" serial PRIMARY KEY,
  "transaction" varchar UNIQUE,
  "request_id" varchar,
  "currency" varchar NOT NULL,
  "provider" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "payment_dt" bigint NOT NULL,
  "bank" varchar NOT NULL,
  "delivery_cost" bigint NOT NULL,
  "goods_total" bigint NOT NULL,
  "custom_fee" bigint NOT NULL,
  "date_created" timestamp DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (transaction) REFERENCES orders(order_uid) ON DELETE CASCADE
);

CREATE TABLE "items" (
  "id" serial PRIMARY KEY,
  "chrt_id" bigint NOT NULL UNIQUE,
  "track_number" varchar NOT NULL,
  "price" bigint NOT NULL,
  "rid" varchar NOT NULL,
  "name" varchar NOT NULL,
  "sale" bigint NOT NULL,
  "size" varchar NOT NULL,
  "total_price" bigint NOT NULL,
  "nm_id" bigint NOT NULL,
  "brand" varchar NOT NULL,
  "status" bigint NOT NULL,
  "date_created" timestamp DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_track_number FOREIGN KEY (track_number)  REFERENCES orders(track_number) ON DELETE CASCADE
);
