-- +goose Up
-- +goose StatementBegin
CREATE TABLE stocks (
 id bigserial primary key,
 sku bigint not null unique,
 count bigint not null,
 reserved bigint not null,
 constraint count_nonnegative check (stocks.count >= 0),
 constraint reserved_nonnegative check (stocks.reserved >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks;
-- +goose StatementEnd
