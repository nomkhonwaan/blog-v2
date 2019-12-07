import { Injectable } from '@angular/core';
import { CanActivate, ActivatedRouteSnapshot, Router, RouterStateSnapshot } from '@angular/router';

import { AuthModule } from './auth.module';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: AuthModule,
  deps: [AuthService],
})
export class AuthGuard implements CanActivate {

  constructor(private auth: AuthService, private router: Router) { }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean {
    if (this.auth.isAuthenticated()) {
      return true;
    }

    this.auth.login(this.router.url);
  }

}
