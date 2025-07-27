-- name: CreateSubscription :one
INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
VALUES (
    $1,
    $2,
    $3,
    TO_DATE($4, 'MM-YYYY'), -- start_date
    CASE 
        WHEN $5 = '' THEN NULL
        ELSE TO_DATE($5, 'MM-YYYY')
    END  
)
RETURNING *;

-- name: GetSubscription :one
SELECT * FROM subscriptions
WHERE id = $1;


-- name: GetSubscriptionByServiceAndUser :one
SELECT * FROM subscriptions
WHERE service_name = $1
AND user_id = $2;

-- name: GetSubscriptionList :many
SELECT * FROM subscriptions;

-- name: DeleteSubscription :execrows
DELETE FROM subscriptions
WHERE id = $1;

-- name: UpdateSubscription :one
UPDATE subscriptions
SET
    service_name = $1,
    price = $2,
    end_date = CASE 
        WHEN $3 = '' THEN NULL
        ELSE TO_DATE($3, 'MM-YYYY')
    END,
    updated_at = NOW()
WHERE id = $4
RETURNING *;

-- name: CalculateTotalCost :one
SELECT 
    SUM(price) AS total_cost,
    COUNT(*) AS subscriptions_count
FROM 
    subscriptions
WHERE 
    start_date >= TO_DATE($1, 'MM-YYYY') AND
    start_date <= TO_DATE($2, 'MM-YYYY') AND
    user_id = $3::UUID AND
    service_name = $4 AND
    (end_date IS NULL OR end_date >= CURRENT_DATE);