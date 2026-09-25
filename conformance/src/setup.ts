import { afterAll } from 'vitest';
import { closeSharedGoClient } from './adapter';

afterAll(async () => {
  await closeSharedGoClient();
});
