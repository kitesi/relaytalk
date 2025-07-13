ALTER TABLE channels DROP CONSTRAINT IF EXISTS channels_server_id_fkey;

ALTER TABLE channels
ADD CONSTRAINT channels_server_id_fkey
FOREIGN KEY (server_id) REFERENCES servers(id);
