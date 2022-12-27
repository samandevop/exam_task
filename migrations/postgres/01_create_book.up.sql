CREATE TABLE categories (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL UNIQUE,
    parent_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT false,
    deleted_at TIMESTAMP
);

CREATE TABLE products (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR NOT NULL,
    price NUMERIC NOT NULL DEFAULT 0,
    category_id UUID NOT NULL REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT false,
    deleted_at TIMESTAMP
);

CREATE TABLE orders (
    id UUID PRIMARY KEY NOT NULL,
    description VARCHAR,
    product_id UUID NOT NULL REFERENCES products(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT false,
    deleted_at TIMESTAMP
)

-- SELECT
--     c1.id,
--     c1.name,
--     ARRAY_AGG(c2.id),
--     ARRAY_AGG(c2.name),
--     ARRAY_AGG(c2.parent_id)
-- FROM
--     categories as c1
-- JOIN categories as c2 ON c1.id = c2.parent_id
-- WHERE c1.is_deleted = false AND c2.is_deleted = false AND c1.id = 'c025ff25-8706-44d8-9fba-2284c67425a9'
-- GROUP BY c1.id, c1.name;


-- SELECT
--     orders.id,
--     orders.description,
--     products.id,
--     products.name,
--     categories.id,
--     categories.name,
--     categories.parent_id
-- FROM
--     orders
-- JOIN products ON orders.product_id = products.id
-- JOIN categories ON products.category_id = categories.id
-- WHERE orders.is_deleted = false AND products.is_deleted = false AND categories.is_deleted = false AND orders.id = $1;