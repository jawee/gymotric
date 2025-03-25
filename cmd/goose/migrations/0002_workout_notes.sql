-- +goose Up
-- +goose StatementBegin
ALTER TABLE workouts
ADD COLUMN note text null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workouts
DROP COLUMN note;
-- +goose StatementEnd
