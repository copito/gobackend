{# 
-- Metric Level: Column
-- Metric Name: Value Match Regex
-- ID: column.row_count
-- AllowedDatabases: ['postgres', 'sqlite', 'mssql']
-- Created By: copito
#}

{%- if database_type = 'postgres' %}

-- Query running using postgres
SELECT 
    CASE 
        WHEN {{column}} ~ '{{regex}}' THEN 1 
    ELSE 0 
    END as value
FROM {{table}} 
LIMIT 1

{%- elif database_type = 'sqlite' %}

-- Query running using sqlite
SELECT 
    CASE 
        WHEN {{column}} REGEXP '{{regex}}' THEN 1 
    ELSE 0 
    END as value
FROM {{table}} 
LIMIT 1

{%- elif database_type = 'mssql' %}

-- Query running using mssql
SELECT 
    CASE 
        WHEN {{column}} LIKE '{{regex}}' THEN 1 
    ELSE 0 
    END as value
FROM {{table}} 
LIMIT 1

{%- else %}

-- Query using `{{database_type}}` type
SELECT 0

{%- endif %}