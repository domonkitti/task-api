-- +goose Up
ALTER TABLE items ADD price real NOT NULL;

ALTER TABLE items ADD status text NOT NULL;

ALTER TABLE items ADD owner text NOT NULL;

-- +goose StatementBegin
SELECT
    'up SQL query';

-- +goose StatementEnd
-- +goose Down
ALTER TABLE items
DROP COLUMN price;

ALTER TABLE items
DROP COLUMN status;

ALTER TABLE items
DROP COLUMN owner;


-- +goose StatementBegin
SELECT
    'down SQL query';

-- +goose StatementEnd