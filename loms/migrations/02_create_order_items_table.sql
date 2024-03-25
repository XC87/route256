-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_items (
 id bigserial primary key,
 order_id bigint not null references orders(id) ON DELETE CASCADE,
 sku bigint not null,
 count bigint not null
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items;
-- +goose StatementEnd
