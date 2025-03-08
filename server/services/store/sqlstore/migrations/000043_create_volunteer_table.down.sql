{{if .mysql}}
ALTER TABLE {{.prefix}}users
DROP FOREIGN KEY volunteer_id,
DROP COLUMN volunteer_id;

DROP TABLE {{.prefix}}volunteer;
{{end}}

{{if .postgres}}
ALTER TABLE {{.prefix}}users
DROP CONSTRAINT IF EXISTS {{.prefix}}users_volunteer_id_fkey,
DROP COLUMN IF EXISTS volunteer_id;

DROP TABLE IF EXISTS {{.prefix}}volunteer;
{{end}}

{{if .sqlite}}
PRAGMA foreign_keys=off;

CREATE TABLE {{.prefix}}users_new AS SELECT * FROM {{.prefix}}users WHERE 1=0;
INSERT INTO {{.prefix}}users_new SELECT * FROM {{.prefix}}users;
DROP TABLE {{.prefix}}users;
ALTER TABLE {{.prefix}}users_new RENAME TO {{.prefix}}users;

PRAGMA foreign_keys=on;

DROP TABLE IF EXISTS {{.prefix}}volunteer;
{{end}}
