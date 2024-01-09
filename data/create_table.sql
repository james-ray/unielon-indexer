DROP TABLE IF EXISTS block;
DROP TABLE IF EXISTS cardinals_revert;
DROP TABLE IF EXISTS cardinals_info;
DROP TABLE IF EXISTS drc20_info;
DROP TABLE IF EXISTS swap_info;
DROP TABLE IF EXISTS swap_liquidity;

CREATE TABLE block
(
    block_hash   VARCHAR(255) PRIMARY KEY,
    block_number INT,
    create_date  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date  TIMESTAMP
);

CREATE TABLE cardinals_revert
(
    id           INT AUTO_INCREMENT PRIMARY KEY,
    tick         VARCHAR(255),
    from_address VARCHAR(255),
    to_address   VARCHAR(255),
    amt          VARCHAR(255),
    block_number INT,
    create_date  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cardinals_info
(
    id
    INT
    AUTO_INCREMENT
    PRIMARY
    KEY,
    order_id
    VARCHAR
(
    255
) NOT NULL,
    p VARCHAR
(
    255
) NOT NULL,
    op VARCHAR
(
    255
) NOT NULL,
    tick VARCHAR
(
    255
) NOT NULL,
    amt VARCHAR
(
    255
) NOT NULL,
    max_ VARCHAR
(
    255
) NOT NULL,
    lim_ VARCHAR
(
    255
) NOT NULL,
    dec_ INT NOT NULL,
    burn_ VARCHAR
(
    255
) NOT NULL,
    func_ VARCHAR
(
    255
) NOT NULL,
    receive_address VARCHAR
(
    255
) NOT NULL,
    fee_address VARCHAR
(
    255
) NOT NULL,
    to_address VARCHAR
(
    255
) NOT NULL,
    drc20_tx_hash VARCHAR
(
    255
) NOT NULL,
    repeat_mint INT NOT NULL,
    block_number INT,
    block_hash VARCHAR
(
    255
),
    order_status INT,
    create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    update_date TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS drc20_info
(
    id
    INT
    AUTO_INCREMENT
    PRIMARY
    KEY,
    tick
    VARCHAR
(
    255
) NOT NULL,
    max_ bigint NOT NULL,
    lim_ bigint NOT NULL,
    receive_address VARCHAR
(
    255
) NOT NULL,
    drc20_tx_hash VARCHAR
(
    255
) NOT NULL,
    create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    update_date TIMESTAMP
    );

CREATE TABLE swap_info
(
    id                INT AUTO_INCREMENT PRIMARY KEY,
    order_id          INT,
    op                VARCHAR(255),
    tick0             VARCHAR(255),
    tick1             VARCHAR(255),
    amt0              VARCHAR(255),
    amt1              VARCHAR(255),
    amt0_min          VARCHAR(255),
    amt1_min          VARCHAR(255),
    liquidity         VARCHAR(255),
    fee_address       VARCHAR(255),
    holder_address    VARCHAR(255),
    swap_tx_hash      VARCHAR(255),
    amt0_out          VARCHAR(255),
    amt1_out          VARCHAR(255),
    swap_block_hash   VARCHAR(255),
    swap_block_number INT,
    order_status      INT,
    create_date       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date       TIMESTAMP
);

CREATE TABLE swap_liquidity
(
    id               INT AUTO_INCREMENT PRIMARY KEY,
    tick             VARCHAR(255),
    tick0            VARCHAR(255),
    tick1            VARCHAR(255),
    holder_address   VARCHAR(255),
    reserves_address VARCHAR(255),
    liquidity_total  VARCHAR(255),
    create_date  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    update_date      TIMESTAMP
);

DROP TRIGGER IF EXISTS update_block;
DROP TRIGGER IF EXISTS update_cardinals_revert;
DROP TRIGGER IF EXISTS update_cardinals_info;
DROP TRIGGER IF EXISTS update_drc20_info;
DROP TRIGGER IF EXISTS update_swap_info;
DROP TRIGGER IF EXISTS update_swap_liquidity;

CREATE TRIGGER update_block
    AFTER UPDATE ON block
BEGIN
    UPDATE block SET update_date = CURRENT_TIMESTAMP WHERE block_hash = NEW.block_hash;
END;

CREATE TRIGGER update_cardinals_revert
    AFTER UPDATE ON cardinals_revert
BEGIN
    UPDATE cardinals_revert SET update_date = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_cardinals_info
    AFTER UPDATE ON cardinals_info
BEGIN
    UPDATE cardinals_info SET update_date = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_drc20_info
    AFTER UPDATE ON drc20_info
BEGIN
    UPDATE drc20_info SET update_date = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_swap_info
    AFTER UPDATE ON swap_info
BEGIN
    UPDATE swap_info SET update_date = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER update_swap_liquidity
    AFTER UPDATE ON swap_liquidity
BEGIN
    UPDATE swap_liquidity SET update_date = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;