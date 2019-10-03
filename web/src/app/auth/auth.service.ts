import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { WebAuth } from 'auth0-js';

import { AuthModule } from './auth.module';

import { LocalStorageService } from '../storage/local-storage.service';

@Injectable({
  providedIn: AuthModule,
  deps: [LocalStorageService, Router, WebAuth],
})
export class AuthService {
  private accessToken?: string;
  private idToken?: string;
  private expiresAt?: number;

  constructor(
    private localStorage: LocalStorageService,
    private router: Router,
    private webAuth: WebAuth,
  ) {
    this.accessToken = this.localStorage.get('accessToken');
    this.idToken = this.localStorage.get('idToken');
    this.expiresAt = this.localStorage.getNumber('expiresAt');
  }

  /**
   * Redirects to the Auth0 login page.
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
   * Parses the authentication result from URL hash.
   */
  handleAuthentication(): void {
    this.webAuth.parseHash((_: Error, authResult: AuthResult) => {
      if (authResult && authResult.accessToken && authResult.idToken) {
        this.localLogin(authResult);

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
   * Stores the authentication result in class properties.
   *
   * @param authResult object An authentication result which contains accessToken, idToken and expiresAt
   */
  localLogin(authResult): void {
    this.accessToken = authResult.accessToken;
    this.idToken = authResult.idToken;
    this.expiresAt = authResult.expiresIn * 1000 + Date.now();

    this.localStorage.setAll({
      accessToken: this.accessToken,
      idToken: this.idToken,
      expiresAt: this.expiresAt.toString()
    });
  }

  /**
   * Performs silent authentication to renew the session.
   */
  renewTokens(): void {
    this.webAuth.checkSession({}, (_, authResult) => {
      if (authResult && authResult.accessToken && authResult.idToken) {
        this.localLogin(authResult);
      }
    });
  }

  /**
   * Removes the user's tokens and expiry time from class properties.
   */
  logout(): void {
    this.accessToken = null;
    this.idToken = null;
    this.expiresAt = null;

    this.localStorage.clear();
  }

  /**
   * Checks whether the user's Access Token is set and its expiry time has passed.
   */
  isAuthenticated(): boolean {
    return this.accessToken && Date.now() < this.expiresAt;
  }
}
