import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { afterEach, describe, expect, test } from 'vitest';
import { GoProcessClient, GoProcessError } from './client';

const here = path.dirname(fileURLToPath(import.meta.url));
const clients: GoProcessClient[] = [];

afterEach(async () => {
  await Promise.all(clients.splice(0).map((client) => client.close()));
});

function fakeClient(): GoProcessClient {
  const client = new GoProcessClient({
    command: process.execPath,
    args: [path.join(here, '../testdata/fake-cli.mjs')],
  });
  clients.push(client);
  return client;
}

describe('GoProcessClient', () => {
  test('correlates JSON-lines responses', async () => {
    const client = fakeClient();
    const [first, second] = await Promise.all([
      client.request<{ value: string }>('echo', { value: 'first' }),
      client.request<{ value: string }>('echo', { value: 'second' }),
    ]);
    expect(first.value).toBe('first');
    expect(second.value).toBe('second');
  });

  test('surfaces protocol error and code', async () => {
    const client = fakeClient();
    expect(() => client.requestSync('fail')).toThrowError(
      expect.objectContaining<Partial<GoProcessError>>({
      message: 'deliberate failure',
      code: 'FAKE',
      }),
    );
  });
});
