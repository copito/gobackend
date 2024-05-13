{# 
-- Get all tables
-- AllowedDatabases: ['postgres', 'sqlite', 'mssql']
-- Created By: copito
#}

{%- if database_type = 'postgres' %}

-- Query running using postgres
SELECT * FROM pg_catalog.pg_tables;


{%- elif database_type = 'sqlite' %}

-- Query running using sqlite
.tables


{%- elif database_type = 'mssql' %}

-- Query running using mssql
SELECT
    TABLE_NAME
FROM {{database}}.INFORMATION_SCHEMA.TABLES
WHERE 
    TABLE_TYPE = 'BASE TABLE';

{%- else %}

-- Query using `{{database_type}}` type

{%- endif %}