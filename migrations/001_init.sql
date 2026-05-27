CREATE TYPE order_status  AS ENUM ('PENDING','PAID','SHIPPED','COMPLETED','CANCELED');
CREATE TYPE draft_status  AS ENUM ('PENDING_PAYMENT','PAYMENT_FAILED','COMPLETED','EXPIRED');

CREATE TABLE users (
  id           VARCHAR PRIMARY KEY,
  email        VARCHAR UNIQUE NOT NULL,
  name         VARCHAR,
  image        VARCHAR,
  phone_number VARCHAR,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE addresses (
  id         VARCHAR PRIMARY KEY,
  user_id    VARCHAR NOT NULL REFERENCES users(id),
  full_name  VARCHAR NOT NULL,
  phone      VARCHAR NOT NULL,
  line1      VARCHAR NOT NULL,
  line2      VARCHAR,
  city       VARCHAR NOT NULL,
  state      VARCHAR,
  postal     VARCHAR NOT NULL,
  country    VARCHAR NOT NULL,
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE products (
  id          VARCHAR PRIMARY KEY,
  name        VARCHAR NOT NULL,
  description TEXT,
  images      TEXT[] NOT NULL DEFAULT '{}',
  category    VARCHAR,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE skus (
  id          VARCHAR PRIMARY KEY,
  product_id  VARCHAR NOT NULL REFERENCES products(id),
  sku_code    VARCHAR UNIQUE NOT NULL,
  price       INTEGER NOT NULL,          -- stored in smallest currency unit (cents/HKD cents)
  stock       INTEGER NOT NULL DEFAULT 0,
  attributes  JSONB NOT NULL,            -- e.g. {"size":"M","color":"red"}
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Replaces Prisma's implicit _ProductToUser join table
CREATE TABLE wishlists (
  user_id    VARCHAR NOT NULL REFERENCES users(id),
  product_id VARCHAR NOT NULL REFERENCES products(id),
  PRIMARY KEY (user_id, product_id)
);

CREATE TABLE draft_orders (
  id                  VARCHAR PRIMARY KEY,
  user_id             VARCHAR NOT NULL REFERENCES users(id),
  total               INTEGER NOT NULL,
  status              draft_status NOT NULL DEFAULT 'PENDING_PAYMENT',
  shipping_address_id VARCHAR NOT NULL REFERENCES addresses(id),
  payment_intent_id   VARCHAR,
  stripe_session_id   VARCHAR,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at          TIMESTAMPTZ
);

CREATE TABLE draft_order_items (
  id             VARCHAR PRIMARY KEY,
  draft_order_id VARCHAR NOT NULL REFERENCES draft_orders(id),
  sku_id         VARCHAR NOT NULL REFERENCES skus(id),
  quantity       INTEGER NOT NULL,
  price          INTEGER NOT NULL
);

CREATE TABLE orders (
  id                  VARCHAR PRIMARY KEY,
  user_id             VARCHAR NOT NULL REFERENCES users(id),
  total               INTEGER NOT NULL,
  status              order_status NOT NULL DEFAULT 'PENDING',
  shipping_address_id VARCHAR NOT NULL REFERENCES addresses(id),
  payment_intent_id   VARCHAR,
  stripe_session_id   VARCHAR,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  draft_order_id      VARCHAR UNIQUE NOT NULL REFERENCES draft_orders(id)
);

CREATE TABLE order_items (
  id        VARCHAR PRIMARY KEY,
  order_id  VARCHAR NOT NULL REFERENCES orders(id),
  sku_id    VARCHAR NOT NULL REFERENCES skus(id),
  quantity  INTEGER NOT NULL,
  price     INTEGER NOT NULL
);