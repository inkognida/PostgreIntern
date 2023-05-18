create table if not exists system_file_events (
                                             id SERIAL PRIMARY KEY,
                                             event_type TEXT NOT NULL,
                                             path TEXT NOT NULL,
                                             file_name TEXT NOT NULL,
                                             created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

create table if not exists commands_events (
                                               id SERIAL PRIMARY KEY,
                                               command TEXT NOT NULL,
                                               args TEXT NOT NULL,
                                               executed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                               system_event_id INTEGER REFERENCES commands_events(id)

);


