CREATE TABLE IF NOT EXISTS {{.prefix}}roles (
	id VARCHAR(100),
	name VARCHAR(100),
  PRIMARY KEY (id)
) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};

CREATE TABLE IF NOT EXISTS {{.prefix}}permissions (
	id VARCHAR(100),
	name VARCHAR(100) UNIQUE,
	description VARCHAR(255),
  PRIMARY KEY (id)
) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};

CREATE TABLE IF NOT EXISTS {{.prefix}}role_permissions (
	role_id VARCHAR(100),
	permission_id VARCHAR(100),
  PRIMARY KEY (role_id, permission_id)
) {{if .mysql}}DEFAULT CHARACTER SET utf8mb4{{end}};

INSERT INTO {{.prefix}}roles (id, name) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', 'admin');
INSERT INTO {{.prefix}}roles (id, name) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', 'manager');
INSERT INTO {{.prefix}}roles (id, name) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', 'curator');
INSERT INTO {{.prefix}}roles (id, name) VALUES ('e0fd6b7b-4610-4367-9892-83e024f5fecf', 'volunteer');

INSERT INTO {{.prefix}}permissions (id, name) VALUES ('05e3394c-fe1c-4a0e-a57a-cadebe40270b', 'manage_system');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('f32a33ca-bc27-4070-8f1f-74ebc2dd4125', 'view_team');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('85f90379-c10b-4d19-a17f-a713710823e2', 'manage_team');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('74c9d0eb-c084-4fd5-ba04-34382f5ba30e', 'view_members');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('f218c55c-6cad-41d4-971d-d4ad566c25a7', 'manage_members');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('d35a4120-bc7f-4591-ad70-746cd302f25a', 'manage_board_type');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('7f56cefb-c8d7-45df-bc78-86b113e69f2e', 'delete_board');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('6939ce93-47b5-405e-99a7-4fef97f56202', 'view_board');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('755d4cb4-5db0-4f8e-9cf7-3f797540659b', 'manage_board_roles');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('d535f9a1-3b9c-4876-8105-2668189986b6', 'share_board');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('987ce8de-2974-4817-99da-92b3995e414d', 'manage_board_cards');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('799dc98f-96aa-4e85-ba83-f1b617f77844', 'manage_board_properties');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('71c7af54-4c18-4476-9466-7f0e09c0a733', 'comment_board_cards');
INSERT INTO {{.prefix}}permissions (id, name) VALUES ('5c62b3ea-20e7-4ea7-9cd2-23a009d4e777', 'delete_others_comments');

INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '05e3394c-fe1c-4a0e-a57a-cadebe40270b');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', 'f32a33ca-bc27-4070-8f1f-74ebc2dd4125');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '85f90379-c10b-4d19-a17f-a713710823e2');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '74c9d0eb-c084-4fd5-ba04-34382f5ba30e');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', 'f218c55c-6cad-41d4-971d-d4ad566c25a7');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', 'd35a4120-bc7f-4591-ad70-746cd302f25a');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '7f56cefb-c8d7-45df-bc78-86b113e69f2e');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '6939ce93-47b5-405e-99a7-4fef97f56202');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '755d4cb4-5db0-4f8e-9cf7-3f797540659b');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', 'd535f9a1-3b9c-4876-8105-2668189986b6');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '987ce8de-2974-4817-99da-92b3995e414d');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '799dc98f-96aa-4e85-ba83-f1b617f77844');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '71c7af54-4c18-4476-9466-7f0e09c0a733');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('ea0c35bb-4a38-4a21-b331-3a5e9e50dc1d', '5c62b3ea-20e7-4ea7-9cd2-23a009d4e777');

INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', 'f32a33ca-bc27-4070-8f1f-74ebc2dd4125');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', '85f90379-c10b-4d19-a17f-a713710823e2');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', '74c9d0eb-c084-4fd5-ba04-34382f5ba30e');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', 'd35a4120-bc7f-4591-ad70-746cd302f25a');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', '7f56cefb-c8d7-45df-bc78-86b113e69f2e');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', '6939ce93-47b5-405e-99a7-4fef97f56202');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', '755d4cb4-5db0-4f8e-9cf7-3f797540659b');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', 'd535f9a1-3b9c-4876-8105-2668189986b6');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', '987ce8de-2974-4817-99da-92b3995e414d');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', '799dc98f-96aa-4e85-ba83-f1b617f77844');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', '71c7af54-4c18-4476-9466-7f0e09c0a733');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('16920b18-e451-4160-a037-d5f683406e6c', '5c62b3ea-20e7-4ea7-9cd2-23a009d4e777');

INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', 'f32a33ca-bc27-4070-8f1f-74ebc2dd4125');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', '85f90379-c10b-4d19-a17f-a713710823e2');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', '74c9d0eb-c084-4fd5-ba04-34382f5ba30e');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', '6939ce93-47b5-405e-99a7-4fef97f56202');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', '755d4cb4-5db0-4f8e-9cf7-3f797540659b');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', 'd535f9a1-3b9c-4876-8105-2668189986b6');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', '987ce8de-2974-4817-99da-92b3995e414d');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', '799dc98f-96aa-4e85-ba83-f1b617f77844');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('b0e84d16-4b16-42c6-acae-002f90f3f70f', '71c7af54-4c18-4476-9466-7f0e09c0a733');

INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('e0fd6b7b-4610-4367-9892-83e024f5fecf', '6939ce93-47b5-405e-99a7-4fef97f56202');
INSERT INTO {{.prefix}}role_permissions (role_id, permission_id) VALUES ('e0fd6b7b-4610-4367-9892-83e024f5fecf', 'f32a33ca-bc27-4070-8f1f-74ebc2dd4125');

ALTER TABLE {{.prefix}}users ADD COLUMN role_id VARCHAR(100) NOT NULL DEFAULT 'e0fd6b7b-4610-4367-9892-83e024f5fecf';
