/* Return active projects
   and their owners. */
SELECT project, owner
FROM projects
WHERE enabled = TRUE AND owner <> 'nobody';
