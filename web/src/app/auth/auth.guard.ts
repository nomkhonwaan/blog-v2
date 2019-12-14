import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';

import { AuthModule } from './auth.module';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: AuthModule,
  deps: [AuthService],
})
export class AuthGuard implements CanActivate {

  constructor(private auth: AuthService) { }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
    if (this.auth.isLoggedIn()) {
      return true;
    }

    this.auth.login(state.url);
  }

}
