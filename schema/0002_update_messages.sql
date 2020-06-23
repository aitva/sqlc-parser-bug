ALTER TABLE messages ALTER external_id TYPE UUID ARRAY USING ARRAY[external_id];
