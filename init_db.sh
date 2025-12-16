#!/usr/bin/env bash
set -euo pipefail

# Project DB initializer for secure-notes
# - Creates the database if missing
# - Creates required tables
# - Creates required global indexes (parent table indexes are kept minimal; child
#   indexes are created by the app's partition maintenance if you enabled it)

# ------------------------------
# Configuration
# ------------------------------
: "${DB_HOST:=localhost}"
: "${DB_PORT:=5432}"
: "${DB_USER:=postgres}"
: "${DB_PASSWORD:=postgres}"
: "${DB_NAME:=notepad}"

export PGPASSWORD="${DB_PASSWORD}"

# Resolve docker compose command if needed
get_compose_cmd() {
  if command -v docker >/dev/null 2>&1 && docker compose version >/dev/null 2>&1; then
    echo "docker compose"
  elif command -v docker-compose >/dev/null 2>&1; then
    echo "docker-compose"
  else
    echo ""
  fi
}

have_psql() { command -v psql >/dev/null 2>&1; }

run_psql_cmd() {
  local db="$1"; shift
  local sql="$1"; shift || true
  if have_psql; then
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "$db" -v ON_ERROR_STOP=1 -tAc "$sql"
  else
    local dcmd
    dcmd=$(get_compose_cmd)
    if [[ -z "$dcmd" ]]; then
      echo "Error: psql is not installed and docker compose is unavailable." >&2
      return 127
    fi
    # Use psql inside the postgres container
    ${dcmd} exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
      psql -h localhost -p 5432 -U "${DB_USER}" -d "$db" -v ON_ERROR_STOP=1 -tAc "$sql"
  fi
}

echo "Checking database '${DB_NAME}' on ${DB_HOST}:${DB_PORT} as ${DB_USER}..."
DB_EXISTS=$(run_psql_cmd postgres "SELECT 1 FROM pg_database WHERE datname='${DB_NAME}';" || echo "")
if [[ "$DB_EXISTS" != "1" ]]; then
  echo "Creating database '${DB_NAME}'..."
  if have_psql; then
    psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -v ON_ERROR_STOP=1 -c "CREATE DATABASE ${DB_NAME};"
  else
    dcmd=$(get_compose_cmd)
    if [[ -z "$dcmd" ]]; then
      echo "Error: cannot create DB — psql missing and docker compose not found." >&2
      exit 1
    fi
    ${dcmd} exec -T -e PGPASSWORD="${DB_PASSWORD}" postgres \
      psql -h localhost -p 5432 -U "${DB_USER}" -v ON_ERROR_STOP=1 -c "CREATE DATABASE ${DB_NAME};"
  fi
else
  echo "Database '${DB_NAME}' already exists."
fi

#
#echo "Creating tables and indexes in '${DB_NAME}'..."
#if have_psql; then
#psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -v ON_ERROR_STOP=1 <<'SQL'
#-- Time-based daily partitions (RANGE) for all except prefixes
#
#-- user_requests partitioned by create_at (daily). No global PK; per-partition PK enforced.
#CREATE TABLE IF NOT EXISTS user_requests (
#  request_id      uuid NOT NULL,
#  client_id       TEXT NOT NULL,
#  total_messages  SMALLINT NOT NULL,
#  create_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
#) PARTITION BY RANGE (create_at);
#
#-- request_messages partitioned by create_at (daily). No FK, no global PK.
#CREATE TABLE IF NOT EXISTS request_messages (
#  group_id       uuid NOT NULL,
#  request_id     uuid NOT NULL,
#  client_id      TEXT NOT NULL,
#  track_id       TEXT,
#  sender         TEXT    NOT NULL,
#  receiver       TEXT    NOT NULL,
#  message        TEXT    NOT NULL,
#  overall_status TEXT    DEFAULT 'PENDING',
#  delivered      SMALLINT  DEFAULT 0,
#  failed         SMALLINT  DEFAULT 0,
#  total_part     SMALLINT  DEFAULT 0,
#  create_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
#  complete_at    TIMESTAMPTZ,
#    PRIMARY KEY (group_id, create_at)
#) PARTITION BY RANGE (create_at);
#
#-- submit_messages partitioned by create_at (daily). No FK to request_messages, per-partition PK on id.
#CREATE TABLE IF NOT EXISTS submit_messages (
#  id            uuid NOT NULL,
#  group_id      uuid NOT NULL,
#  sender        TEXT NOT NULL,
#  receiver      TEXT NOT NULL,
#  message_seq   SMALLINT NOT NULL DEFAULT 0,
#  total_part    SMALLINT NOT NULL DEFAULT 0,
#  data_coding   SMALLINT NOT NULL DEFAULT 0,
#  message_ref   SMALLINT NOT NULL DEFAULT 0,
#  message_part  BYTEA NOT NULL,
#  -- Send status fields
#  submit_send_status TEXT,
#  submit_resp_status TEXT,
#  deliver_status   TEXT,
#  -- Response status fields
#  server_id     TEXT,
#  smsc_message_id TEXT,
#  resp_error    TEXT,
#  dlr_error     TEXT,
#  -- Delivery status fields
#  create_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
#  sent_at       TIMESTAMPTZ,
#  resp_at       TIMESTAMPTZ,
#  deliver_resp_at  TIMESTAMPTZ,
#    PRIMARY KEY (id, create_at)
#) PARTITION BY RANGE (create_at);
#
#-- deliver_messages partitioned by create_at (daily). Per-partition PK on id.
#CREATE TABLE IF NOT EXISTS deliver_messages (
#  id           uuid NOT NULL,
#  client_id    TEXT NOT NULL,
#  server_id    TEXT NOT NULL,
#  receiver     TEXT NOT NULL,
#  total_part   SMALLINT NOT NULL DEFAULT 0,
#  message_ref  SMALLINT NOT NULL DEFAULT 0,
#  sender       TEXT NOT NULL,
#  data_coding  SMALLINT NOT NULL DEFAULT 0,
#  message_seq  SMALLINT NOT NULL DEFAULT 0,
#  message      TEXT NOT NULL,
#  create_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
#    PRIMARY KEY (id, create_at)
#) PARTITION BY RANGE (create_at);
#
#-- complete_deliver_messages partitioned by received_at (daily)
#CREATE TABLE IF NOT EXISTS complete_deliver_messages (
#  id           INTEGER NOT NULL GENERATED ALWAYS AS IDENTITY,
#  client_id    TEXT NOT NULL,
#  receiver     TEXT NOT NULL,
#  server_id    TEXT NOT NULL,
#  sender       TEXT NOT NULL,
#  data_coding  SMALLINT NOT NULL DEFAULT 0,
#  total_part   SMALLINT NOT NULL DEFAULT 0,
#  message      TEXT NOT NULL,
#  message_ref  SMALLINT NOT NULL DEFAULT 0,
#  received_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
#    PRIMARY KEY (id, received_at)
#) PARTITION BY RANGE (received_at);
#
#-- prefixes remains unpartitioned
#CREATE TABLE IF NOT EXISTS prefixes (
#  id         INTEGER GENERATED ALWAYS AS IDENTITY,
#  client_id  TEXT NOT NULL,
#  prefix     TEXT NOT NULL UNIQUE,
#  is_prefix  bool NOT NULL,
#  create_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
#);
#
#CREATE OR REPLACE FUNCTION create_daily_partitions_from_today(
#    parent regclass,
#    days_ahead int DEFAULT 3  -- how many total days (starting from today)
#)
#    RETURNS void
#    LANGUAGE plpgsql AS $$
#DECLARE
#    base date := CURRENT_DATE;  -- today's date in server timezone
#    i int;
#    child_name text;
#BEGIN
#    IF days_ahead < 1 THEN
#        RAISE NOTICE 'Nothing to do: days_ahead=%', days_ahead;
#        RETURN;
#    END IF;
#
#    -- Create partitions for: [today .. today + (days_ahead - 1)]
#    FOR i IN 0..(days_ahead - 1) LOOP
#            child_name := format('%s_%s', parent::text, to_char(base + i, 'YYYY_MM_DD'));
#
#            EXECUTE format(
#                    'CREATE TABLE IF NOT EXISTS %I PARTITION OF %s
#                       FOR VALUES FROM (%L) TO (%L);',
#                    child_name,
#                    parent::text,
#                    (base + i)::text,        -- lower bound (start of this day)
#                    (base + i + 1)::text     -- upper bound (start of next day)
#                    );
#        END LOOP;
#
#    RAISE NOTICE 'Created daily partitions for %: [% .. %]', parent::text, base, base + (days_ahead - 1);
#END$$;
#
#
#------ past days
#
#CREATE OR REPLACE FUNCTION create_past_daily_partitions(
#    parent regclass,
#    days_back int DEFAULT 35  -- how many days to create, starting from yesterday
#)
#    RETURNS void
#    LANGUAGE plpgsql AS $$
#DECLARE
#    base date := CURRENT_DATE;  -- today's date in server timezone
#    i int;
#    child_name text;
#BEGIN
#    IF days_back < 1 THEN
#        RAISE NOTICE 'Nothing to do: days_back=%', days_back;
#        RETURN;
#    END IF;
#
#    -- Create partitions for: [base-1 .. base-days_back], i.e., yesterday down to N days ago
#    FOR i IN 1..days_back LOOP
#            child_name := format('%s_%s', parent::text, to_char(base - i, 'YYYY_MM_DD'));
#
#            EXECUTE format(
#                    'CREATE TABLE IF NOT EXISTS %I PARTITION OF %s
#                       FOR VALUES FROM (%L) TO (%L);',
#                    child_name,
#                    parent::text,
#                    (base - i)::text,          -- lower bound: start of (today - i)
#                    (base - i + 1)::text       -- upper bound: start of the next day
#                    );
#        END LOOP;
#
#    RAISE NOTICE 'Created past daily partitions for %: [% .. %] (yesterday back % days)',
#        parent::text, base - days_back, base - 1, days_back;
#END$$;
#
#
#--- remove partitions older than 30 days
#CREATE OR REPLACE FUNCTION drop_partitions_older_than(
#    parent regclass,
#    days_back int DEFAULT 30
#) RETURNS void
#    LANGUAGE plpgsql AS $$
#DECLARE
#    -- We'll try to anchor on the real "today" partition; otherwise use local midnight
#    t_lower timestamptz := NULL;
#    t_upper timestamptz := NULL;
#
#    fallback_today_start timestamptz := date_trunc('day', now());
#    cutoff timestamptz;  -- drop if partition upper_bound <= cutoff
#
#    r record;
#    lower_txt text;
#    upper_txt text;
#    p_lower timestamptz;
#    p_upper timestamptz;
#BEGIN
#    -- 1) Try to discover today's partition and its [lower, upper)
#    FOR r IN
#        SELECT c.relname AS child_name, pg_get_expr(c.relpartbound, c.oid) AS bound_expr
#        FROM pg_inherits i
#                 JOIN pg_class    c ON c.oid = i.inhrelid
#        WHERE i.inhparent = parent
#        LOOP
#            SELECT (regexp_matches(r.bound_expr, 'FROM \(''([^'']+)''\)\s+TO \(''([^'']+)''\)'))[1],
#                   (regexp_matches(r.bound_expr, 'FROM \(''([^'']+)''\)\s+TO \(''([^'']+)''\)'))[2]
#            INTO lower_txt, upper_txt;
#            IF lower_txt IS NULL OR upper_txt IS NULL THEN
#                CONTINUE;
#            END IF;
#
#            EXECUTE format('SELECT %L::timestamptz, %L::timestamptz', lower_txt, upper_txt)
#                INTO p_lower, p_upper;
#
#            IF now() >= p_lower AND now() < p_upper THEN
#                t_lower := p_lower;
#                t_upper := p_upper;
#                EXIT;
#            END IF;
#        END LOOP;
#
#    -- 2) Compute cutoff
#    IF t_lower IS NOT NULL THEN
#        -- Align to your actual partition schedule (e.g., 03:30 local = 00:00 UTC)
#        cutoff := t_lower - make_interval(days => GREATEST(days_back, 0));
#    ELSE
#        -- Fallback: local day start minus N days
#        cutoff := fallback_today_start - make_interval(days => GREATEST(days_back, 0));
#    END IF;
#
#    RAISE NOTICE 'Dropping partitions in % with upper_bound <= % (days_back=%)', parent::text, cutoff, days_back;
#
#    -- 3) Walk children and drop those fully older than the cutoff
#    FOR r IN
#        SELECT c.relname AS child_name, pg_get_expr(c.relpartbound, c.oid) AS bound_expr
#        FROM pg_inherits i
#                 JOIN pg_class    c ON c.oid = i.inhrelid
#        WHERE i.inhparent = parent
#        LOOP
#            SELECT (regexp_matches(r.bound_expr, 'FROM \(''([^'']+)''\)\s+TO \(''([^'']+)''\)'))[1],
#                   (regexp_matches(r.bound_expr, 'FROM \(''([^'']+)''\)\s+TO \(''([^'']+)''\)'))[2]
#            INTO lower_txt, upper_txt;
#
#            IF lower_txt IS NULL OR upper_txt IS NULL THEN
#                RAISE NOTICE 'Skip %: cannot parse bounds: %', r.child_name, r.bound_expr;
#                CONTINUE;
#            END IF;
#
#            EXECUTE format('SELECT %L::timestamptz, %L::timestamptz', lower_txt, upper_txt)
#                INTO p_lower, p_upper;
#
#            IF p_upper <= cutoff THEN
#                RAISE NOTICE 'DROP  %  [% .. %]', r.child_name, p_lower, p_upper;
#                EXECUTE format('DROP TABLE IF EXISTS %I', r.child_name);
#            ELSE
#                RAISE NOTICE 'KEEP  %  [% .. %]', r.child_name, p_lower, p_upper;
#            END IF;
#        END LOOP;
#END$$;
#
#--- make sure i have 5 future days partition .
#
#CREATE OR REPLACE FUNCTION ensure_min_future_partitions(
#    parent regclass,
#    min_future_days int DEFAULT 5
#) RETURNS void
#    LANGUAGE plpgsql AS $$
#DECLARE
#    -- Anchor at the end of "today" if it exists; otherwise the first future lower bound;
#    -- if neither exists, fall back to local today 00:00 + 1 day.
#    anchor_upper timestamptz := NULL;
#
#    -- temp vars
#    r record;
#    lower_txt text; upper_txt text;
#    lower_bound timestamptz; upper_bound timestamptz;
#
#    future_count int := 0;
#    need int := 0;
#
#    k int;
#    start_ts timestamptz;
#    end_ts   timestamptz;
#    child_name text;
#BEGIN
#    IF min_future_days <= 0 THEN
#        RAISE NOTICE 'min_future_days=%; nothing to do.', min_future_days;
#        RETURN;
#    END IF;
#
#    -- Pass 1: find "today" (bounds containing now())
#    FOR r IN
#        SELECT c.relname AS child_name, pg_get_expr(c.relpartbound, c.oid) AS bound_expr
#        FROM pg_inherits i
#                 JOIN pg_class    c ON c.oid = i.inhrelid
#        WHERE i.inhparent = parent
#        LOOP
#            SELECT (regexp_matches(r.bound_expr, 'FROM \(''([^'']+)''\)\s+TO \(''([^'']+)''\)'))[1],
#                   (regexp_matches(r.bound_expr, 'FROM \(''([^'']+)''\)\s+TO \(''([^'']+)''\)'))[2]
#            INTO lower_txt, upper_txt;
#            IF lower_txt IS NULL OR upper_txt IS NULL THEN
#                CONTINUE;
#            END IF;
#
#            EXECUTE format('SELECT %L::timestamptz, %L::timestamptz', lower_txt, upper_txt)
#                INTO lower_bound, upper_bound;
#
#            IF now() >= lower_bound AND now() < upper_bound THEN
#                anchor_upper := upper_bound;  -- end of "today"
#                EXIT;
#            END IF;
#        END LOOP;
#
#    -- If no "today", use the nearest future lower bound as the anchor
#    IF anchor_upper IS NULL THEN
#        SELECT MIN(lb) INTO anchor_upper
#        FROM (
#                 SELECT ( (regexp_matches(pg_get_expr(c.relpartbound, c.oid),
#                                          'FROM \(''([^'']+)''\)\s+TO \(''([^'']+)''\)'))[1] )::timestamptz AS lb
#                 FROM pg_inherits i
#                          JOIN pg_class    c ON c.oid = i.inhrelid
#                 WHERE i.inhparent = parent
#             ) s
#        WHERE lb > now();
#    END IF;
#
#    -- If still NULL (no partitions at all), fall back to local midnight + 1 day
#    IF anchor_upper IS NULL THEN
#        anchor_upper := date_trunc('day', now()) + interval '1 day';
#    END IF;
#
#    -- Count how many future partitions already exist (lower_bound >= anchor_upper)
#    SELECT COUNT(*) INTO future_count
#    FROM (
#             SELECT ( (regexp_matches(pg_get_expr(c.relpartbound, c.oid),
#                                      'FROM \(''([^'']+)''\)\s+TO \(''([^'']+)''\)'))[1] )::timestamptz AS lb
#             FROM pg_inherits i
#                      JOIN pg_class    c ON c.oid = i.inhrelid
#             WHERE i.inhparent = parent
#         ) q
#    WHERE lb >= anchor_upper;
#
#    need := GREATEST(min_future_days - future_count, 0);
#    RAISE NOTICE 'Parent=%, future_count=% (>= %?), need to create=%; anchor_upper=%',
#        parent::text, future_count, min_future_days, need, anchor_upper;
#
#    IF need = 0 THEN
#        RETURN;
#    END IF;
#
#    -- Create exactly the missing future partitions at [anchor_upper + k, anchor_upper + k + 1)
#    FOR k IN 0..(need - 1) LOOP
#            start_ts := anchor_upper + make_interval(days => k);
#            end_ts   := anchor_upper + make_interval(days => k + 1);
#            child_name := format('%s_%s', parent::text, to_char(start_ts::date, 'YYYY_MM_DD'));
#
#            EXECUTE format(
#                    'CREATE TABLE IF NOT EXISTS %I PARTITION OF %s
#                       FOR VALUES FROM (%L) TO (%L);',
#                    child_name, parent::text, start_ts, end_ts
#                    );
#            RAISE NOTICE 'Created (or exists) %  [FROM % TO %)', child_name, start_ts, end_ts;
#        END LOOP;
#END$$;
#
#DO $$
#DECLARE
#  tbl text;
#  tables text[] := ARRAY[
#    'request_messages',
#    'deliver_messages',
#    'user_requests',
#    'complete_deliver_messages',
#    'submit_messages'
#  ];
#BEGIN
#  FOREACH tbl IN ARRAY tables LOOP
#    RAISE NOTICE 'Creating daily partitions for %', tbl;
#    EXECUTE format(
#      'SELECT create_daily_partitions_from_today(%L::regclass, 35);',
#      tbl
#    );
#  END LOOP;
#END$$;
#
#-- Non-unique, non-partial indexes can be created on the parent
#
#-- request_messages
#CREATE INDEX IF NOT EXISTS idx_request_messages_client_group
#  ON request_messages (client_id, group_id);
#
#CREATE INDEX IF NOT EXISTS idx_request_messages_req_client_group
#  ON request_messages (request_id, client_id);
#
#-- Accelerate pending cut-off lookup used by retry worker
#CREATE INDEX IF NOT EXISTS idx_request_messages_pending_cutoff
#  ON request_messages (create_at)
#  WHERE overall_status = 'PENDING';
#
#CREATE INDEX IF NOT EXISTS idx_submit_messages_group
#  ON submit_messages (group_id);
#
#-- Helps ORDER BY (group_id, message_seq) scans
#CREATE INDEX IF NOT EXISTS idx_submit_messages_group_seq
#  ON submit_messages (group_id, message_seq);
#
#-- deliver_messages (PK per partition; keep this index at parent for compatibility)
#CREATE INDEX IF NOT EXISTS idx_deliver_messages_id
#  ON deliver_messages (id);
#
#-- complete_deliver_messages
#CREATE INDEX IF NOT EXISTS idx_complete_deliver_messages_client
#  ON complete_deliver_messages (client_id, receiver);
#
#-- Prefixes lookup by client
#CREATE INDEX IF NOT EXISTS idx_prefixes_client_id
#  ON prefixes (client_id);
#
#-- Per-partition unique/partial indexes for request_messages are applied in the partition creation loop above.
#
#SQL
#fi
#echo "Initialization complete."
