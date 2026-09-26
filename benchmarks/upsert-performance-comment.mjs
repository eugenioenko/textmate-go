#!/usr/bin/env node

import fs from 'node:fs';

const marker = '<!-- textmate-go-performance -->';
const token = process.env.GITHUB_TOKEN;
const repository = process.env.GITHUB_REPOSITORY;
const pullNumber = process.env.PR_NUMBER;
const reportFile = process.env.REPORT_FILE;

if (!token || !repository || !pullNumber || !reportFile) {
  throw new Error('GITHUB_TOKEN, GITHUB_REPOSITORY, PR_NUMBER, and REPORT_FILE are required');
}
const body = fs.readFileSync(reportFile, 'utf8');
if (!body.includes(marker)) throw new Error(`performance report is missing marker ${marker}`);

const api = async (route, options = {}) => {
  const response = await fetch(`https://api.github.com${route}`, {
    ...options,
    headers: {
      Accept: 'application/vnd.github+json',
      Authorization: `Bearer ${token}`,
      'X-GitHub-Api-Version': '2022-11-28',
      ...options.headers,
    },
  });
  if (!response.ok) throw new Error(`${options.method ?? 'GET'} ${route}: ${response.status} ${await response.text()}`);
  return response.status === 204 ? null : response.json();
};

let existing;
for (let page = 1; !existing; page++) {
  const comments = await api(`/repos/${repository}/issues/${pullNumber}/comments?per_page=100&page=${page}`);
  existing = comments.find((comment) => comment.user?.type === 'Bot' && comment.body?.includes(marker));
  if (comments.length < 100) break;
}

if (existing) {
  await api(`/repos/${repository}/issues/comments/${existing.id}`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ body }),
  });
  console.log(`updated performance comment ${existing.id}`);
} else {
  const comment = await api(`/repos/${repository}/issues/${pullNumber}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ body }),
  });
  console.log(`created performance comment ${comment.id}`);
}
