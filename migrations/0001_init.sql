-- +goose Up
CREATE TABLE users (
  id bigint generated always as identity primary key,
  telegram_user_id bigint not null unique,
  created_at timestamptz not null default now()
);

CREATE TABLE plants (
  id bigint generated always as identity primary key,
  user_id bigint not null references users(id) on delete cascade,
  name text not null,
  created_at timestamptz not null default now()
);

create index IDX_plants__user_id on plants(user_id);

CREATE TABLE care_events (
  id bigint generated always as identity primary key,
  plant_id bigint not null references plants(id) on delete cascade,
  event_type text not null,
  occurred_at timestamptz not null,
  created_at timestamptz not null default now()
);

create index IDX_care_events__last_event on care_events(plant_id, occurred_at desc);

CREATE TABLE reminders (
  id bigint generated always as identity primary key,
  plant_id bigint not null references plants(id) on delete cascade,
  reminder_type text not null,
  next_at timestamptz not null,
  created_at timestamptz not null default now()
);

create index IDX_reminders__last_reminders on reminders(next_at);

-- +goose Down
DROP TABLE reminders;
DROP TABLE care_events;
DROP TABLE plants;
DROP TABLE users;
