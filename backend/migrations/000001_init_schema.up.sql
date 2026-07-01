CREATE TABLE users (
  id uuid PRIMARY KEY,
  name varchar(100) NOT NULL,
  mobile varchar(20) NOT NULL UNIQUE,
  email varchar(100) NOT NULL UNIQUE,
  password_hash varchar(255) NOT NULL,
  role varchar(20) DEFAULT 'user',
  status varchar(20) DEFAULT 'pending',
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp DEFAULT NULL
);

CREATE TABLE sessions (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL,
  access_token text NOT NULL,
  refresh_token varchar(255) NOT NULL UNIQUE,
  expires_at timestamp NOT NULL,
  ip_address varchar(45) DEFAULT NULL,
  user_agent varchar(255) DEFAULT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE personal_details (
  user_id uuid PRIMARY KEY,
  father_name varchar(100) DEFAULT NULL,
  mother_name varchar(100) DEFAULT NULL,
  country varchar(100) DEFAULT NULL,
  state varchar(100) DEFAULT NULL,
  city varchar(100) DEFAULT NULL,
  address text,
  pincode varchar(20) DEFAULT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_personal_details_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE bank_details (
  user_id uuid PRIMARY KEY,
  ifsc varchar(20) DEFAULT NULL,
  bank_name varchar(100) DEFAULT NULL,
  branch varchar(100) DEFAULT NULL,
  account_number varchar(50) UNIQUE DEFAULT NULL,
  income_range varchar(50) DEFAULT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_bank_details_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE wallets (
  user_id uuid PRIMARY KEY,
  wallet_balance decimal(15,2) DEFAULT 0.00,
  blocked_balance decimal(15,2) DEFAULT 0.00,
  available_balance decimal(15,2) DEFAULT 0.00,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_wallets_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE transactions (
  id uuid PRIMARY KEY,
  user_id uuid DEFAULT NULL,
  type varchar(50) DEFAULT NULL,
  amount decimal(15,2) DEFAULT NULL,
  reference_id varchar(100) UNIQUE DEFAULT NULL,
  description varchar(255) DEFAULT NULL,
  status varchar(20) DEFAULT 'completed',
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_transactions_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE segments (
  id bigserial PRIMARY KEY,
  name varchar(20) UNIQUE DEFAULT NULL
);

CREATE TABLE shares (
  id uuid PRIMARY KEY,
  symbol varchar(50) NOT NULL UNIQUE,
  name varchar(150) NOT NULL,
  price decimal(15,2) NOT NULL,
  segment_id bigint DEFAULT NULL,
  total_shares bigint DEFAULT NULL,
  available_shares bigint DEFAULT NULL,
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
  deleted_at timestamp DEFAULT NULL,
  CONSTRAINT fk_shares_segment FOREIGN KEY (segment_id) REFERENCES segments (id)
);

CREATE TABLE trades (
  id uuid PRIMARY KEY,
  user_id uuid DEFAULT NULL,
  share_id uuid DEFAULT NULL,
  quantity bigint NOT NULL,
  price decimal(15,2) NOT NULL,
  type varchar(20) NOT NULL,
  status varchar(20) DEFAULT 'completed',
  created_at timestamp DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_trades_share FOREIGN KEY (share_id) REFERENCES shares (id),
  CONSTRAINT fk_trades_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
