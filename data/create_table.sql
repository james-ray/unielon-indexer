CREATE TABLE block (
  block_hash VARCHAR(255) PRIMARY KEY,
  block_number INT
);

CREATE TABLE cardinals_revert (
  id INT AUTO_INCREMENT PRIMARY KEY,
  tick VARCHAR(255),
  from_address VARCHAR(255),
  to_address VARCHAR(255),
  amt VARCHAR(255),
  block_number INT
);


CREATE TABLE IF NOT EXISTS cardinals_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL,
    p VARCHAR(255) NOT NULL,
    op VARCHAR(255) NOT NULL,
    tick VARCHAR(255) NOT NULL,
    amt VARCHAR(255) NOT NULL,
    max_ VARCHAR(255) NOT NULL,
    lim_ VARCHAR(255) NOT NULL,
    dec_ INT NOT NULL,
    burn_ VARCHAR(255) NOT NULL,
    func_ VARCHAR(255) NOT NULL,
    receive_address VARCHAR(255) NOT NULL,
    fee_address VARCHAR(255) NOT NULL,
    to_address VARCHAR(255) NOT NULL,
    drc20_tx_hash VARCHAR(255) NOT NULL,
    repeat_mint INT NOT NULL,
    block_number INT,
    block_hash VARCHAR(255),
    order_status INT
);

CREATE TABLE IF NOT EXISTS drc20_info (
    id INT AUTO_INCREMENT PRIMARY KEY,
    tick VARCHAR(255) NOT NULL,
    max_ bigint NOT NULL,
    lim_ bigint NOT NULL,
    receive_address VARCHAR(255) NOT NULL,
    drc20_tx_hash VARCHAR(255) NOT NULL
);


CREATE TABLE swap_info (
  id INT AUTO_INCREMENT PRIMARY KEY,
  order_id INT,
  op VARCHAR(255),
  tick0 VARCHAR(255),
  tick1 VARCHAR(255),
  amt0 VARCHAR(255),
  amt1 VARCHAR(255),
  amt0_min VARCHAR(255),
  amt1_min VARCHAR(255),
  liquidity VARCHAR(255),
  fee_address VARCHAR(255),
  holder_address VARCHAR(255),
  swap_tx_hash VARCHAR(255),
  amt0_out VARCHAR(255),
  amt1_out VARCHAR(255),
  swap_block_hash VARCHAR(255),
  swap_block_number INT,
  order_status INT
);

CREATE TABLE swap_liquidity (
  id INT AUTO_INCREMENT PRIMARY KEY,
  tick VARCHAR(255),
  tick0 VARCHAR(255),
  tick1 VARCHAR(255),
  holder_address VARCHAR(255),
  reserves_address VARCHAR(255),
  liquidity_total VARCHAR(255)
);


