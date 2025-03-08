{{if .mysql}}
CREATE TABLE {{.prefix}}volunteers (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(512) NOT NULL,
    contact VARCHAR(512) NOT NULL,
    files JSON NOT NULL,
    description TEXT NULL,
    create_at BIGINT,
    update_at BIGINT,
    delete_at BIGINT
);

ALTER TABLE {{.prefix}}users
ADD COLUMN volunteer_id INT,
ADD FOREIGN KEY (volunteer_id) REFERENCES {{.prefix}}volunteers(id);
{{end}}

{{if .postgres}}
CREATE TABLE {{.prefix}}volunteers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(512) NOT NULL,
    contact VARCHAR(512) NOT NULL,
    files TEXT[] NOT NULL,
    description TEXT NULL,
    create_at BIGINT,
    update_at BIGINT,
    delete_at BIGINT
);

ALTER TABLE {{.prefix}}users
ADD COLUMN volunteer_id INT,
ADD FOREIGN KEY (volunteer_id) REFERENCES {{.prefix}}volunteers(id);
{{end}}

{{if .sqlite}}
CREATE TABLE {{.prefix}}volunteers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    contact TEXT NOT NULL,
    files TEXT NOT NULL,
    description TEXT NULL,
    create_at BIGINT,
    update_at BIGINT,
    delete_at BIGINT
);

ALTER TABLE {{.prefix}}users
ADD COLUMN volunteer_id INTEGER,
ADD FOREIGN KEY (volunteer_id) REFERENCES {{.prefix}}volunteers(id);
{{end}}
