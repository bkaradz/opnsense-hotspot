-- 1️⃣ Providers table
CREATE TABLE IF NOT EXISTS providers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE
);

-- 2️⃣ Voucher groups table
CREATE TABLE IF NOT EXISTS voucher_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    provider_id INTEGER NOT NULL REFERENCES providers(id),
    name TEXT NOT NULL,
    UNIQUE(provider_id, name)
);

-- 3️⃣ Vouchers table
CREATE TABLE IF NOT EXISTS vouchers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL REFERENCES voucher_groups(id),
    username TEXT NOT NULL,
    validity INTEGER NOT NULL,
    expirytime INTEGER,
    starttime REAL,
    endtime INTEGER,
    printed TEXT NOT NULL DEFAULT "unprinted",
    state TEXT NOT NULL,
    UNIQUE(username, group_id)
);

-- 4️⃣ FTS5 virtual table for username search
CREATE VIRTUAL TABLE IF NOT EXISTS vouchers_fts
USING fts5(username, content='vouchers', content_rowid='id', prefix='2 3 4 5 6 7', tokenize="unicode61");

-- 5️⃣ Triggers to sync FTS index
CREATE TRIGGER IF NOT EXISTS vouchers_ai AFTER INSERT ON vouchers BEGIN
  INSERT INTO vouchers_fts(rowid, username) VALUES (new.id, new.username);
END;

CREATE TRIGGER IF NOT EXISTS vouchers_au AFTER UPDATE ON vouchers BEGIN
  UPDATE vouchers_fts SET username = new.username WHERE rowid = old.id;
END;

CREATE TRIGGER IF NOT EXISTS vouchers_ad AFTER DELETE ON vouchers BEGIN
  DELETE FROM vouchers_fts WHERE rowid = old.id;
END;
