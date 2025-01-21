-- Create the tables that do not depend on other tables first
CREATE TABLE IF NOT EXISTS fertilizers (
  fertilizer_id uuid PRIMARY KEY,
  name varchar
);

CREATE TABLE IF NOT EXISTS material_types (
  material_type_id uuid PRIMARY KEY,
  name varchar
);

CREATE TABLE IF NOT EXISTS order_types (
  order_type_id uuid PRIMARY KEY,
  name varchar
);

CREATE TABLE IF NOT EXISTS seeds (
  seed_id uuid PRIMARY KEY,
  name varchar,
  type varchar
);

CREATE TABLE IF NOT EXISTS suppliers (
  supplier_id uuid PRIMARY KEY,
  name varchar,
  url varchar,
  suiss boolean
);

CREATE TABLE IF NOT EXISTS "orders" (
  order_id uuid PRIMARY KEY,
  total float,
  shipping float
);

-- Create the tables that reference the previously created tables
CREATE TABLE IF NOT EXISTS materials (
  material_id uuid PRIMARY KEY,
  name varchar,
  material_type_id uuid,
  main boolean,
  FOREIGN KEY (material_type_id) REFERENCES material_types (material_type_id)
);

CREATE TABLE IF NOT EXISTS order_items (
  order_item_id uuid PRIMARY KEY,
  order_id uuid,
  order_type_id uuid,
  name varchar,
  seed_id uuid,
  material_id uuid,
  supplier_id uuid,
  price float,
  FOREIGN KEY (order_id) REFERENCES "orders" (order_id),
  FOREIGN KEY (supplier_id) REFERENCES suppliers (supplier_id),
  FOREIGN KEY (material_id) REFERENCES materials (material_id),
  FOREIGN KEY (order_type_id) REFERENCES order_types (order_type_id)
);

CREATE TABLE IF NOT EXISTS seed_instructions (
  seed_instruction_id uuid PRIMARY KEY,
  seed_id uuid,
  seed_grams int,
  soaking_hours int,
  stacking_hours int,
  blackout_hours int,
  lights_hours int,
  yield_grams int,
  special_treatment varchar,
  FOREIGN KEY (seed_id) REFERENCES seeds (seed_id)
);

CREATE TABLE IF NOT EXISTS crops (
  crop_id uuid PRIMARY KEY,
  seed_id uuid,
  soaking_start timestamptz,
  stacking_start timestamptz,
  blackout_start timestamptz,
  lights_start timestamptz,
  harvest timestamptz,
  code varchar,
  yield_grams int null,
  FOREIGN KEY (seed_id) REFERENCES seeds (seed_id)
);

CREATE TABLE IF NOT EXISTS waterings (
  watering_id uuid PRIMARY KEY,
  crop_id uuid,
  quantity_ml int,
  fertilizer_ml int,
  fertilizer_id uuid,
  FOREIGN KEY (crop_id) REFERENCES crops (crop_id),
  FOREIGN KEY (fertilizer_id) REFERENCES fertilizers (fertilizer_id)
);

-- Create indexes
CREATE INDEX idx_crops_seed ON crops (seed_id);
CREATE INDEX idx_waterings_crop ON waterings (crop_id);
CREATE INDEX idx_order_items_order ON order_items (order_id);
CREATE INDEX idx_order_items_supplier ON order_items (supplier_id);
