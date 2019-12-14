import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';
import { WebAuth } from 'auth0-js';
import { environment } from '../../environments/environment';
import { StorageModule } from '../storage';

@NgModule({
  imports: [
    StorageModule,
    StoreModule,
  ],
  providers: [
    {
      provide: WebAuth,
      useFactory: () => new WebAuth({
        domain: 'nomkhonwaan.auth0.com',
        clientID: environment.auth0.clientId,
        responseType: 'token id_token',
        redirectUri: environment.auth0.redirectUri,
        scope: 'email openid profile',
        audience: environment.auth0.audience,
      }),
    },
  ],
})
export class AuthModule { }
