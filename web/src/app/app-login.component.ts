import { Component, OnInit } from '@angular/core';

import { AuthService } from './auth/auth.service';

@Component({
  selector: 'app-login',
  template: ``,
})
export class AppLoginComponent implements OnInit {

  constructor(private auth: AuthService) { }

  ngOnInit(): void {
    this.auth.handleAuthentication();
  }

}
