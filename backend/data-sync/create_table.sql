create table bitcoin_block
(
    hash              varchar(128),
    `size`            bigint,
    `stripped_size`   bigint,
    weight            bigint,
    number            bigint,
    version           bigint,
    `merkle_root`     varchar(128),
    timestamp         timestamp,
    nonce             varchar(64),
    bits              varchar(64),
    coinbase_param    text,
    transaction_count bigint,
    unique key uniq_hash (hash);
);


CREATE TABLE bitcoin_transaction
(
    `hash`          varchar(128) NOT NULL COMMENT 'The hash of this transaction',
    size            BIGINT COMMENT 'The size of this transaction in bytes',
    virtual_size    BIGINT COMMENT 'The virtual transaction size (differs from size for witness transactions)',
    version         BIGINT COMMENT 'Protocol version specified in block which contained this transaction',
    lock_time       BIGINT COMMENT 'Earliest time that miners can include the transaction in their hashing of the Merkle root to attach it in the latest block of the blockchain',
    block_hash      varchar(128) NOT NULL COMMENT 'Hash of the block which contains this transaction',
    block_number    BIGINT       NOT NULL COMMENT 'Number of the block which contains this transaction',
    block_timestamp TIMESTAMP    NOT NULL COMMENT 'Timestamp of the block which contains this transaction',
    input_count     BIGINT COMMENT 'The number of inputs in the transaction',
    output_count    BIGINT COMMENT 'The number of outputs in the transaction',
    input_value     NUMERIC COMMENT 'Total value of inputs in the transaction',
    output_value    NUMERIC COMMENT 'Total value of outputs in the transaction',
    is_coinbase     BOOL COMMENT 'true if this transaction is a coinbase transaction',
    fee             NUMERIC COMMENT 'The fee paid by this transaction',
    index idx_block_timestamp (block_timestamp),
    unique key uniq_hash (hash)
);


CREATE TABLE bitcoin_transaction_input
(
    transaction_hash       varchar(128) comment 'The hash of the transaction which contains the output that this input spends',
    `index`                BIGINT NOT NULL comment '0-indexed number of an input within a transaction',
    spent_transaction_hash varchar(128) comment 'The hash of the transaction which contains the output that this input spends',
    spent_output_index     BIGINT comment 'The index of the output this input spends',
    script_asm             text comment 'Symbolic representation of the bitcoin s script language op-codes',
    script_hex             text comment 'Hexadecimal representation of the bitcoin s script language op-codes',
    sequence               BIGINT comment 'A number intended to allow unconfirmed time-locked transactions to be updated before being finalized; not currently used except to disable locktime in a transaction',
    required_signatures    BIGINT comment 'The number of signatures required to authorize the spent output',
    type                   varchar(32) comment 'The address type of the spent output',
    addresses              json comment 'Addresses which own the spent output',
    value                  NUMERIC comment 'The value in base currency attached to the spent output',
    unique index uniq_tx_idx(transaction_hash,`index`)
);

CREATE TABLE bitcoin_transaction_output
(
    transaction_hash    varchar(128) comment 'The hash of the transaction which contains the output that this input spends',
    `index`             BIGINT NOT NULL comment '0-indexed number of an output within a transaction used by a later transaction to refer to that specific output',
    script_asm          text comment 'Symbolic representation of the bitcoin s script language op-codes',
    script_hex          text comment 'Hexadecimal representation of the bitcoin s script language op-codes',
    required_signatures BIGINT comment 'The number of signatures required to authorize spending of this output',
    type                varchar(32) comment 'The address type of the output',
    addresses           json comment 'Addresses which own this output',
    `value`               NUMERIC comment 'The value in base currency attached to this output',
    unique index uniq_tx_idx(transaction_hash,`index`)
);

