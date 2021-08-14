import equal from "fast-deep-equal/es6";
import {
  get,
  Subscriber,
  Unsubscriber,
  Updater,
  writable,
  Writable,
} from "svelte/store";

export class apiResource<T> {
  url: string;
  listener: (v: T) => void;

  constructor(url: string, listener?: (v: T) => void) {
    this.url = url;
    if (listener) {
      this.listener = listener;
    } else {
      this.listener = () => {};
    }
  }

  get(): Promise<T> {
    const p = fetch(this.url).then((r) => r.json() as Promise<T>);
    p.then((v) => this.listener(v));
    return p;
  }

  put(v: T): Promise<T> {
    return fetch(this.url, {
      method: "PUT",
      body: JSON.stringify(v),
      headers: { "Content-Type": "application/json" },
    }).then((r) => {
      const contentType = r.headers.get("content-type");
      if (contentType && contentType.indexOf("application/json") >= 0) {
        const p = r.json() as Promise<T>;
        p.then((v) => this.listener(v));
        return p;
      } else {
        return this.get();
      }
    });
  }

  remove(): Promise<void> {
    return fetch(this.url, { method: "DELETE" }).then(() =>
      this.listener(undefined)
    );
  }
}

export enum LoadingState {
  Uninitialized,
  Loading,
  Loaded,
}

export class resourceWritable<T> extends apiResource<T> implements Writable<T> {
  public state: LoadingState;
  private _w: Writable<T>;

  constructor(url: string, listener?: (v: T) => void, w?: Writable<T>) {
    super(url, listener);
    if (w) {
      this._w = w;
    } else {
      this._w = writable(undefined, () => {
        this.reload();
      });
    }
    this.state = LoadingState.Uninitialized;
  }

  get(): Promise<T> {
    if (this.state == LoadingState.Uninitialized) {
      this.state = LoadingState.Loading;
    }
    return super.get().then((v) => {
      this._w.set(v);
      this.state = LoadingState.Loaded;
      return v;
    });
  }

  put(v: T): Promise<T> {
    this._w.set(v);
    return super.put(v).then((v2) => {
      if (!equal(get(this._w), v2)) {
        this._w.set(v2);
      }
      return v2;
    });
  }

  remove(): Promise<void> {
    this._w.set(undefined);
    return super.remove();
  }

  reload(): Promise<void> {
    return this.get().then(() => {});
  }

  subscribe(
    run: Subscriber<T>,
    invalidate?: (value?: T) => void
  ): Unsubscriber {
    return this._w.subscribe(run, invalidate);
  }

  set(value: T): void {
    this.put(value);
  }

  update(updater: Updater<T>): void {
    const oldV = get(this._w);
    const newV = updater(oldV);
    if (!equal(oldV, newV)) {
      this.set(newV);
    }
  }
}
