ALTER TABLE tasks ADD COLUMN jira_id text;

UPDATE tasks SET
    jira_id = substring(description from '\[(.*?)\]'),
    description = regexp_replace(description, '\[(.*?)\]', '', 'g');
