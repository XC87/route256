-- +goose Up
-- +goose StatementBegin
CREATE TABLE orderStatus
(
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
INSERT INTO orderStatus (name)
VALUES ('new'),
       ('awaiting payment'),
       ('failed'),
       ('paid'),
       ('cancelled');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orderStatus;
-- +goose StatementEnd
