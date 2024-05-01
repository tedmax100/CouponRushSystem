INSERT INTO coupon_active (date, begin_time, end_time, purchase_begin_time, purchase_end_time, state)
VALUES (
    CURRENT_DATE,
    NOW() + INTERVAL '1 minute',
    NOW() + INTERVAL '6 minutes',
    NOW() + INTERVAL '7 minutes',
    NOW() + INTERVAL '8 minutes',
    'OPENING'
)
RETURNING *;