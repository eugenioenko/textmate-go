import { spawn, type ChildProcessWithoutNullStreams } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

interface WireResponse {
  id: number;
  error?: string;
  code?: string;
  [key: string]: unknown;
}

export class GoProcessError extends Error {
  constructor(message: string, readonly code?: string) {
    super(message);
    this.name = 'GoProcessError';
  }
}

export interface GoProcessClientOptions {
  command?: string;
  args?: string[];
  cwd?: string;
  env?: NodeJS.ProcessEnv;
}

export class GoProcessClient {
  private child: ChildProcessWithoutNullStreams | null = null;
  private nextId = 1;
  private stdoutBuffer = '';

  constructor(private readonly options: GoProcessClientOptions = {}) {}

  async request<T extends object>(op: string, payload: object = {}): Promise<T> {
    return this.requestSync(op, payload);
  }

  requestSync<T extends object>(op: string, payload: object = {}): T {
    this.start();
    const id = this.nextId++;
    writeAll(streamFd(this.child!.stdin), `${JSON.stringify({ id, op, ...payload })}\n`);
    const value = this.readResponse(id);
    if (value.error) throw new GoProcessError(value.error, value.code);
    return value as T;
  }

  async close(): Promise<void> {
    const child = this.child;
    this.child = null;
    if (!child) return;
    child.stdin.end();
    if (child.exitCode === null) child.kill();
  }

  private start(): void {
    if (this.child) return;
    const here = path.dirname(fileURLToPath(import.meta.url));
    const repository = path.resolve(here, '../..');
    const command = this.options.command ?? process.env.TEXTMATE_GO_CLI ?? 'go';
    const args =
      this.options.args ??
      parseArgs(process.env.TEXTMATE_GO_CLI_ARGS) ??
      (process.env.TEXTMATE_GO_CLI ? [] : ['run', './cmd/textmate-cli']);
    const child = spawn(command, args, {
      cwd: this.options.cwd ?? repository,
      env: { ...process.env, ...this.options.env },
      stdio: ['pipe', 'pipe', 'pipe'],
    });
    this.child = child;

    child.stderr.on('data', (chunk: Buffer) => process.stderr.write(chunk));
    child.once('exit', (code, signal) => {
      if (this.child === child) this.child = null;
    });
  }

  private readResponse(expectedId: number): WireResponse {
    for (;;) {
      const newline = this.stdoutBuffer.indexOf('\n');
      if (newline >= 0) {
        const line = this.stdoutBuffer.slice(0, newline);
        this.stdoutBuffer = this.stdoutBuffer.slice(newline + 1);
        let response: WireResponse;
        try {
          response = JSON.parse(line) as WireResponse;
        } catch (error) {
          throw new Error(`invalid JSON from textmate-cli: ${line}`, { cause: error });
        }
        if (response.id !== expectedId) {
          throw new Error(`textmate-cli replied with id ${response.id}, expected ${expectedId}`);
        }
        return response;
      }
      const buffer = Buffer.allocUnsafe(16 * 1024);
      let count: number;
      try {
        count = fs.readSync(streamFd(this.child!.stdout), buffer, 0, buffer.length, null);
      } catch (error) {
        if (isAgain(error)) {
          shortWait();
          continue;
        }
        throw error;
      }
      if (count === 0) throw new Error('textmate-cli exited before replying');
      this.stdoutBuffer += buffer.toString('utf8', 0, count);
    }
  }
}

function writeAll(fd: number, value: string): void {
  const buffer = Buffer.from(value);
  let offset = 0;
  while (offset < buffer.length) {
    try {
      offset += fs.writeSync(fd, buffer, offset, buffer.length - offset);
    } catch (error) {
      if (!isAgain(error)) throw error;
      shortWait();
    }
  }
}

function isAgain(error: unknown): boolean {
  return (error as NodeJS.ErrnoException).code === 'EAGAIN';
}

const sleepSignal = new Int32Array(new SharedArrayBuffer(4));
function shortWait(): void {
  Atomics.wait(sleepSignal, 0, 0, 1);
}

function streamFd(stream: NodeJS.ReadableStream | NodeJS.WritableStream): number {
  const fd = (stream as unknown as { _handle?: { fd?: number } })._handle?.fd;
  if (typeof fd !== 'number') throw new Error('could not access textmate-cli pipe');
  return fd;
}

function parseArgs(value: string | undefined): string[] | undefined {
  if (!value) return undefined;
  const parsed: unknown = JSON.parse(value);
  if (!Array.isArray(parsed) || !parsed.every((part) => typeof part === 'string')) {
    throw new Error('TEXTMATE_GO_CLI_ARGS must be a JSON array of strings');
  }
  return parsed;
}
