-- auto-generated definition
create table charging_history
(
    time          timestamp with time zone not null,
    total_wh      double precision         not null,
    wh_difference double precision         not null
);

alter table charging_history
    owner to energy_user;

-- auto-generated definition
create table comed_price
(
    time  timestamp with time zone not null,
    price double precision         not null
);

alter table comed_price
    owner to energy_user;

create or replace view tesla_cost_view as
SELECT cs.time charge_time, cp.time price_time, cs.wh_difference, cp.price, (cp.price + 0.06) * cs.wh_difference as cost
from charging_history cs join  comed_price cp
                               on  to_char((cs.time), 'mm/dd/yyyy hh24') = to_char((cp.time), 'mm/dd/yyyy hh24');