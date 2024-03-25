-- name: CreateOrder :one
insert into orders (created_at, updated_at, user_id, status) values ($1,$2, $3, $4) returning id;

-- name: CreateOrderItems :exec
insert into order_items (order_id, sku, count) values ($1, $2, $3);

-- name: UpdateOrder :exec
update orders set status = $1, user_id = $2, created_at = $3, updated_at = $4 where id = $5;

-- name: UpdateOrderStatus :exec
update orders set status = $1 where id = $2;

-- name: GetOrderById :many
select sqlc.embed(o), sqlc.embed(oi) from orders o
    left join order_items oi on o.id = oi.order_id
where o.id = $1;

-- name: GetBySku :one
select count, reserved from stocks where sku = $1;

-- name: UpdateReserveBySku :execresult
update stocks set reserved = reserved + $1 where sku = $2 returning id;

-- name: UpdateCountBySku :execresult
update stocks set count = count + $1 where sku = $2 returning id;
