CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    login         VARCHAR(256) NOT NULL,
    name          VARCHAR(256) NOT NULL,
    surname       VARCHAR(256) NOT NULL,
    password_hash TEXT         NOT NULL,
    created_at    timestamptz           DEFAULT now(),
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE
);

CREATE TABLE refresh_tokens
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES Users (id) ON DELETE CASCADE,
    token_hash TEXT        NOT NULL,
    expired_at timestamptz NOT NULL,
    revoked    BOOLEAN     NOT NULL DEFAULT FALSE
);

