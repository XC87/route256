-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
 id bigserial primary key,
 created_at timestamp not null,
 updated_at timestamp null,
 user_id bigint not null,
 status int not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
