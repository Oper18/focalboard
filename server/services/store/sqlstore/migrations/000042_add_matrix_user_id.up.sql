{{if .mysql}}
ALTER TABLE {{.prefix}}users ADD COLUMN matrix_user_id VARCHAR(255) NULL;
{{end}}

{{if .postgres}}
ALTER TABLE {{.prefix}}users ADD COLUMN matrix_user_id VARCHAR(255) NULL;
{{end}}

{{if .sqlite}}
ALTER TABLE {{.prefix}}users ADD COLUMN matrix_user_id TEXT NULL;
{{end}}
