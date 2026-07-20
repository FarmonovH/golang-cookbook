package main

import (
	"errors"
	"fmt"
)



func main(){
	
}


oracle

SELECT *  FROM ibs.transacts_cur@iabs;

select /* +  index (t IDX_TRANSACTS_ACC_ID)
                   index (a XPKACCOUNTS) index (l XPKLEADS) */
 th.amount,
 th.op_dc,
 lc.cl_name,
 lc.co_name,
 substr(lc.cl_acc, -20) cl_acc,
 substr(lc.co_acc, -20) co_acc,
 lc.code_currency
  from ibs.transacts_cur@iabs th, Leads_Cur lc
 where lc.id = th.lead_id
union all
select /* +  index (t IDX_TRANSACTS_ACC_ID)
                   index (a XPKACCOUNTS) index (l XPKLEADS) */
 th.amount,
 th.op_dc,
 lc.cl_name,
 lc.co_name,
 substr(lc.cl_acc, -20) cl_acc,
 substr(lc.co_acc, -20) co_acc,
 lc.code_currency
  from ibs.transacts_cur@iabs th, ibs.Leads_fs@iabs lc
 where lc.id = th.lead_id


select
       lc.date_enter,
       a.client_id,
       lc.date_execute,
       lc.state_id,
       lc.id,
       th.amount,
       th.op_dc,
       lc.cl_name,
       lc.co_name,
       substr(lc.cl_acc, -20) cl_acc,
       substr(lc.co_acc, -20) co_acc,
       lc.code_currency
from ibs.transacts_cur@iabs th
join leads_cur lc
  on lc.id = th.lead_id
left join accounts a
  on a.code = lc.co_acc where lc.state_id = 41 AND lc.date_execute > TO_DATE('2026-06-03 12:17:39', 'YYYY-MM-DD HH24:MI:SS')
union all
select
       lc.date_enter,
       a.client_id,
       lc.date_execute,
       lc.state_id,
       lc.id,
       th.amount,
       th.op_dc,
       lc.cl_name,
       lc.co_name,
       substr(lc.cl_acc, -20) cl_acc,
       substr(lc.co_acc, -20) co_acc,
       lc.code_currency
from ibs.transacts_cur@iabs th
join ibs.Leads_fs@iabs lc
  on lc.id = th.lead_id
left join accounts a
  on a.code = lc.co_acc where lc.state_id = 41 AND lc.date_execute > TO_DATE('2026-06-03 12:17:39', 'YYYY-MM-DD HH24:MI:SS');


select max(date_execute), min(date_execute) from ibs.leads_cur.@iabs where state_id = 41;


select
    lc.date_enter,
    a.client_id,
    lc.date_execute,
    lc.state_id,
    lc.id,
    th.amount,
    th.op_dc,
    lc.cl_name,
    lc.co_name,
    substr(lc.cl_acc, -20) cl_acc,
    substr(lc.co_acc, -20) co_acc,
    lc.code_currency
from ibs.transacts_cur@iabs th
join leads_cur lc
on lc.id = th.lead_id
left join accounts a
on a.code = lc.co_acc
where lc.state_id = 41 AND  lc.date_execute > TO_DATE('2026-06-03 12:17:39', 'YYYY-MM-DD HH24:MI:SS')

union all

select
    lc.date_enter,
    a.client_id,
    lc.date_execute,
    lc.state_id,
    lc.id,
    th.amount,
    th.op_dc,
    lc.cl_name,
    lc.co_name,
    substr(lc.cl_acc, -20) cl_acc,
    substr(lc.co_acc, -20) co_acc,
    lc.code_currency
from ibs.transacts_cur@iabs th
join ibs.Leads_fs@iabs lc
on lc.id = th.lead_id
left join accounts a
on a.code = lc.co_acc
where lc.state_id = 41 AND  lc.date_execute > TO_DATE('2026-06-03 12:17:39', 'YYYY-MM-DD HH24:MI:SS')
------------------------------------------------------------------------------------
-- 204

WITH const AS (

    SELECT

        TO_DATE('03.10.2025', 'DD.MM.YYYY') AS report_date,

        TO_DATE('02.06.2025', 'DD.MM.YYYY') AS sleeping_start,

        TO_DATE('02.10.2025', 'DD.MM.YYYY') AS sleeping_end,

        2000000000  AS sleeping_limit,

        40000000000 AS spike_limit

    FROM dual

),

acc AS (

    SELECT DISTINCT a.client_uid, a.code

    FROM accounts a

    JOIN client_current cc ON cc.client_uid = a.client_uid

    WHERE a.code_currency = '000'

      AND a.code_coa IN ('20208','20218','20210','20212','20214','20296')

      AND cc.subject   = 'J'

      AND cc.condition = 'A'

),


sleeping_clients AS (
    SELECT client_uid
    FROM acc
    GROUP BY client_uid
    HAVING
        SUM(
            NVL(GET_DEBIT_ALL(code, (SELECT sleeping_end   FROM const)), 0)
          - NVL(GET_DEBIT_ALL(code, (SELECT sleeping_start FROM const) - 1), 0)
          + NVL(GET_CREDIT_ALL(code, (SELECT sleeping_end   FROM const)), 0)
          - NVL(GET_CREDIT_ALL(code, (SELECT sleeping_start FROM const) - 1), 0)
        ) <= (SELECT sleeping_limit FROM const)
),
-- 2. Генерируем даты после 04.11.2025 и считаем дневной оборот
daily_after AS (
    SELECT
        a.client_uid,
        dt.d AS op_date,
        SUM(
            NVL(GET_DEBIT_ALL(a.code, dt.d), 0)     - NVL(GET_DEBIT_ALL(a.code, dt.d - 1), 0) +
            NVL(GET_CREDIT_ALL(a.code, dt.d), 0)     - NVL(GET_CREDIT_ALL(a.code, dt.d - 1), 0)
        ) AS day_amt
    FROM acc a
    CROSS JOIN (
        SELECT (SELECT sleeping_end FROM const) + LEVEL - 1 AS d
        FROM dual
        CONNECT BY LEVEL <= (SELECT report_date FROM const) - (SELECT sleeping_end FROM const) + 10  -- +10 на всякий
    ) dt
    WHERE dt.d <= (SELECT report_date FROM const)
    GROUP BY a.client_uid, dt.d
    HAVING SUM(
            NVL(GET_DEBIT_ALL(a.code, dt.d), 0)     - NVL(GET_DEBIT_ALL(a.code, dt.d - 1), 0) +
            NVL(GET_CREDIT_ALL(a.code, dt.d), 0)     - NVL(GET_CREDIT_ALL(a.code, dt.d - 1), 0)
           ) > 0
),
-- 3. Ищем два подряд дня ≥ 400 млн
spikes AS (
    SELECT
        d1.client_uid,
        d1.op_date          AS day1,
        d1.op_date + 30      AS day2,
        d1.day_amt + NVL(d2.day_amt, 0) AS two_days_amt
    FROM daily_after d1
    LEFT JOIN daily_after d2
           ON d2.client_uid = d1.client_uid
         -- AND d2.op_date = d1.op_date +
    WHERE d1.day_amt + NVL(d2.day_amt, 0) >= (SELECT spike_limit FROM const)
)
-- ФИНАЛ — всё в SELECT, чтобы не было ORA-01791
SELECT
    cc.*,
    s.day1,
    s.day2,
    ROUND(s.two_days_amt / 1000000, 2)          AS two_days_mln,
    s.two_days_amt                              AS sort_amt,
    ROUND((
        SELECT SUM(
                  NVL(GET_DEBIT_ALL(a.code, c.sleeping_end), 0)
                - NVL(GET_DEBIT_ALL(a.code, c.sleeping_start - 1), 0)
                + NVL(GET_CREDIT_ALL(a.code, c.sleeping_end), 0)
                - NVL(GET_CREDIT_ALL(a.code, c.sleeping_start - 1), 0)
               )
        FROM acc a
        WHERE a.client_uid = cc.client_uid
    ) / 1000000, 2)                             AS sleeping_90d_mln
FROM client_current cc
JOIN sleeping_clients sc ON sc.client_uid = cc.client_uid
JOIN spikes s            ON s.client_uid  = cc.client_uid
CROSS JOIN const c
ORDER BY s.two_days_amt DESC;

-- 204 END
-----------------------------------------------------------------------------------------------------------------
SELECT object_name, object_type
FROM user_objects
WHERE UPPER(object_name) = 'GET_CREDIT_ALL';



select *
from ALL_SCHEDULER_JOBS;

SELECT column_name, data_type
FROM all_tab_columns@iabs
WHERE table_name = 'LN_AVP_DEBTOR_SOURCES';

SELECT BANK.GET_CREDIT_ALL@iabs('12345', SYSDATE)
FROM dual;

SELECT owner,
       object_name,
       procedure_name
FROM all_procedures@iabs
WHERE UPPER(procedure_name) = 'GET_CREDIT_ALL';

SELECT BANK.GET_CREDIT_ALL@iabs('12345', SYSDATE)
FROM dual;


select * from saldo;


CREATE OR REPLACE FUNCTION GET_DEBIT_ALL(a_code VARCHAR2, o_day DATE) RETURN NUMBER
AS
    result NUMBER;
BEGIN
    SELECT --+ index (sl UK_SALDO_ACCOUNT_DAY)
        TURNOVER_ALL_DEBIT
    into result
    FROM SALDO sl
    WHERE ACCOUNT_CODE = a_code AND OPER_DAY = ADD_MONTHS(o_day, -3)
    AND ROWNUM = 1;
    RETURN result;
END;

CREATE OR REPLACE FUNCTION GET_CREDIT_ALL(a_code VARCHAR2, o_day DATE) RETURN NUMBER
AS
    result NUMBER;
BEGIN
    SELECT --+ index (sl UK_SALDO_ACCOUNT_DAY)
        TURNOVER_ALL_CREDIT
    into result
    FROM SALDO sl
    WHERE ACCOUNT_CODE = a_code AND OPER_DAY = ADD_MONTHS(o_day, -3)
    AND ROWNUM = 1;
    RETURN result;
END;

SELECT
    (select --+ index (sl UK_SALDO_ACCOUNT_DAY)
oper_day, TURNOVER_ALL_CREDIT, TURNOVER_ALL_DEBIT
from saldo sl
where account_code='110000020814000900009013001'
and rownum = 1)
FROM DUAL;
SELECT GET_CREDIT_ALL('110000020814000900009013001', SYSDATE -2 ) FROM DUAL;

SELECT * FROM IBS.SALDO@iabs WHERE ACCOUNT_CODE = '110200096381000699963675001'  AND OPER_DAY >= ADD_MONTHS(SYSDATE -2, -3) AND TURNOVER_ALL_DEBIT > 0;


SELECT MIN(OPER_DAY),
       MAX(OPER_DAY)
FROM IBS.SALDO@iabs
WHERE ACCOUNT_CODE = '110000020814000900009013001';

SELECT * FROM SALDO WHERE ACCOUNT_CODE = '110000020814000900009013001';

SELECT *
FROM IBS.SALDO@iabs ORDER BY OPER_DAY DESC ;

select * from ALL_IND_COLUMNS@iabs where table_name = 'SALDO';


select --+ index (sl UK_SALDO_ACCOUNT_DAY)
oper_day, TURNOVER_ALL_CREDIT, TURNOVER_ALL_DEBIT
from saldo sl
where account_code='110000020814000900009013001'
and rownum = 1;

select /*+ FIRST_ROWS(10) */
oper_day, TURNOVER_ALL_CREDIT, TURNOVER_ALL_DEBIT
from saldo sl
where account_code='110000020814000900009013001';



WITH const AS (
    SELECT
        TO_DATE('03.10.2025', 'DD.MM.YYYY') AS report_date,
        TO_DATE('02.06.2025', 'DD.MM.YYYY') AS sleeping_start,
        TO_DATE('02.10.2025', 'DD.MM.YYYY') AS sleeping_end,
        2000000000  AS sleeping_limit,
        40000000000 AS spike_limit
    FROM dual
),

-- 🔹 SALDO'dan kerakli ma'lumotlarni oldindan olish
saldo_cache AS (
    SELECT
        account_code,
        oper_day,
        turnover_all_debit,
        turnover_all_credit
    FROM saldo
    WHERE oper_day >= ADD_MONTHS((SELECT sleeping_end FROM const), -3)
),

-- 🔹 Hisoblar ro'yxati
acc AS (
    SELECT DISTINCT a.client_uid, a.code
    FROM accounts a
    JOIN client_current cc ON cc.client_uid = a.client_uid
    WHERE a.code_currency = '000'
      AND a.code_coa IN ('20208','20218','20210','20212','20214','20296')
      AND cc.subject   = 'J'
      AND cc.condition = 'A'
),

-- 🔹 Uyqu holati: 3 oylik oborot kam bo'lgan klientlar
sleeping_clients AS (
    SELECT
        a.client_uid,
        SUM(
            NVL(s.turnover_all_debit, 0) + NVL(s.turnover_all_credit, 0)
        ) AS sleeping_total
    FROM acc a
    LEFT JOIN saldo_cache s
        ON s.account_code = a.code
        AND s.oper_day >= (SELECT sleeping_start FROM const)
        AND s.oper_day <= (SELECT sleeping_end FROM const)
    GROUP BY a.client_uid
    HAVING SUM(NVL(s.turnover_all_debit, 0) + NVL(s.turnover_all_credit, 0))
           <= (SELECT sleeping_limit FROM const)
),

-- 🔹 Balanslarni oldindan hisoblash (funksiyalar chaqirish o'rniga)
daily_with_balances AS (
    SELECT
        a.client_uid,
        dt.d AS op_date,
        a.code,
        -- Bugungi balans
        (
            SELECT turnover_all_debit
            FROM saldo_cache s
            WHERE s.account_code = a.code
              AND s.oper_day <= dt.d
            ORDER BY s.oper_day DESC
            FETCH FIRST 1 ROW ONLY
        ) AS debit_today,
        -- Kechagi balans
        (
            SELECT turnover_all_debit
            FROM saldo_cache s
            WHERE s.account_code = a.code
              AND s.oper_day <= dt.d - 1
            ORDER BY s.oper_day DESC
            FETCH FIRST 1 ROW ONLY
        ) AS debit_yesterday,
        -- Bugungi kredit
        (
            SELECT turnover_all_credit
            FROM saldo_cache s
            WHERE s.account_code = a.code
              AND s.oper_day <= dt.d
            ORDER BY s.oper_day DESC
            FETCH FIRST 1 ROW ONLY
        ) AS credit_today,
        -- Kechagi kredit
        (
            SELECT turnover_all_credit
            FROM saldo_cache s
            WHERE s.account_code = a.code
              AND s.oper_day <= dt.d - 1
            ORDER BY s.oper_day DESC
            FETCH FIRST 1 ROW ONLY
        ) AS credit_yesterday
    FROM acc a
    CROSS JOIN (
        SELECT (SELECT sleeping_end FROM const) + LEVEL - 1 AS d
        FROM dual
        CONNECT BY LEVEL <= (
            (SELECT report_date FROM const) -
            (SELECT sleeping_end FROM const) + 10
        )
    ) dt
),

-- 🔹 Kunlik oborot hisobini
daily_after AS (
    SELECT
        client_uid,
        op_date,
        SUM(
            NVL(debit_today, 0) - NVL(debit_yesterday, 0) +
            NVL(credit_today, 0) - NVL(credit_yesterday, 0)
        ) AS day_amt
    FROM daily_with_balances
    WHERE op_date > (SELECT sleeping_end FROM const)
    GROUP BY client_uid, op_date
    HAVING SUM(
        NVL(debit_today, 0) - NVL(debit_yesterday, 0) +
        NVL(credit_today, 0) - NVL(credit_yesterday, 0)
    ) > 0
),

-- 🔹 Ikki kunlik piki (spike) - ketma-ket yoki har qanday)
spikes AS (
    SELECT
        d1.client_uid,
        d1.op_date AS day1,
        d2.op_date AS day2,
        d1.day_amt + NVL(d2.day_amt, 0) AS two_days_amt
    FROM daily_after d1
    LEFT JOIN daily_after d2
        ON d2.client_uid = d1.client_uid
        AND d2.op_date = d1.op_date + 1  -- ✅ Ketma-ket kunlar
    WHERE d1.day_amt + NVL(d2.day_amt, 0) >= (SELECT spike_limit FROM const)
),

-- 🔹 3 oylik uyqu davri oborot
sleeping_90_days AS (
    SELECT
        a.client_uid,
        SUM(
            NVL(s.turnover_all_debit, 0) + NVL(s.turnover_all_credit, 0)
        ) AS sleeping_90d_amt
    FROM acc a
    LEFT JOIN saldo_cache s
        ON s.account_code = a.code
        AND s.oper_day >= (SELECT sleeping_start FROM const)
        AND s.oper_day <= (SELECT sleeping_end FROM const)
    GROUP BY a.client_uid
)

-- 🔹 FINAL: Hamma ma'lumotlarni birlashtirib chiqarish
SELECT
    cc.*,
    s.day1,
    s.day2,
    ROUND(s.two_days_amt / 1000000, 2) AS two_days_mln,
    s.two_days_amt AS sort_amt,
    ROUND(s90.sleeping_90d_amt / 1000000, 2) AS sleeping_90d_mln
FROM client_current cc
INNER JOIN sleeping_clients sc ON sc.client_uid = cc.client_uid
INNER JOIN spikes s ON s.client_uid = cc.client_uid
INNER JOIN sleeping_90_days s90 ON s90.client_uid = cc.client_uid
ORDER BY s.two_days_amt DESC;


WITH const AS (
    SELECT
        TO_DATE('03.10.2025', 'DD.MM.YYYY') AS report_date,
        TO_DATE('02.06.2025', 'DD.MM.YYYY') AS sleeping_start,
        TO_DATE('02.10.2025', 'DD.MM.YYYY') AS sleeping_end,
        2000000000  AS sleeping_limit,
        40000000000 AS spike_limit
    FROM dual
),

-- 🔹 Bir marta SALDO'dan kerakli ma'lumotlarni olish
saldo_with_lag AS (
    SELECT
        account_code,
        oper_day,
        turnover_all_debit,
        turnover_all_credit,
        LAG(turnover_all_debit, 1, 0) OVER (
            PARTITION BY account_code
            ORDER BY oper_day
        ) AS debit_prev,
        LAG(turnover_all_credit, 1, 0) OVER (
            PARTITION BY account_code
            ORDER BY oper_day
        ) AS credit_prev
    FROM saldo
    WHERE oper_day >= ADD_MONTHS((SELECT sleeping_start FROM const), -3)
      AND oper_day <= (SELECT report_date FROM const)
),

acc AS (
    SELECT DISTINCT a.client_uid, a.code
    FROM accounts a
    JOIN client_current cc ON cc.client_uid = a.client_uid
    WHERE a.code_currency = '000'
      AND a.code_coa IN ('20208','20218','20210','20212','20214','20296')
      AND cc.subject   = 'J'
      AND cc.condition = 'A'
),

sleeping_clients AS (
    SELECT a.client_uid
    FROM acc a
    LEFT JOIN saldo_with_lag s ON s.account_code = a.code
        AND s.oper_day >= (SELECT sleeping_start FROM const)
        AND s.oper_day <= (SELECT sleeping_end FROM const)
    GROUP BY a.client_uid
    HAVING SUM(
        NVL(s.turnover_all_debit, 0) + NVL(s.turnover_all_credit, 0)
    ) <= (SELECT sleeping_limit FROM const)
),

daily_after AS (
    SELECT
        a.client_uid,
        s.oper_day AS op_date,
        SUM(
            NVL(s.turnover_all_debit, 0) - NVL(s.debit_prev, 0) +
            NVL(s.turnover_all_credit, 0) - NVL(s.credit_prev, 0)
        ) AS day_amt
    FROM acc a
    JOIN saldo_with_lag s ON s.account_code = a.code
    WHERE s.oper_day > (SELECT sleeping_end FROM const)
    GROUP BY a.client_uid, s.oper_day
    HAVING SUM(
        NVL(s.turnover_all_debit, 0) - NVL(s.debit_prev, 0) +
        NVL(s.turnover_all_credit, 0) - NVL(s.credit_prev, 0)
    ) > 0
),

spikes AS (
    SELECT
        d1.client_uid,
        d1.op_date AS day1,
        d2.op_date AS day2,
        d1.day_amt + NVL(d2.day_amt, 0) AS two_days_amt
    FROM daily_after d1
    LEFT JOIN daily_after d2 ON d2.client_uid = d1.client_uid
        AND d2.op_date = d1.op_date + 1
    WHERE d1.day_amt + NVL(d2.day_amt, 0) >= (SELECT spike_limit FROM const)
),

sleeping_balance AS (
    SELECT
        a.client_uid,
        SUM(
            NVL(s.turnover_all_debit, 0) + NVL(s.turnover_all_credit, 0)
        ) AS sleeping_90d_amt
    FROM acc a
    LEFT JOIN saldo_with_lag s ON s.account_code = a.code
        AND s.oper_day >= (SELECT sleeping_start FROM const)
        AND s.oper_day <= (SELECT sleeping_end FROM const)
    GROUP BY a.client_uid
)

SELECT
    cc.*,
    s.day1,
    s.day2,
    ROUND(s.two_days_amt / 1000000, 2) AS two_days_mln,
    s.two_days_amt AS sort_amt,
    ROUND(sb.sleeping_90d_amt / 1000000, 2) AS sleeping_90d_mln
FROM client_current cc
INNER JOIN sleeping_clients sc ON sc.client_uid = cc.client_uid
INNER JOIN spikes s ON s.client_uid = cc.client_uid
LEFT JOIN sleeping_balance sb ON sb.client_uid = cc.client_uid
ORDER BY s.two_days_amt DESC;


select *
  from mr_deals_hist t
 where t.operation_id = 1
   and t.deal_date between to_date('01.01.2026', 'dd.mm.yyyy') and
       to_date('01.01.2026', 'dd.mm.yyyy') -- 1 отправка 2 выдача
   and t.country_code in (  862,
                            492,
                            728,
                            332,
                            418,
                            524,
                            516,
                            404,
                            408,
                            887,
                            760,
                            364,
                            104,
                            092,
                            068,
                            422,
                            384,
                            024,
                            012,
                            100,
                            704,
                            120,
                            180
                         );



CREATE TABLE AML_HIGH_RISK_COUNTRIES (
    ID            NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    COUNTRY_CODE  NUMBER(3)       NOT NULL,
    COUNTRY_NAME  VARCHAR2(200)   NOT NULL,
    ALPHA2_CODE   VARCHAR2(2)     NOT NULL,
    ALPHA3_CODE   VARCHAR2(3)     NOT NULL,
    IS_ACTIVE     NUMBER(1) DEFAULT 1 NOT NULL,
    CREATED_AT    DATE DEFAULT SYSDATE NOT NULL,
    UPDATED_AT    DATE
);

SELECT * FROM AML_HIGH_RISK_COUNTRIES;

SELECT  t.DEAL_DATE,   cc.PHONE, d.city, d.country_code
FROM mr_deals_hist t join CLIENT_CURRENT cc on t.client_code = cc.code join mr_deals d on cc.code = d.client_code
WHERE t.operation_id = 1 and cc.SUBJECT = 'P'
  AND t.deal_date BETWEEN TO_DATE('01.01.2026', 'dd.mm.yyyy')
                      AND TO_DATE('28.06.2026', 'dd.mm.yyyy')
  AND EXISTS (
      SELECT 1
      FROM aml_high_risk_countries c
      WHERE c.country_code = t.country_code
        AND c.is_active = 1
  );

select country_code from mr_deals_hist;
select * from MR_DEALS;
select deal_date from mr_deals_hist order by deal_date desc;
SELECT * FROM SALDO;



drop table REPORT_207_RAW_TRANSACTIONS;

CREATE TABLE REPORT_207_RAW_TRANSACTIONS (
    OPER_DAY               DATE,
    MFO_CODE               VARCHAR2(20),

    CL_ACC                 VARCHAR2(300),
    CO_ACC                 VARCHAR2(300),
    CO_NAME                VARCHAR2(1000),
    PAY_PURPOSE            VARCHAR2(4000),

    SUM_PAY                NUMBER(20,2),

    STATE_NAME             VARCHAR2(100),

    CLIENT_CODE            VARCHAR2(50),
    CLIENT_UID             NUMBER(20),
    CLIENT_INN             VARCHAR2(50),
    CLIENT_NAME            VARCHAR2(1000),
    CLIENT_MFO             VARCHAR2(20),

    ACCOUNT_OPEN_DATE      VARCHAR2(20),

    OKED                   VARCHAR2(100),
    DIRECTOR_NAME          VARCHAR2(1000),
    DIRECTOR_PASSPORT      VARCHAR2(4000)
);

SELECT COUNT(*) FROM CR.REPORT_207_RAW_TRANSACTIONS where OPER_DAY = TO_DATE('2026-07-01', 'YYYY-MM-DD');
SELECT
    client_code,
    co_name,
    pay_purpose,
    sum_pay,
    oper_day
FROM CR.REPORT_207_RAW_TRANSACTIONS
WHERE (client_code, co_name) IN (
    SELECT
        client_code,
        co_name
    FROM CR.REPORT_207_RAW_TRANSACTIONS
    GROUP BY client_code, co_name
    HAVING COUNT(*) >= 2
)
ORDER BY
    client_code,
    co_name,
    oper_day;

select  OPER_DAY,
            MFO_CODE,
            CL_ACC,
            CO_ACC,
            CO_NAME,
            PAY_PURPOSE,
            SUM_PAY,
            STATE_NAME,
            CLIENT_CODE,
            CLIENT_UID,
            CLIENT_INN,
            CLIENT_NAME,
            CLIENT_MFO,
            ACCOUNT_OPEN_DATE,
            OKED,
            DIRECTOR_NAME,
            DIRECTOR_PASSPORT
            from REPORT_207_RAW_TRANSACTIONS;

INSERT INTO REPORT_207_RAW_TRANSACTIONS (
    OPER_DAY,
    MFO_CODE,
    CL_ACC,
    CO_ACC,
    CO_NAME,
    PAY_PURPOSE,
    SUM_PAY,
    STATE_NAME,
    CLIENT_CODE,
    CLIENT_UID,
    CLIENT_INN,
    CLIENT_NAME,
    CLIENT_MFO,
    ACCOUNT_OPEN_DATE,
    OKED,
    DIRECTOR_NAME,
    DIRECTOR_PASSPORT
)
SELECT
    TRUNC(th.curr_day),
    th.code_filial,
    l.cl_acc,
    l.co_acc,
    l.co_name,
    l.pay_purpose,
    ROUND(l.sum_pay / 100),
    'Проведен',
    cc.code,
    cc.client_uid,
    cc.inn,
    cc.name,
    cc.code_filial,
    TO_CHAR(cc.date_open, 'YYYY-MM-DD'),
    cjc.oked,
    cjc.director_name,
    cjc.director_passport
FROM transacts_history th
        JOIN leads_history l
            ON th.lead_id = l.id
        JOIN accounts a
            ON a.id = th.acc_id
        JOIN client_current cc
            ON cc.code      = SUBSTR(l.cl_acc, 17, 8)
           AND cc.condition = 'A'
           AND cc.subject   = 'J'
        LEFT JOIN client_juridical_current cjc
            ON cjc.id = cc.id
        WHERE TRUNC(th.curr_day) BETWEEN TO_DATE('2026-07-01', 'YYYY-MM-DD')
                                     AND TO_DATE('2026-07-02',   'YYYY-MM-DD')
          AND th.state_id   = 41
          AND th.op_dc      = 1
          AND a.code_coa    IN ('20208', '20210', '20212', '20214', '20216', '20218')
          AND SUBSTR(l.co_acc, 8, 5) NOT IN ('17101')
          AND SUBSTR(l.cl_acc, 17, 8) <> SUBSTR(l.co_acc, 17, 8)
        ORDER BY cc.code, l.co_name, l.pay_purpose, th.curr_day;
commit;

SELECT SYS_CONTEXT('USERENV', 'DB_NAME'),
                              SYS_CONTEXT('USERENV', 'SERVICE_NAME'),
                              SYS_CONTEXT('USERENV', 'INSTANCE_NAME')
                       FROM dual;

BEGIN
    DBMS_SCHEDULER.CREATE_JOB (
        job_name        => 'JOB_REPORT_207_RAW_TRANSACTIONS',

        job_type        => 'PLSQL_BLOCK',

        job_action      => q'[
BEGIN

    INSERT INTO REPORT_207_RAW_TRANSACTIONS (
        OPER_DAY,
        MFO_CODE,
        CL_ACC,
        CO_ACC,
        CO_NAME,
        PAY_PURPOSE,
        SUM_PAY,
        STATE_NAME,
        CLIENT_CODE,
        CLIENT_UID,
        CLIENT_INN,
        CLIENT_NAME,
        CLIENT_MFO,
        ACCOUNT_OPEN_DATE,
        OKED,
        DIRECTOR_NAME,
        DIRECTOR_PASSPORT
    )
    SELECT
        TRUNC(th.curr_day),
        th.code_filial,
        l.cl_acc,
        l.co_acc,
        l.co_name,
        l.pay_purpose,
        ROUND(l.sum_pay / 100),
        'Проведен',
        cc.code,
        cc.client_uid,
        cc.inn,
        cc.name,
        cc.code_filial,
        TO_CHAR(cc.date_open, 'YYYY-MM-DD'),
        cjc.oked,
        cjc.director_name,
        cjc.director_passport
    FROM transacts_history th
        JOIN leads_history l
            ON th.lead_id = l.id
        JOIN accounts a
            ON a.id = th.acc_id
        JOIN client_current cc
            ON cc.code = SUBSTR(l.cl_acc, 17, 8)
           AND cc.condition = 'A'
           AND cc.subject = 'J'
        LEFT JOIN client_juridical_current cjc
            ON cjc.id = cc.id
    WHERE TRUNC(th.curr_day) = TRUNC(SYSDATE)
      AND th.state_id = 41
      AND th.op_dc = 1
      AND a.code_coa IN ('20208','20210','20212','20214','20216','20218')
      AND SUBSTR(l.co_acc,8,5) <> '17101'
      AND SUBSTR(l.cl_acc,17,8) <> SUBSTR(l.co_acc,17,8);

    COMMIT;

END;
]',

        start_date      => TO_TIMESTAMP_TZ(
                               '2026-07-02 23:00:00 Asia/Tashkent',
                               'YYYY-MM-DD HH24:MI:SS TZR'
                           ),

        repeat_interval => 'FREQ=DAILY;BYHOUR=23;BYMINUTE=0;BYSECOND=0',

        enabled         => TRUE,

        comments        => 'Refresh REPORT_207_RAW_TRANSACTIONS every day at 23:00 Asia/Tashkent'
    );
END;


SELECT
    job_name,
    enabled,
    state,
    repeat_interval,
    start_date,
    last_start_date,
    next_run_date
FROM USER_SCHEDULER_JOBS
ORDER BY job_name;


SELECT * FROM CR.REPORT_207_RAW_TRANSACTIONS where OPER_DAY = TO_DATE('2026-07-02', 'YYYY-MM-DD') and CLIENT_NAME = 'BAXTIYOR-RAMIL TRANS MCHJ';
SELECT count(*) FROM CR.REPORT_207_RAW_TRANSACTIONS where OPER_DAY = TO_DATE('2026-07-2', 'YYYY-MM-DD');

SELECT CLIENT_CODE
FROM CR.REPORT_207_RAW_TRANSACTIONS
WHERE OPER_DAY = TO_DATE('2026-07-02', 'YYYY-MM-DD')
GROUP BY CLIENT_CODE
HAVING COUNT(*) > 1
   AND SUM(SUM_PAY) >= 412000000;

SELECT
            cjc.director_passport,
            (SELECT cpc.passport_type
             FROM   client_physical_current cpc
             WHERE  TRIM(cpc.passport_serial || cpc.passport_number)
                        = TRIM(cjc.director_passport)
             AND    ROWNUM = 1)                              AS director_doc_type,
            (SELECT TO_CHAR(cpc.birthday, 'YYYY-MM-DD')
             FROM   client_physical_current cpc
             WHERE  TRIM(cpc.passport_serial || cpc.passport_number)
                        = TRIM(cjc.director_passport)
             AND    ROWNUM = 1)                              AS director_birthday,
            (SELECT TO_CHAR(cpc.pinfl)
             FROM   client_physical_current cpc
             WHERE  TRIM(cpc.passport_serial || cpc.passport_number)
                        = TRIM(cjc.director_passport)
             AND    ROWNUM = 1)                              AS director_pinfl,
            (SELECT cpc.passport_serial
             FROM   client_physical_current cpc
             WHERE  TRIM(cpc.passport_serial || cpc.passport_number)
                        = TRIM(cjc.director_passport)
             AND    ROWNUM = 1)                              AS director_passport_serial,
            (SELECT cpc.passport_number
             FROM   client_physical_current cpc
             WHERE  TRIM(cpc.passport_serial || cpc.passport_number)
                        = TRIM(cjc.director_passport)
             AND    ROWNUM = 1)                              AS director_passport_number,
            (SELECT TO_CHAR(cpc.passport_registration_date, 'YYYY-MM-DD')
             FROM   client_physical_current cpc
             WHERE  TRIM(cpc.passport_serial || cpc.passport_number)
                        = TRIM(cjc.director_passport)
             AND    ROWNUM = 1)                              AS director_doc_issue_date

        from client_juridical_current cjc
        WHERE cjc.code IN ('05019563');

select * from CLIENT_CURRENT where code = '04916582';
select * from CLIENT_CURRENT where CLIENT_UID = '164369';


select * from ibs.ln_avp_debtor_sources@iabs;
select count(*) from ibs.ln_avp_debtor_sources@iabs;

SELECT COUNT(*) FROM ibs.saldo@iabs where oper_day >= TO_DATE('2026-07-02', 'YYYY-MM-DD');
SELECT COUNT(*) FROM ibs.saldo@iabs where oper_day = TO_DATE('2026-07-02', 'YYYY-MM-DD');



WITH qualifying_clients AS (
    SELECT r.client_code
    FROM   CR.REPORT_207_RAW_TRANSACTIONS r
    WHERE  r.oper_day BETWEEN TO_DATE('2026-07-03', 'YYYY-MM-DD')
                          AND TO_DATE('2026-07-03',   'YYYY-MM-DD')
    GROUP BY r.client_code
    HAVING COUNT(*) >= 2
       AND SUM(r.sum_pay) >= 412000000
)
SELECT
    r.OPER_DAY,
    r.MFO_CODE,
    r.CL_ACC,
    r.CO_ACC,
    r.CO_NAME,
    r.PAY_PURPOSE,
    r.SUM_PAY,
    r.STATE_NAME,
    r.CLIENT_CODE,
    r.CLIENT_UID,
    r.CLIENT_INN,
    r.CLIENT_NAME,
    r.CLIENT_MFO,
    r.ACCOUNT_OPEN_DATE,
    r.OKED,
    r.DIRECTOR_NAME,
    r.DIRECTOR_PASSPORT
FROM   CR.REPORT_207_RAW_TRANSACTIONS r
JOIN   qualifying_clients q ON q.client_code = r.client_code
WHERE  r.oper_day BETWEEN TO_DATE('2026-07-03', 'YYYY-MM-DD')
                      AND TO_DATE('2026-07-03',   'YYYY-MM-DD')
ORDER BY r.client_code;


select count(*) from client_physical_current;
select count(*) from client_current;
select count(*) from client_juridical_current;
select * from accounts;
SELECT
    client_uid,
    COUNT(*) cnt,
    SUM(turnover_all_debit) debit,
    SUM(turnover_all_credit) credit
FROM accounts
WHERE client_uid = 16757815
GROUP BY client_uid;


SELECT
    code,
    turnover_all_debit,
    turnover_all_credit
FROM accounts
WHERE client_uid = 16757815;


SELECT
    SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA') AS current_schema,
    d.name AS database_name,
    i.instance_name,
    SYS_CONTEXT('USERENV', 'DB_NAME') AS db_name
FROM v$database d
CROSS JOIN v$instance i;

select * from REPORT_207_RAW_TRANSACTIONS order by OPER_DAY = '2026-07-06';


SELECT
    SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA') AS current_schema,
    d.name AS database_name,
    i.instance_name,
    SYS_CONTEXT('USERENV', 'DB_NAME') AS db_name,

    c.owner,
    c.table_name,
    c.column_id,
    c.column_name,
    c.data_type,
    c.data_length,
    c.nullable,
    c.data_default
FROM all_tab_columns c
CROSS JOIN v$database d
CROSS JOIN v$instance i
WHERE c.owner = 'CR'
  AND c.table_name = 'REPORT_207_RAW_TRANSACTIONS'
ORDER BY c.column_id;


WITH director_physical AS (
    SELECT
        cpc.passport_serial,
        cpc.passport_number,
        cpc.passport_type,
        cpc.birthday,
        cpc.pinfl,
        cpc.passport_registration_date,
        ROW_NUMBER() OVER (
            PARTITION BY TRIM(cpc.passport_serial || cpc.passport_number)
            ORDER BY cpc.id
        ) AS rn
    FROM client_physical_current cpc
)
SELECT
    cc.code                                               AS client_code,
    cjc.director_passport,
    dp.passport_type                                     AS director_doc_type,
    TO_CHAR(dp.birthday, 'YYYY-MM-DD')                    AS director_birthday,
    TO_CHAR(dp.pinfl)                                     AS director_pinfl,
    dp.passport_serial                                    AS director_passport_serial,
    dp.passport_number                                    AS director_passport_number,
    TO_CHAR(dp.passport_registration_date, 'YYYY-MM-DD')  AS director_doc_issue_date
FROM client_current cc
JOIN client_juridical_current cjc ON cjc.id = cc.id
LEFT JOIN director_physical dp
       ON dp.rn = 1
      AND TRIM(dp.passport_serial || dp.passport_number) = TRIM(cjc.director_passport)
WHERE cc.code IN ('00011999'
);


SELECT
        TRUNC(th.curr_day),
        th.code_filial,
        l.cl_acc,
        l.co_acc,
        l.co_name,
        l.pay_purpose,
        ROUND(l.sum_pay / 100),
        'Проведен',
        cc.code,
        cc.client_uid,
        cc.inn,
        cc.name,
        cc.code_filial,
        TO_CHAR(cc.date_open, 'YYYY-MM-DD'),
        cjc.oked,
        cjc.director_name,
        cjc.director_passport
    FROM transacts_history th
        JOIN leads_history l
            ON th.lead_id = l.id
        JOIN accounts a
            ON a.id = th.acc_id
        JOIN client_current cc
            ON cc.code = SUBSTR(l.cl_acc, 17, 8)
           AND cc.condition = 'A'
           AND cc.subject = 'J'
        LEFT JOIN client_juridical_current cjc
            ON cjc.id = cc.id
    WHERE TRUNC(th.curr_day) = TO_DATE('2026-07-06', 'YYYY-MM-DD')
      AND th.state_id = 41
      AND th.op_dc = 1
      AND a.code_coa IN ('20208','20210','20212','20214','20216','20218')
      AND SUBSTR(l.co_acc,8,5) <> '17101'
      AND SUBSTR(l.cl_acc,17,8) <> SUBSTR(l.co_acc,17,8);


SELECT
        *
        FROM transacts_history th
        JOIN leads_history l
            ON  l.id = th.lead_id
        JOIN accounts a
            ON  a.id             = th.acc_id
            AND a.code_coa       IN ('20208', '20218', '20210', '20214')
            AND a.code_currency  = '000'
            AND a.client_uid     IN ('00938218', '07455311', '07263793', '07462591', '02242550')
        WHERE th.curr_day  BETWEEN TO_DATE('2026-07-16', 'YYYY-MM-DD')
                               AND TO_DATE('2026-07-16',   'YYYY-MM-DD')
          AND th.state_id  = 41
        ORDER BY a.client_uid, th.curr_day


clickhouse 

SELECT COUNT(*) FROM saldo WHERE OPER_DAY = '2026-07-04';


SELECT count(*) FROM saldo WHERE OPER_DAY >= '2026-01-01'; --102632751

WITH
    toDate('2026-01-01') AS start_date,
    today() AS end_date

SELECT calendar.day AS missing_day
FROM
(
    SELECT start_date + number AS day
    FROM numbers(dateDiff('day', start_date, end_date) + 1)
    WHERE toDayOfWeek(start_date + number) BETWEEN 1 AND 5
) AS calendar
LEFT JOIN
(
    SELECT DISTINCT OPER_DAY
    FROM saldo
    WHERE OPER_DAY BETWEEN start_date AND end_date
) AS t
ON calendar.day = t.OPER_DAY
WHERE t.OPER_DAY IS NULL
ORDER BY calendar.day;


OPTIMIZE TABLE saldo PARTITION '202602' FINAL;
OPTIMIZE TABLE saldo PARTITION '202603' FINAL;
OPTIMIZE TABLE saldo PARTITION '202604' FINAL;
OPTIMIZE TABLE saldo PARTITION '202605' FINAL;
OPTIMIZE TABLE saldo PARTITION '202606' FINAL;
OPTIMIZE TABLE saldo PARTITION '202607' FINAL;

SELECT
    query_id,
    elapsed,
    query
FROM system.processes
WHERE user = 'dwh_admin';

select * from mob_users;
show create table mob_users;