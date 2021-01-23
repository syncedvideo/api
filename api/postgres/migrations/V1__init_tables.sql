CREATE TABLE IF NOT EXISTS sv_user (
    id UUID PRIMARY KEY,
    name VARCHAR,
    color VARCHAR,
    is_admin BOOLEAN DEFAULT false,
    ip_address INET,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS sv_room (
    id UUID PRIMARY KEY,
    owner_user_id UUID,
    name VARCHAR,
    CONSTRAINT sv_room_owner_user_id_fkey FOREIGN KEY(owner_user_id) REFERENCES sv_user(id)
);

CREATE TABLE IF NOT EXISTS sv_room_user (
    room_id UUID NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT sv_room_user_room_id_fkey FOREIGN KEY(room_id) REFERENCES sv_room(id) ON DELETE CASCADE,
    CONSTRAINT sv_room_user_user_id_fkey FOREIGN KEY(user_id) REFERENCES sv_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sv_queue_item (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT sv_queue_item_room_id_fkey FOREIGN KEY(room_id) REFERENCES sv_room(id) ON DELETE CASCADE,
    CONSTRAINT sv_queue_item_user_id_fkey FOREIGN KEY(user_id) REFERENCES sv_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sv_queue_item_vote (
    id UUID PRIMARY KEY,
    queue_item_id UUID NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT sv_queue_item_vote_queue_item_id_fkey FOREIGN KEY(queue_item_id) REFERENCES sv_queue_item(id) ON DELETE CASCADE,
    CONSTRAINT sv_queue_item_vote_user_id_fkey FOREIGN KEY(user_id) REFERENCES sv_user(id) ON DELETE CASCADE
);
