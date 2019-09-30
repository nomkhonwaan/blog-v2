import { NgModule } from '@angular/core';
import { WebAuth } from 'auth0-js';

import { AuthService } from './auth.service';

import { environment } from '../../environments/environment';
import { LocalStorageModule } from '../storage/local-storage.module';

@NgModule({
  imports: [
    LocalStorageModule,
  ],
  providers: [
    {
      provide: AuthService,
      deps: [WebAuth],
    },
    {
      provide: WebAuth,
      useFactory: () => new WebAuth({
        clientID: environment.auth0.clientId,
        domain: 'nomkhonwaan.auth0.com',
        responseType: 'token id_token',
        redirectUri: environment.auth0.redirectUri,
        scope: 'email openid profile'
      }),
    },
  ],
})
export class AuthModule { }
