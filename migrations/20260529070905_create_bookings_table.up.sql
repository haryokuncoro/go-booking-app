CREATE TABLE bookings (

    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,

    room_id BIGINT NOT NULL,

    booking_date DATE NOT NULL,

    status VARCHAR(50) NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),

    updated_at TIMESTAMP DEFAULT NOW(),

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id),

    CONSTRAINT idx_room_date
        UNIQUE(room_id, booking_date)
);