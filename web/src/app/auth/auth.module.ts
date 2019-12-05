import { NgModule } from '@angular/core';
import { StoreModule } from '@ngrx/store';
import { WebAuth } from 'auth0-js';

import { StorageModule } from '../storage/storage.module';

import { environment } from '../../environments/environment';

@NgModule({
  imports: [
    StorageModule,
    StoreModule,
  ],
  providers: [
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
