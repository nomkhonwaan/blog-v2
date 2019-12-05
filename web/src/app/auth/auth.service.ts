import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Store } from '@ngrx/store';
import { WebAuth } from 'auth0-js';

import { AuthModule } from './auth.module';

import { LocalStorageService } from '../storage/local-storage.service';
import { setAuthentication } from '../app.actions';

@Injectable({
  providedIn: AuthModule,
  deps: [LocalStorageService, Router, Store, WebAuth],
})
export class AuthService {
  private accessToken?: string;
  private idToken?: string;
  private expiresAt?: number;

  constructor(
    private localStorage: LocalStorageService,
    private router: Router,
    private store: Store<AppState>,
    private webAuth: WebAuth,
  ) {
    this.accessToken = this.localStorage.get('accessToken');
    this.idToken = this.localStorage.get('idToken');
    this.expiresAt = this.localStorage.getNumber('expiresAt');

    this.dispatchNgrxStore();
  }

  /**
   * Dispatch the @ngrx/store for updating `accessToken` and `idToken` values
   */
  dispatchNgrxStore(): void {
    if (this.isAuthenticated()) {
      this.store.dispatch(setAuthentication({ accessToken: this.accessToken, idToken: this.idToken }));
    }
  }

  /**
   * Redirect to the Auth0 login page.
   *
   * @param redirectPath string
   */
  login(redirectPath?: string): void {
    if (redirectPath !== null) {
      this.localStorage.set('redirectPath', redirectPath);
    }

    this.webAuth.authorize();
  }

  /**
   * Parse the authentication result from URL hash.
   */
  handleAuthentication(): void {
    this.webAuth.parseHash((_: Error, authResult: AuthResult) => {
      if (authResult && authResult.accessToken && authResult.idToken) {
        this.storeCredentials(authResult);

        // Redirect back to the previous path (that was saved in the login step) or home path
        const redirectPath: string = this.localStorage.get('redirectPath');
        if (redirectPath !== null) {
          this.localStorage.remove('redirectPath');
        }

        this.router.navigate([redirectPath || '/']);
      }
    });
  }

  /**
   * Store the authentication result in class properties.
   *
   * @param authResult object An authentication result which contains accessToken, idToken and expiresAt
   */
  storeCredentials(authResult): void {
    this.accessToken = authResult.accessToken;
    this.idToken = authResult.idToken;
    this.expiresAt = authResult.expiresIn * 1000 + Date.now();

    this.localStorage.setAll({
      accessToken: this.accessToken,
      idToken: this.idToken,
      expiresAt: this.expiresAt.toString()
    });

    this.dispatchNgrxStore();
  }

  /**
   * Perform silent authentication to renew the session.
   */
  renewTokens(): void {
    this.webAuth.checkSession({}, (_: Error, authResult: AuthResult) => {
      if (authResult && authResult.accessToken && authResult.idToken) {
        this.storeCredentials(authResult);
      }
    });
  }

  /**
   * Remove the user's tokens and expiry time from class properties.
   */
  logout(): void {
    this.accessToken = null;
    this.idToken = null;
    this.expiresAt = null;

    this.localStorage.clear();
  }

  /**
   * Check whether the user's Access Token is set and its expiry time has passed.
   */
  isAuthenticated(): boolean {
    return this.accessToken !== '' && this.idToken !== '' && Date.now() < this.expiresAt;
  }
}
