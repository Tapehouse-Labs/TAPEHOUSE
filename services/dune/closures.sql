-- SPDX-License-Identifier: MIT OR Apache-2.0
-- Each closure of Robinhood Chain's Chainlink equity feeds since the chain launched: the feed's last round before a
-- silence of a day or more and its first round after, and the move between them, from each aggregator's
-- AnswerUpdated(int256 indexed current, uint256 indexed roundId, uint256 updatedAt). The aggregators are those the
-- feeds in deployments/4663.json pointed to on 2 October 2026.
with feeds (asset, aggregator) as (
    values
        ('NVDA', 0xC9d16E4f2569b9E3ea0468fD85844953713DC2a2),
        ('TSLA', 0x7A6b81ba7FbCB90104d8C496158Cf383cD7233b1),
        ('AAPL', 0xBb11A21267cFDb63d4935d99a499133DD1744ACb),
        ('MSFT', 0xc3b117F52cf17Dd4369eaF5eaf7cF0E2f91b4E30),
        ('GOOGL', 0x11eD6d598eF565DDA86fAfE7E779303e7CC6b2Bd),
        ('SPY', 0x78BCB218fA04B9b3a278eBc865Ed320BF8DEFBAc)
),
rounds as (
    select
        f.asset,
        from_unixtime(bytearray_to_uint256(l.data)) as updated_at,
        cast(bytearray_to_int256(l.topic1) as double) / 1e8 as answer
    from robinhood.logs as l
    inner join feeds as f on l.contract_address = f.aggregator
    where l.topic0 = 0x0559884fd3a460db3073b7fc896cc77986f16e378210ded43186175bf646fc5f
        and l.block_time >= timestamp '2026-07-01'
),
ordered as (
    select
        asset,
        updated_at,
        answer,
        lag(updated_at) over (partition by asset order by updated_at) as last_updated_at,
        lag(answer) over (partition by asset order by updated_at) as last_answer
    from rounds
)
select
    asset,
    last_updated_at as last_round,
    updated_at as reopening_round,
    last_answer,
    answer,
    answer / last_answer - 1 as move
from ordered
where updated_at > last_updated_at + interval '1' day
order by reopening_round, asset
