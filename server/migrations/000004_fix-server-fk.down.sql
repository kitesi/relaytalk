-- Drop correct FK (undo)
ALTER TABLE channels DROP CONSTRAINT IF EXISTS channels_server_id_fkey;

-- (Optional) Re-add incorrect FK if needed
ALTER TABLE channels
ADD CONSTRAINT channels_server_id_fkey
FOREIGN KEY (server_id) REFERENCES server(id);
