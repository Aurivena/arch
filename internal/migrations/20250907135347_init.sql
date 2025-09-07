-- +goose Up
-- +goose StatementBegin
CREATE TABLE "Test"(
    id serial
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
