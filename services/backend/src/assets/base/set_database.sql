{# 
-- Assigns the database to query
-- AllowedDatabases: ['postgres', 'sqlite', 'mssql']
-- Created By: copito
#}

{%- if database_type = 'postgres' %}

-- Query running using postgres
USE {{database}};

{%- elif database_type = 'mssql' %}

-- Query running using mssql
USE {{database}};

{%- elif database_type = 'sqlite' %}

-- Query running using sqlite
USE {{database}};

{%- else %}

-- Query using `{{database_type}}` type
USE {{database}};

{%- endif %}