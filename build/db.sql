-- User table
CREATE TABLE public.user (
    id BIGINT  NOT NULL PRIMARY KEY,
    name VARCHAR(30) NOT NULL
);

-- Coupon activaty
CREATE TYPE coupon_active_state AS ENUM ('NOT_OPEN', 'OPENING', 'CLOSED');

CREATE TABLE coupon_active (
    id BIGINT NOT NULL PRIMARY KEY,
    date Date NOT NULL,
    begin_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    end_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    state coupon_active_state NOT NULL
);

-- history about user reserve coupon
CREATE TABLE user_reserved_coupon_active_history (
    user_id BIGINT NOT NULL,
    active_id BIGINT NOT NULL,
    serial_num BIGINT NOT NULL,
    reserved_at BIGINT,
    PRIMARY KEY (user_id, active_id),
    FOREIGN KEY (user_id) REFERENCES public.user (user_id),
    FOREIGN KEY (active_id) REFERENCES public.coupon_active (id)
);

CREATE TYPE coupon_state AS ENUM ('UNRESERVED', 'RESERVED', 'USED');

-- Coupon table
CREATE TABLE coupon (
    id BIGINT PRIMARY KEY,
    active_id BIGINT,
    coupon_code VARCHAR(50) NOT NULL,
    state coupon_state NOT NULL,
    buyer BIGINT,
    buy_time BIGINT,
    created_at BIGINT NOT NULL,
    FOREIGN KEY (active_id) REFERENCES public.coupon_active (id),
    FOREIGN KEY (buyer) REFERENCES public.user (user_id)
);
