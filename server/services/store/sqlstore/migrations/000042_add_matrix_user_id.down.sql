{{if .mysql}}
ALTER TABLE {{.prefix}}users DROP COLUMN matrix_user_id;
{{end}}

{{if .postgres}}
ALTER TABLE {{.prefix}}users DROP COLUMN matrix_user_id;
{{end}}

{{if .sqlite}}
-- SQLite does not support dropping columns directly. You would need to recreate the table without the column.
-- This is a placeholder for manual intervention if needed.
{{end}}
