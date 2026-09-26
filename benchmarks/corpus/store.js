// Minimal reactive store with selectors and middleware.
import { EventEmitter } from 'node:events';

const DEFAULT_OPTIONS = Object.freeze({
  strict: true,
  history: 50,
  name: 'store',
});

export class Store extends EventEmitter {
  #state;
  #middleware = [];
  #history = [];

  constructor(initialState = {}, options = {}) {
    super();
    this.options = { ...DEFAULT_OPTIONS, ...options };
    this.#state = structuredClone(initialState);
  }

  get state() {
    return this.options.strict ? Object.freeze({ ...this.#state }) : this.#state;
  }

  use(fn) {
    if (typeof fn !== 'function') {
      throw new TypeError(`middleware must be a function, got ${typeof fn}`);
    }
    this.#middleware.push(fn);
    return () => {
      this.#middleware = this.#middleware.filter((m) => m !== fn);
    };
  }

  async dispatch(action) {
    const chain = this.#middleware.reduceRight(
      (next, mw) => (act) => mw(this, act, next),
      async (act) => this.#reduce(act),
    );
    return chain(action);
  }

  #reduce({ type, payload = null }) {
    const previous = this.#state;
    switch (type) {
      case 'set':
        this.#state = { ...previous, [payload.key]: payload.value };
        break;
      case 'remove': {
        const { [payload.key]: _removed, ...rest } = previous;
        this.#state = rest;
        break;
      }
      default:
        return previous;
    }
    this.#history.push(previous);
    if (this.#history.length > this.options.history) this.#history.shift();
    this.emit('change', this.#state, previous);
    return this.#state;
  }

  select(selector, onChange) {
    let last = selector(this.#state);
    const listener = (next) => {
      const value = selector(next);
      if (!Object.is(value, last)) onChange((last = value));
    };
    this.on('change', listener);
    return () => this.off('change', listener);
  }
}

export const logger = (store, action, next) => {
  const started = performance.now();
  const result = next(action);
  console.debug(`[${store.options.name}] ${action.type} in ${(performance.now() - started).toFixed(2)}ms`);
  return result;
};
