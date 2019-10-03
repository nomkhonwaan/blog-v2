import { Injectable } from '@angular/core';

import { StorageModule } from './storage.module';

@Injectable({
  providedIn: StorageModule,
})
export class LocalStorageService {
  /**
   * Empties the localStorage.
   */
  clear(): void {
    localStorage.clear();
  }

  /**
   * Removes value by its key from the localStorage API.
   *
   * @param key string
   */
  remove(key: string): void {
    localStorage.removeItem(key);
  }

  /**
   * Retrieves value by its key from the localStorage API.
   *
   * @param key string
   */
  get(key: string): string {
    return localStorage.getItem(key);
  }

  /**
   * Likes get() but will return as a number.
   *
   * @param key string
   */
  getNumber(key: string): number {
    return parseInt(this.get(key), 10);
  }

  /**
   * Creates or updates specific value along with the key to the localStorage API.
   *
   * @param key string
   * @param value string
   */
  set(key: string, value: string): void {
    localStorage.setItem(key, value);
  }

  /**
   * Performs multiple records creating or updating.
   *
   * @param keyValuePair object
   */
  setAll(keyValuePair: { [name: string]: string }): void {
    Object.keys(keyValuePair).map((key) => {
      this.set(key, keyValuePair[key]);
    });
  }
}
