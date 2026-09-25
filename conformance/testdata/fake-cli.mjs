import readline from 'node:readline';

const input = readline.createInterface({ input: process.stdin });
input.on('line', (line) => {
  const request = JSON.parse(line);
  const response =
    request.op === 'fail'
      ? { id: request.id, error: 'deliberate failure', code: 'FAKE' }
      : { id: request.id, value: request.value };
  process.stdout.write(`${JSON.stringify(response)}\n`);
});
