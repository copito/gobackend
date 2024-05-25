{# 
-- Metric Level: Table
-- Metric Name: Row Count
-- ID: table.row_count
-- AllowedDatabases: ['postgres', 'sqlite', 'mssql']
-- Created By: copito
#}

{%- if database_type = 'postgres' %}

-- Query running using postgres
SELECT 
    COUNT(*) as value
FROM {{table}} 
LIMIT 1

{%- elif database_type = 'sqlite' %}

-- Query running using sqlite
SELECT 
    COUNT(*) as value
FROM {{table}} 
LIMIT 1

{%- elif database_type = 'mssql' %}

-- Query running using sqlite
SELECT 
    COUNT(*) as value
FROM {{table}} 
LIMIT 1

{%- else %}

-- Query using `{{database_type}}` type
SELECT 0 as value

{%- endif %}